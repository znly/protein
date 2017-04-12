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
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"sort"

	"github.com/pkg/errors"
)

// -----------------------------------------------------------------------------

// Hasher takes a pre-sorted slice of byte-slices as input and outputs a hashed
// representation of this data as a result.
//
// This Hasher, provided by the end-user, will be used to version every schema
// and associated dependencies found by the Protoscan package.
//
// This package provides some basic, ready-to-use hashers: MD5, SHA1, SHA256,
// SHA512.
type Hasher func(bss ByteSSlice) ([]byte, error)

// MD5 implements a Hasher using the MD5 hashing algorithm.
func MD5(bss ByteSSlice) ([]byte, error) {
	h := md5.New()
	var err error
	for _, bs := range bss {
		_, err = h.Write(bs)
		if err != nil {
			return nil, errors.WithStack(err)
		}
	}
	return h.Sum(nil), nil
}

// SHA1 implements a Hasher using the SHA1 hashing algorithm.
func SHA1(bss ByteSSlice) ([]byte, error) {
	h := sha1.New()
	var err error
	for _, bs := range bss {
		_, err = h.Write(bs)
		if err != nil {
			return nil, errors.WithStack(err)
		}
	}
	return h.Sum(nil), nil
}

// SHA256 implements a Hasher using the SHA256 hashing algorithm.
func SHA256(bss ByteSSlice) ([]byte, error) {
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

// SHA512 implements a Hasher using the SHA512 hashing algorithm.
func SHA512(bss ByteSSlice) ([]byte, error) {
	h := sha512.New()
	var err error
	for _, bs := range bss {
		_, err = h.Write(bs)
		if err != nil {
			return nil, errors.WithStack(err)
		}
	}
	return h.Sum(nil), nil
}

// -----------------------------------------------------------------------------

// ByteSSlice is a sortable slice of byte-slices.
//
// It is used to compute the schema hashes in DescriptorTree's implementation.
type ByteSSlice [][]byte

func (bss ByteSSlice) Len() int      { return len(bss) }
func (bss ByteSSlice) Swap(i, j int) { bss[i], bss[j] = bss[j], bss[i] }
func (bss ByteSSlice) Less(i, j int) bool {
	return bytes.Compare(bss[i], bss[j]) < 0
}
func (bss ByteSSlice) Sort() { sort.Sort(bss) }
