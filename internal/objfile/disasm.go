// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package objfile

import (
	"bufio"
	"bytes"
	"github.com/znly/protein/internal/src"
	"container/list"
	"debug/gosym"
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"text/tabwriter"

	"golang.org/x/arch/arm/armasm"
	"golang.org/x/arch/ppc64/ppc64asm"
	"golang.org/x/arch/x86/x86asm"
)

// Disasm is a disassembler for a given File.
type Disasm struct {
	syms      []Sym            //symbols in file, sorted by address
	pcln      Liner            // pcln table
	text      []byte           // bytes of text segment (actual instructions)
	textStart uint64           // start PC of text
	textEnd   uint64           // end PC of text
	goarch    string           // GOARCH string
	disasm    disasmFunc       // disassembler function for goarch
	byteOrder binary.ByteOrder // byte order for goarch
}

// Disasm returns a disassembler for the file f.
func (f *File) Disasm() (*Disasm, error) {
	syms, err := f.Symbols()
	if err != nil {
		return nil, err
	}

	pcln, err := f.PCLineTable()
	if err != nil {
		return nil, err
	}

	textStart, textBytes, err := f.Text()
	if err != nil {
		return nil, err
	}

	goarch := f.GOARCH()
	disasm := disasms[goarch]
	byteOrder := byteOrders[goarch]
	if disasm == nil || byteOrder == nil {
		return nil, fmt.Errorf("unsupported architecture")
	}

	// Filter out section symbols, overwriting syms in place.
	keep := syms[:0]
	for _, sym := range syms {
		switch sym.Name {
		case "runtime.text", "text", "_text", "runtime.etext", "etext", "_etext":
			// drop
		default:
			keep = append(keep, sym)
		}
	}
	syms = keep
	d := &Disasm{
		syms:      syms,
		pcln:      pcln,
		text:      textBytes,
		textStart: textStart,
		textEnd:   textStart + uint64(len(textBytes)),
		goarch:    goarch,
		disasm:    disasm,
		byteOrder: byteOrder,
	}

	return d, nil
}

// lookup finds the symbol name containing addr.
func (d *Disasm) lookup(addr uint64) (name string, base uint64) {
	i := sort.Search(len(d.syms), func(i int) bool { return addr < d.syms[i].Addr })
	if i > 0 {
		s := d.syms[i-1]
		if s.Addr != 0 && s.Addr <= addr && addr < s.Addr+uint64(s.Size) {
			return s.Name, s.Addr
		}
	}
	return "", 0
}

