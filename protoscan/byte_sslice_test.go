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
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
)

// -----------------------------------------------------------------------------

func TestProtoscan_ByteSSlice_Hash(t *testing.T) {
	myNameIsGiovanniGiorgio := []byte("My name is Giovanni Giorgio,")
	butEverybodyCallsMe := []byte("but everybody calls me...")
	giorgio := []byte("Giorgo.")

	bss := ByteSSlice{myNameIsGiovanniGiorgio, butEverybodyCallsMe, giorgio}

	bss.Sort()
	assert.Equal(t, giorgio, bss[0])
	assert.Equal(t, myNameIsGiovanniGiorgio, bss[1])
	assert.Equal(t, butEverybodyCallsMe, bss[2])

	expectedH := "f27edac4d321e0b20a955c3b2d1d77cb6331eab6954e02cb5621c37a9775869f"
	h, err := SHA256(bss)
	assert.Nil(t, err)
	assert.NotNil(t, h)
	assert.Equal(t, expectedH, hex.EncodeToString(h))
}
