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
	"bytes"
	"crypto/sha256"
	"sort"

	"github.com/pkg/errors"
)

// -----------------------------------------------------------------------------

// ByteSSlice is a sortable slice of slice of bytes.
//
// It is mostly used for computing recursive hashes in DescriptorTree's
// implementation.
type ByteSSlice [][]byte

func (bss ByteSSlice) Len() int      { return len(bss) }
func (bss ByteSSlice) Swap(i, j int) { bss[i], bss[j] = bss[j], bss[i] }
func (bss ByteSSlice) Less(i, j int) bool {
	return bytes.Compare(bss[i], bss[j]) < 0
}
func (bss ByteSSlice) Sort() { sort.Sort(bss) }

// Hash implements the hashing logic behind every hash value generated inside
// the protoscan package.
//
// It first sorts `bss`, then computes the SHA-1 hash of the result of appending
// every entry in the sorted slice.
func (bss ByteSSlice) Hash() ([]byte, error) {
	bss.Sort()
	h := sha256.New()
	var err error
	for _, bs := range bss {
		_, err = h.Write(bs)
		if err != nil {
			return nil, errors.WithStack(err)
		}
	}
	return h.Sum(nil), nil
}
