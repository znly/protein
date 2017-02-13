// Copyright © 2016 Zenly <hello@zen.ly>.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package protoscan

import (
	"reflect"
	"strings"
	"unsafe"

	"github.com/gogo/protobuf/protoc-gen-gogo/descriptor"
	"github.com/kardianos/osext"
	"github.com/pkg/errors"
	"github.com/znly/protein/protobuf/schemas"
	"github.com/znly/protein/protoscan/internal/objfile"
)

// -----------------------------------------------------------------------------

// ScanSchemas retrieves every protobuf schema instanciated by any of the
// currently loaded protobuf libraries (e.g. golang/protobuf, gogo/protobuf...),
// computes the dependency graphs that link them, then finally returns a map of
// ProtobufSchema objects (which are protobuf objects themselves) using each
// schema's unique, deterministic & versioned identifier as key.
//
// Note that ProtobufSchemas' UIDs are always prefixed with "PROT-".
//
// This unique key is generated based on the binary representation of the
// schema and of its dependency graph: this implies that the key will change if
// any of the schema's dependency is modified in any way.
// In the end, this means that, as the schema and/or its dependencies follow
// their natural evolution, each and every historic version of it will have
// been stored with their own unique identifier.
//
// `failOnDuplicate` is an optional parameter that defaults to true; have a
// look at ScanSchemas' implementation to understand what it does and when (if
// ever) would you need to set it to false instead.
//
// Have a look at 'protoscan.go' and 'descriptor_tree.go' for more information
// about how all of this works; the code is heavily documented.
func ScanSchemas(failOnDuplicate ...bool) (map[string]*schemas.ProtobufSchema, error) {
	fod := true
	if len(failOnDuplicate) > 0 {
		fod = failOnDuplicate[0]
	}

	// get local pointers to proto.protoFiles instances
	protoFiles, err := BindProtofileSymbols()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	// unzip everything into a map of FileDescriptorProtos using the path of
	// the original .proto as key
	fdps := map[string]*descriptor.FileDescriptorProto{}
	for _, maps := range protoFiles {
		for file, descr := range *maps {
			// If a FileDescriptorProto already exists for this .proto
			// (i.e. another protobuf package has already instanciated a type of
			//  the same name) and `failOnDuplicate` is true (which is what it
			// defaults to), then we immediately stop everything and return
			// an error.
			//
			// You can disable this check by setting `failOnDuplicate` to false,
			// but be aware that if this condition ever returns true, either:
			// - you know exactly what you're doing and that is what you expected
			//   to happen (i.e. some FDPs will be overwritten)
			// - there is something seriously wrong with your setup and things
			//   are going to take a turn for the worst pretty soon; hence you're
			//   better off crashing right now
			if _, ok := fdps[file]; ok && fod {
				return nil, errors.Errorf("`%s` is instanciated multiple times", file)
			}
			fdp, err := UnzipAndUnmarshal(descr)
			if err != nil {
				log.Error(err)
				continue
			}
			fdps[file] = fdp
		}
	}

	dtsByUID, err := NewDescriptorTrees(fdps)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	// builds slice of schemas.ProtobufSchema objects
	pss := make(map[string]*schemas.ProtobufSchema, len(dtsByUID))
	for uid, dt := range dtsByUID {
		ps := &schemas.ProtobufSchema{
			UID:    uid,
			FQName: dt.FQName(),
			Deps:   map[string]string{},
		}
		switch descr := dt.descr.(type) {
		case *descriptor.DescriptorProto:
			ps.Descr = &schemas.ProtobufSchema_Message{descr}
		case *descriptor.EnumDescriptorProto:
			ps.Descr = &schemas.ProtobufSchema_Enum{descr}
		default:
			return nil, errors.Errorf("`%v`: illegal type", reflect.TypeOf(descr))
		}
		for _, depUID := range dt.DependencyUIDs() {
			dep, ok := dtsByUID[depUID]
			if !ok {
				return nil, errors.Errorf("missing dependency")
			}
			ps.Deps[depUID] = dep.FQName()
		}
		pss[uid] = ps
	}

	return pss, nil
}

// BindProtofileSymbols loads the currently running executable in memory
// using Go's private objfile API and then loops over its symbols in order
// to find `protoFiles` variables.
// NOTE: since its an internal API, the objfile package and its dependencies
//       had to be fully copied into this project, see protoscan/internal/.
//
// These `protoFiles` variables are maintained by the various protobuf
// libraries out there (i.e. golang/protobuf, gogo/protobuf & other
// implementations) in order to keep track of the FileDescriptorProtos
// that have been loaded at boot-time (see proto.RegisterFile).
// This essentially means that each and every protobuf schema known
// by the currently running program is stored into these maps.
//
//
// There are two main issues that need to be worked around though:
//
// A. `proto.protoFiles` is a package-level private variable and, as such,
//    cannot be (AFAIK) accessed by any means except by forking the original
//    package, which is not a viable option here.
//
// B. Because of how vendoring works, there can actually be an infinite amount
//    of `proto.protoFiles` variables instanciated at runtime, and we must
//    get ahold of each and every one of them.
//
// Considering the above issues, doing some hacking with the symbols seem
// to be the smart(er) way to go here.
// As `proto.protoFiles` variables are declared as a package-level globals,
// their respective virtual addresses are known at compile-time and stored
// in the executable: what we're doing here is we find those addresses then
// apply some unsafe-foo magic in order to create local pointers that point
// to these addresses.
//
// And, voila!
func BindProtofileSymbols() (map[string]*map[string][]byte, error) {
	var protoFilesBindings map[string]*map[string][]byte

	binPath, err := osext.Executable()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	bin, err := objfile.Open(binPath)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer bin.Close()
	syms, err := bin.Symbols()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	protoFilesBindings = make(map[string]*map[string][]byte, len(syms))
	for _, s := range syms {
		if strings.HasSuffix(s.Name, "/proto.protoFiles") {
			p := (*map[string][]byte)(unsafe.Pointer(uintptr(s.Addr)))
			log.Infof("found symbol `%s` @ %p", s.Name, p)
			protoFilesBindings[s.Name] = p
		}
	}

	return protoFilesBindings, nil
}
