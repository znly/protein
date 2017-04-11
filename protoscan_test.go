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
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/znly/protein/protoscan"
)

// -----------------------------------------------------------------------------

// TODO(cmc)
func ExampleScanSchemas() {
	sm, err := ScanSchemas()
	if err != nil {
		log.Fatal(err.Error())
	}

	sm.ForEach(func(ps *ProtobufSchema) error {
		fmt.Printf("[%s] %s\n", ps.GetUID(), ps.GetFQName())
		for uid, name := range ps.GetDeps() {
			fmt.Printf("\tdepends on: [%s] %s\n", uid, name)
		}
		return nil
	})
}

// -----------------------------------------------------------------------------

func TestProtoscan_ScanSchemas(t *testing.T) {
	sm, err := ScanSchemas()
	assert.Nil(t, err)

	// should at least find the `.protoscan.TestSchema` and its nested
	// `DepsEntry` message in the returned protobuf schemas

	ps := sm.GetByUID(protoscan.TEST_TSKnownHashRecurse)
	assert.NotNil(t, ps)
	assert.Equal(t, protoscan.TEST_TSKnownName, ps.GetFQName())
	assert.Equal(t, protoscan.TEST_TSKnownHashRecurse, ps.GetUID())
	assert.NotNil(t, ps.GetDescr())
	assert.NotEmpty(t, ps.GetDeps())
	assert.NotNil(t, ps.GetDeps()[protoscan.TEST_DEKnownHashRecurse])

	de := sm.GetByUID(protoscan.TEST_DEKnownHashRecurse)
	assert.NotNil(t, de)
	assert.Equal(t, protoscan.TEST_DEKnownName, de.GetFQName())
	assert.Equal(t, protoscan.TEST_DEKnownHashRecurse, de.GetUID())
	assert.NotNil(t, de.GetDescr())
	assert.Empty(t, de.GetDeps())
}
