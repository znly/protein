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

	"go.uber.org/zap"

	"github.com/gogo/protobuf/protoc-gen-gogo/descriptor"
	"github.com/pkg/errors"
	"github.com/znly/protein/failure"
	"github.com/znly/protein/protoscan"
)

// -----------------------------------------------------------------------------

// SchemaMap is a thread-safe mapping & reverse-mapping of `ProtobufSchema`s.
//
// It atomically maintains two data-structures in parallel:
// a map of schemaUIDs to `ProtobufSchema`s, and
// a map of fully-qualified schema names to schemaUIDs.
//
// The `SchemaMap` is the main data-structure behind a `Transcoder`, used to
// store and retrieve every `ProtobufSchema`s that have been cached locally.
type SchemaMap struct {
	lock *sync.RWMutex
	// map[schemaUID]ProtobufSchema
	schemaMap map[string]*ProtobufSchema
	// reverse-mapping of fully-qualified names to UIDs
	// map[FQNname][]{schemaUID}
	revmap map[string][]string
}

// NewSchemaMap returns a new SchemaMap with its maps & locks pre-allocated.
func NewSchemaMap() *SchemaMap {
	return &SchemaMap{
		lock:      &sync.RWMutex{},
		schemaMap: make(map[string]*ProtobufSchema, 0),
		revmap:    make(map[string][]string, 0),
	}
}

// Add walks over the given map of `schemas` and add them to the `SchemaMap`
// while making sure to atomically maintain both the internal map and reverse-map.
func (sm *SchemaMap) Add(schemas map[string]*ProtobufSchema) *SchemaMap {
	sm.lock.Lock()
	for _, s := range schemas {
		zap.L().Debug("schema found",
			zap.String("uid", s.SchemaUID), zap.String("fq-name", s.FQName),
		)
		sm.schemaMap[s.SchemaUID] = s
		sm.revmap[s.FQName] = append(sm.revmap[s.FQName], s.SchemaUID)
	}
	sm.lock.Unlock() // avoid defer()
	return sm
}

// Size returns the number of schemas stored in the `SchemaMap`.
func (sm *SchemaMap) Size() int {
	sm.lock.RLock()
	sz := len(sm.schemaMap)
	sm.lock.RUnlock() // avoid defer()
	return sz
}

// ForEach applies the specified function `f` to every entry in the `SchemaMap`.
//
// ForEach is guaranteed to see a consistent view of the internal mapping.
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

// GetByUID returns the `ProtobufSchema` associated with the specified `schemaUID`.
//
// This is thread-safe.
func (sm *SchemaMap) GetByUID(schemaUID string) *ProtobufSchema {
	sm.lock.RLock()
	s := sm.schemaMap[schemaUID]
	sm.lock.RUnlock() // avoid defer()
	return s
}

// GetByFQName returns the first `ProtobufSchema`s that matches the specified
// fully-qualified name (e.g. `.google.protobuf.timestamp`).
//
// If more than one version of the same schema are stored in the `SchemaMap`,
// a fully-qualified name will naturally point to several distinct schemaUIDs.
// When this happens, `GetByFQName` will always return the first one to had
// been inserted in the map.
//
// This is thread-safe.
func (sm *SchemaMap) GetByFQName(fqName string) *ProtobufSchema {
	sm.lock.RLock()
	uids := sm.revmap[fqName]
	if len(uids) <= 0 {
		sm.lock.RUnlock() // avoid defer()
		return nil
	}
	s := sm.schemaMap[uids[0]]
	sm.lock.RUnlock() // avoid defer()
	return s
}

// -----------------------------------------------------------------------------

