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
	"strings"
	"unsafe"

	"github.com/kardianos/osext"
	"github.com/pkg/errors"
	"github.com/znly/protein/protoscan/internal/objfile"
)

// -----------------------------------------------------------------------------

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
//    cannot be (to my knowledge at least) accessed by any means except
//    by forking the original package, which is not a viable option here.
//
// B. Because of how vendoring works, there can actually be an infinite amount
//    of `proto.protoFiles` variables instanciated at runtime, and we must
//    to get ahold of each and every one of them.
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
func BindProtofileSymbols() ([]*map[string][]byte, error) {
	var protoFilesBindings []*map[string][]byte

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

	protoFilesBindings = make([]*map[string][]byte, 0, len(syms))
	for _, s := range syms {
		if strings.HasSuffix(s.Name, "/proto.protoFiles") {
			p := (*map[string][]byte)(unsafe.Pointer(uintptr(s.Addr)))
			log.Infof("found symbol `%s` @ %p", s.Name, p)
			protoFilesBindings = append(protoFilesBindings, p)
		}
	}

	return protoFilesBindings, nil
}
