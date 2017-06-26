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
	"fmt"
	"sort"
	"strings"
	"testing"

	"go.uber.org/zap"

	"github.com/stretchr/testify/assert"
	"github.com/znly/protein/protobuf/test"
	"github.com/znly/protein/protoscan"
)

// -----------------------------------------------------------------------------

// This example demonstrates how to use Protein's `ScanSchemas` function in
// order to sniff all the locally instanciated schemas into a `SchemaMap`,
// then walk over this map in order to print each schema as well as its
// dependencies.
func ExampleScanSchemas() {
	// sniff local protobuf schemas into a `SchemaMap` using a MD5 hasher,
	// and prefixing each resulting UID with 'PROT-'
	sm, err := ScanSchemas(protoscan.MD5, "PROT-")
	if err != nil {
		zap.L().Fatal(err.Error())
	}

	// walk over the map to print the schemas and their respective dependencies
	var output string
	sm.ForEach(func(ps *ProtobufSchema) error {
		// discard schemas not orginating from protein's test package
		if !strings.HasPrefix(ps.GetFQName(), ".test") {
			return nil
		}
		output += fmt.Sprintf("[%s] %s\n", ps.GetSchemaUID(), ps.GetFQName())
		for uid, name := range ps.GetDeps() {
			// discard dependencies not orginating from protein's test package
			if strings.HasPrefix(name, ".test") {
				output += fmt.Sprintf(
					"[%s] depends on: [%s] %s\n", ps.GetSchemaUID(), uid, name,
				)
			}
		}
		return nil
	})

	// sort rows so the output stays deterministic
	rows := strings.Split(output, "\n")
	sort.Strings(rows)
	for i, r := range rows { // prettying
		if strings.Contains(r, "depends on") {
			rows[i] = "\t" + strings.Join(strings.Split(r, " ")[1:], " ")
		}
	}
	output = strings.Join(rows, "\n")
	fmt.Println(output)

	// Output:
	// [PROT-048ddab197df688302a76296293ba101] .test.OtherTestSchemaXXX
	// [PROT-1ed0887b99e22551676141f133ee3813] .test.TestSchemaXXX
	// 	depends on: [PROT-048ddab197df688302a76296293ba101] .test.OtherTestSchemaXXX
	// 	depends on: [PROT-393cb6dc1b4fc350cf10ca99f429301d] .test.TestSchemaXXX.WeatherType
	// 	depends on: [PROT-6926276ca6306966d1a802c3b8f75298] .test.TestSchemaXXX.IdsEntry
	// 	depends on: [PROT-c43da9745d68bd3cb97dc0f4905f3279] .test.TestSchemaXXX.NestedEntry
	// 	depends on: [PROT-f6be24770f6e8d5edc8ef12c94a23010] .test.TestSchemaXXX.DepsEntry
	// [PROT-393cb6dc1b4fc350cf10ca99f429301d] .test.TestSchemaXXX.WeatherType
	// [PROT-4f6928d2737ba44dac0e3df123f80284] .test.TestSchema.DepsEntry
	// [PROT-528e869395e00dd5525c2c2c69bbd4d0] .test.TestSchema
	// 	depends on: [PROT-4f6928d2737ba44dac0e3df123f80284] .test.TestSchema.DepsEntry
	// [PROT-6926276ca6306966d1a802c3b8f75298] .test.TestSchemaXXX.IdsEntry
	// [PROT-c43da9745d68bd3cb97dc0f4905f3279] .test.TestSchemaXXX.NestedEntry
	// [PROT-f6be24770f6e8d5edc8ef12c94a23010] .test.TestSchemaXXX.DepsEntry
	// 	depends on: [PROT-c43da9745d68bd3cb97dc0f4905f3279] .test.TestSchemaXXX.NestedEntry
}

// -----------------------------------------------------------------------------

func TestProtoscan_ScanSchemas(t *testing.T) {
	sm, err := ScanSchemas(protoscan.SHA256, "PROT-")
	assert.Nil(t, err)

	// should at least find the `.protoscan.TestSchema` and its nested
	// `DepsEntry` message in the returned protobuf schemas

	ps := sm.GetByUID(protoscan.TEST_TSKnownHashRecurse)
	assert.NotNil(t, ps)
	assert.Equal(t, protoscan.TEST_TSKnownName, ps.GetFQName())
	assert.Equal(t, protoscan.TEST_TSKnownHashRecurse, ps.GetSchemaUID())
	assert.NotNil(t, ps.GetDescr())
	assert.NotEmpty(t, ps.GetDeps())
	assert.NotNil(t, ps.GetDeps()[protoscan.TEST_DEKnownHashRecurse])

	de := sm.GetByUID(protoscan.TEST_DEKnownHashRecurse)
	assert.NotNil(t, de)
	assert.Equal(t, protoscan.TEST_DEKnownName, de.GetFQName())
	assert.Equal(t, protoscan.TEST_DEKnownHashRecurse, de.GetSchemaUID())
	assert.NotNil(t, de.GetDescr())
	assert.Empty(t, de.GetDeps())
}

func TestProtoscan_LoadSchemas(t *testing.T) {
	protoFiles := map[string][]byte{
		"test_schema.proto": test.FileDescriptorTestSchema,
	}
	sm, err := LoadSchemas(protoFiles, protoscan.SHA256, "PROT-")
	assert.Nil(t, err)

	// should at least find the `.protoscan.TestSchema` and its nested
	// `DepsEntry` message in the returned protobuf schemas

	ps := sm.GetByUID(protoscan.TEST_TSKnownHashRecurse)
	assert.NotNil(t, ps)
	assert.Equal(t, protoscan.TEST_TSKnownName, ps.GetFQName())
	assert.Equal(t, protoscan.TEST_TSKnownHashRecurse, ps.GetSchemaUID())
	assert.NotNil(t, ps.GetDescr())
	assert.NotEmpty(t, ps.GetDeps())
	assert.NotNil(t, ps.GetDeps()[protoscan.TEST_DEKnownHashRecurse])

	de := sm.GetByUID(protoscan.TEST_DEKnownHashRecurse)
	assert.NotNil(t, de)
	assert.Equal(t, protoscan.TEST_DEKnownName, de.GetFQName())
	assert.Equal(t, protoscan.TEST_DEKnownHashRecurse, de.GetSchemaUID())
	assert.NotNil(t, de.GetDescr())
	assert.Empty(t, de.GetDeps())
}