// base returns the final element in the path.
// It works on both Windows and Unix paths,
// regardless of host operating system.
func base(path string) string {
	path = path[strings.LastIndex(path, "/")+1:]
	path = path[strings.LastIndex(path, `\`)+1:]
	return path
}

// CachedFile contains the content of a file split into lines.
type CachedFile struct {
	FileName string
	Lines    [][]byte
}

// FileCache is a simple LRU cache of file contents.
type FileCache struct {
	files  *list.List
	maxLen int
}

// NewFileCache returns a FileCache which can contain up to maxLen cached file contents.
func NewFileCache(maxLen int) *FileCache {
	return &FileCache{
		files:  list.New(),
		maxLen: maxLen,
	}
}

// Line returns the source code line for the given file and line number.
// If the file is not already cached, reads it , inserts it into the cache,
// and removes the least recently used file if necessary.
// If the file is in cache, moves it up to the front of the list.
func (fc *FileCache) Line(filename string, line int) ([]byte, error) {
	if filepath.Ext(filename) != ".go" {
		return nil, nil
	}

	// Clean filenames returned by src.Pos.SymFilename()
	// or src.PosBase.SymFilename() removing
	// the leading src.FileSymPrefix.
	if strings.HasPrefix(filename, src.FileSymPrefix) {
		filename = filename[len(src.FileSymPrefix):]
	}

	// Expand literal "$GOROOT" rewrited by obj.AbsFile()
	filename = filepath.Clean(os.ExpandEnv(filename))

	var cf *CachedFile
	var e *list.Element

	for e = fc.files.Front(); e != nil; e = e.Next() {
		cf = e.Value.(*CachedFile)
		if cf.FileName == filename {
			break
		}
	}

	if e == nil {
		content, err := ioutil.ReadFile(filename)
		if err != nil {
			return nil, err
		}

		cf = &CachedFile{
			FileName: filename,
			Lines:    bytes.Split(content, []byte{'\n'}),
		}
		fc.files.PushFront(cf)

		if fc.files.Len() >= fc.maxLen {
			fc.files.Remove(fc.files.Back())
		}
	} else {
		fc.files.MoveToFront(e)
	}

	return cf.Lines[line-1], nil
}

// Print prints a disassembly of the file to w.
// If filter is non-nil, the disassembly only includes functions with names matching filter.
// If printCode is true, the disassembly includs corresponding source lines.
// The disassembly only includes functions that overlap the range [start, end).
func (d *Disasm) Print(w io.Writer, filter *regexp.Regexp, start, end uint64, printCode bool) {
	if start < d.textStart {
		start = d.textStart
	}
	if end > d.textEnd {
		end = d.textEnd
	}
	printed := false
	bw := bufio.NewWriter(w)

	var fc *FileCache
	if printCode {
		fc = NewFileCache(8)
	}

	for _, sym := range d.syms {
		symStart := sym.Addr
		symEnd := sym.Addr + uint64(sym.Size)
		relocs := sym.Relocs
		if sym.Code != 'T' && sym.Code != 't' ||
			symStart < d.textStart ||
			symEnd <= start || end <= symStart ||
			filter != nil && !filter.MatchString(sym.Name) {
			continue
		}
		if printed {
			fmt.Fprintf(bw, "\n")
		}
		printed = true

		file, _, _ := d.pcln.PCToLine(sym.Addr)
		fmt.Fprintf(bw, "TEXT %s(SB) %s\n", sym.Name, file)

		tw := tabwriter.NewWriter(bw, 18, 8, 1, '\t', tabwriter.StripEscape)
		if symEnd > end {
			symEnd = end
		}
		code := d.text[:end-d.textStart]

		var lastFile string
		var lastLine int

		d.Decode(symStart, symEnd, relocs, func(pc, size uint64, file string, line int, text string) {
			i := pc - d.textStart

			if printCode {
				if file != lastFile || line != lastLine {
					if srcLine, err := fc.Line(file, line); err == nil {
						fmt.Fprintf(tw, "%s%s%s\n", []byte{tabwriter.Escape}, srcLine, []byte{tabwriter.Escape})
					}

					lastFile, lastLine = file, line
				}

				fmt.Fprintf(tw, "  %#x\t", pc)
			} else {
				fmt.Fprintf(tw, "  %s:%d\t%#x\t", base(file), line, pc)
			}

			if size%4 != 0 || d.goarch == "386" || d.goarch == "amd64" {
				// Print instruction as bytes.
				fmt.Fprintf(tw, "%x", code[i:i+size])
			} else {
				// Print instruction as 32-bit words.
				for j := uint64(0); j < size; j += 4 {
					if j > 0 {
						fmt.Fprintf(tw, " ")
					}
					fmt.Fprintf(tw, "%08x", d.byteOrder.Uint32(code[i+j:]))
				}
			}
			fmt.Fprintf(tw, "\t%s\n", text)
		})
		tw.Flush()
	}
	bw.Flush()
}

// Decode disassembles the text segment range [start, end), calling f for each instruction.
func (d *Disasm) Decode(start, end uint64, relocs []Reloc, f func(pc, size uint64, file string, line int, text string)) {
	if start < d.textStart {
		start = d.textStart
	}
	if end > d.textEnd {
		end = d.textEnd
	}
	code := d.text[:end-d.textStart]
	lookup := d.lookup
	for pc := start; pc < end; {
		i := pc - d.textStart
		text, size := d.disasm(code[i:], pc, lookup, d.byteOrder)
		file, line, _ := d.pcln.PCToLine(pc)
		text += "\t"
		first := true
		for len(relocs) > 0 && relocs[0].Addr < i+uint64(size) {
			if first {
				first = false
			} else {
				text += " "
			}
			text += relocs[0].Stringer.String(pc - start)
			relocs = relocs[1:]
		}
		f(pc, uint64(size), file, line, text)
		pc += uint64(size)
	}
}

type disasmFunc func(code []byte, pc uint64, lookup x86asm.SymLookup, ord binary.ByteOrder) (text string, size int)

func disasm_386(code []byte, pc uint64, lookup x86asm.SymLookup, _ binary.ByteOrder) (string, int) {
	return disasm_x86(code, pc, lookup, 32)
}

func disasm_amd64(code []byte, pc uint64, lookup x86asm.SymLookup, _ binary.ByteOrder) (string, int) {
	return disasm_x86(code, pc, lookup, 64)
}

func disasm_x86(code []byte, pc uint64, lookup x86asm.SymLookup, arch int) (string, int) {
	inst, err := x86asm.Decode(code, 64)
	var text string
	size := inst.Len
	if err != nil || size == 0 || inst.Op == 0 {
		size = 1
		text = "?"
	} else {
		text = x86asm.GoSyntax(inst, pc, lookup)
	}
	return text, size
}

type textReader struct {
	code []byte
	pc   uint64
}

func (r textReader) ReadAt(data []byte, off int64) (n int, err error) {
	if off < 0 || uint64(off) < r.pc {
		return 0, io.EOF
	}
	d := uint64(off) - r.pc
	if d >= uint64(len(r.code)) {
		return 0, io.EOF
	}
	n = copy(data, r.code[d:])
	if n < len(data) {
		err = io.ErrUnexpectedEOF
	}
	return
}

func disasm_arm(code []byte, pc uint64, lookup x86asm.SymLookup, _ binary.ByteOrder) (string, int) {
	inst, err := armasm.Decode(code, armasm.ModeARM)
	var text string
	size := inst.Len
	if err != nil || size == 0 || inst.Op == 0 {
		size = 4
		text = "?"
	} else {
		text = armasm.GoSyntax(inst, pc, lookup, textReader{code, pc})
	}
	return text, size
}

func disasm_ppc64(code []byte, pc uint64, lookup x86asm.SymLookup, byteOrder binary.ByteOrder) (string, int) {
	inst, err := ppc64asm.Decode(code, byteOrder)
	var text string
	size := inst.Len
	if err != nil || size == 0 || inst.Op == 0 {
		size = 4
		text = "?"
	} else {
		text = ppc64asm.GoSyntax(inst, pc, lookup)
	}
	return text, size
}

var disasms = map[string]disasmFunc{
	"386":     disasm_386,
	"amd64":   disasm_amd64,
	"arm":     disasm_arm,
	"ppc64":   disasm_ppc64,
	"ppc64le": disasm_ppc64,
}

var byteOrders = map[string]binary.ByteOrder{
	"386":     binary.LittleEndian,
	"amd64":   binary.LittleEndian,
	"arm":     binary.LittleEndian,
	"ppc64":   binary.BigEndian,
	"ppc64le": binary.LittleEndian,
	"s390x":   binary.BigEndian,
}

type Liner interface {
	// Given a pc, returns the corresponding file, line, and function data.
	// If unknown, returns "",0,nil.
	PCToLine(uint64) (string, int, *gosym.Func)
}
