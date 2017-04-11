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

package protein

import (
	"reflect"
	"sync"

	"github.com/gogo/protobuf/protoc-gen-gogo/descriptor"
	"github.com/pkg/errors"
	"github.com/znly/protein/protoscan"
)

// -----------------------------------------------------------------------------

// TODO(cmc)
type SchemaMap struct {
	lock *sync.RWMutex
	// map[schemaUID]ProtobufSchema
	schemaMap map[string]*ProtobufSchema
	// reverse-mapping of fully-qualified names to UIDs
	// map[FQNname][]{schemaUID}
	revmap map[string][]string
}

// TODO(cmc)
func NewSchemaMap() *SchemaMap {
	return &SchemaMap{
		lock:      &sync.RWMutex{},
		schemaMap: make(map[string]*ProtobufSchema, 0),
		revmap:    make(map[string][]string, 0),
	}
}

// TODO(cmc)
func (sm *SchemaMap) Add(schemas map[string]*ProtobufSchema) *SchemaMap {
	sm.lock.Lock()
	for _, s := range schemas {
		sm.schemaMap[s.UID] = s
		sm.revmap[s.FQName] = append(sm.revmap[s.FQName], s.UID)
	}
	sm.lock.Unlock() // avoid defer()
	return sm
}

// TODO(cmc)
func (sm *SchemaMap) Size() int {
	sm.lock.RLock()
	sz := len(sm.schemaMap)
	sm.lock.RUnlock() // avoid defer()
	return sz
}

// TODO(cmc)
func (sm *SchemaMap) ForEach(f func(s *ProtobufSchema) error) error {
	sm.lock.RLock()
	for _, s := range sm.schemaMap {
		if err := f(s); err != nil {
			sm.lock.RUnlock() // avoid defer()
			return errors.WithStack(err)
		}
	}
	sm.lock.RUnlock() // avoid defer()
	return nil
}

// TODO(cmc)
func (sm *SchemaMap) GetByUID(uid string) *ProtobufSchema {
	sm.lock.RLock()
	s := sm.schemaMap[uid]
	sm.lock.RUnlock() // avoid defer()
	return s
}

// TODO(cmc)
func (sm *SchemaMap) GetByFQName(fqName string) *ProtobufSchema {
	sm.lock.RLock()
	uids := sm.revmap[fqName]
	if len(uids) <= 0 {
		return nil
	}
	s := sm.schemaMap[uids[0]]
	sm.lock.RUnlock() // avoid defer()
	return s
}

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
// Have a look at the `protoscan` sub-package for more information about how all
// of this works; the code is heavily documented.
func ScanSchemas(failOnDuplicate ...bool) (*SchemaMap, error) {
	fod := true
	if len(failOnDuplicate) > 0 {
		fod = failOnDuplicate[0]
	}

	// get local pointers to proto.protoFiles instances
	protoFiles, err := protoscan.BindProtofileSymbols()
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
				// TODO(cmc): real error
				return nil, errors.Errorf("`%s` is instanciated multiple times", file)
			}
			fdp, err := protoscan.UnzipAndUnmarshal(descr)
			if err != nil {
				log.Error(err.Error())
				continue
			}
			fdps[file] = fdp
		}
	}

	dtsByUID, err := protoscan.NewDescriptorTrees(fdps)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	pss := make(map[string]*ProtobufSchema, len(dtsByUID))
	for uid, dt := range dtsByUID {
		ps := &ProtobufSchema{
			UID:    uid,
			FQName: dt.FQName(),
			Deps:   map[string]string{},
		}
		switch descr := dt.Descr().(type) {
		case *descriptor.DescriptorProto:
			ps.Descr = &ProtobufSchema_Message{Message: descr}
		case *descriptor.EnumDescriptorProto:
			ps.Descr = &ProtobufSchema_Enum{Enum: descr}
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

	return NewSchemaMap().Add(pss), nil
}