// ScanSchemas retrieves every protobuf schema instanciated by any of the
// currently loaded protobuf libraries (e.g. `golang/protobuf`,
// `gogo/protobuf`...), computes the dependency graphs that link them, builds
// the `ProtobufSchema` objects then returns a new `SchemaMap` filled with all
// those schemas.
//
// A `ProtobufSchema` is a data-structure that holds a protobuf descriptor as
// well as a map of all its dependencies' schemaUIDs.
// A schemaUID uniquely & deterministically identifies a protobuf schema based
// on its descriptor and all of its dependencies' descriptors.
// It essentially is the versioned identifier of a schema.
// For more information, see `ProtobufSchema`'s documentation as well as the
// `protoscan` implementation, especially `descriptor_tree.go`.
//
// The specified `hasher` is the actual function used to compute these
// schemaUIDs, see the `protoscan.Hasher` documentation for more information.
// The `hashPrefix` string will be preprended to the resulting hash that was
// computed via the `hasher` function.
// E.g. by passing `protoscan.MD5` as a `Hasher` and `PROT-` as a `hashPrefix`,
// the resulting schemaUIDs will be of the form 'PROT-<MD5hex>'.
//
// As a schema and/or its dependencies follow their natural evolution, each
// and every historic version of them will thus have been stored with their
// own unique identifiers.
//
// `failOnDuplicate` is an optional parameter that defaults to true; have a
// look at `ScanSchemas` implementation to understand what it does and when (if
// ever) would you need to set it to false instead.
//
// Finally, have a look at the `protoscan` sub-packages as a whole for more
// information about how all of this machinery works; the code is heavily
// documented.
func ScanSchemas(
	hasher protoscan.Hasher, hashPrefix string, failOnDuplicate ...bool,
) (*SchemaMap, error) {
	// get local pointers to `proto.protoFiles` instances
	protoFiles, err := protoscan.BindProtofileSymbols()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return buildSchemas(protoFiles, hasher, hashPrefix, failOnDuplicate...)
}

// LoadSchemas is the exact same thing as `ScanSchemas` except for the fact
// that it uses the given set of user-specified `fileDescriptorProtos` instead
// of scanning for matching symbols.
//
// See `ScanSchemas`'s documentation for more information.
func LoadSchemas(fileDescriptorProtos map[string][]byte,
	hasher protoscan.Hasher, hashPrefix string, failOnDuplicate ...bool,
) (*SchemaMap, error) {
	protoFiles := map[string]*map[string][]byte{
		"user_specified": &fileDescriptorProtos,
	}
	return buildSchemas(protoFiles, hasher, hashPrefix, failOnDuplicate...)
}

func buildSchemas(protoFiles map[string]*map[string][]byte,
	hasher protoscan.Hasher, hashPrefix string, failOnDuplicate ...bool,
) (*SchemaMap, error) {
	fod := true
	if len(failOnDuplicate) > 0 {
		fod = failOnDuplicate[0]
	}

	// unzip everything into a map of `FileDescriptorProto`s using the path of
	// the original .proto as key
	fdps := map[string]*descriptor.FileDescriptorProto{}
	for _, maps := range protoFiles {
		for file, descr := range *maps {
			// If a `FileDescriptorProto` already exists for this .proto
			// (i.e. another protobuf package has already instanciated a type with
			//  the exact same fully-qualified name) and `failOnDuplicate` is
			// true (which is what it defaults to), then we immediately stop
			// everything and return an error.
			//
			// You can disable this check by setting `failOnDuplicate` to false,
			// but be aware that if this condition ever returns true, either:
			// - you know exactly what you're doing and that is what you expected
			//   to happen (i.e. some `FileDescriptorProto`s will be overwritten)
			// - there is something seriously wrong with your setup and things
			//   are going to take a turn for the worst pretty soon; hence you're
			//   better off crashing right now
			if _, ok := fdps[file]; ok && fod {
				return nil, errors.Wrapf(failure.ErrFDAlreadyInstanciated,
					"`%s` is instanciated multiple times", file,
				)
			}
			fdp, err := protoscan.UnzipAndUnmarshal(descr)
			if err != nil {
				zap.L().Error(err.Error())
				continue
			}
			fdps[file] = fdp
		}
	}

	dtsByUID, err := protoscan.NewDescriptorTrees(hasher, hashPrefix, fdps)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	pss := make(map[string]*ProtobufSchema, len(dtsByUID))
	for schemaUID, dt := range dtsByUID {
		ps := &ProtobufSchema{
			SchemaUID: schemaUID,
			FQName:    dt.FQName(),
			Deps:      map[string]string{},
		}
		switch descr := dt.Descr().(type) {
		case *descriptor.DescriptorProto:
			ps.Descr = &ProtobufSchema_Message{Message: descr}
		case *descriptor.EnumDescriptorProto:
			ps.Descr = &ProtobufSchema_Enum{Enum: descr}
		default:
			return nil, errors.Wrapf(failure.ErrFDUnknownType,
				"`%v`: unknown type", reflect.TypeOf(descr),
			)
		}
		for _, depUID := range dt.DependencyUIDs() {
			dep, ok := dtsByUID[depUID]
			if !ok {
				return nil, errors.Wrapf(failure.ErrDependencyNotFound,
					"`%s`: no dependency with this schemaUID", depUID,
				)
			}
			ps.Deps[depUID] = dep.FQName()
		}
		pss[schemaUID] = ps
	}

	return NewSchemaMap().Add(pss), nil
}
