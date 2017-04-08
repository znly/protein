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

// -----------------------------------------------------------------------------

// `TestSchema` and `DepsEntry` are both immutable, for-testing-purposes
// only schemas and, as such, both their respective single & recursive
// hashes can be known in advance and shouldn't ever change unless the
// `ByteSSlice.Hash` method is ever modified.
//
// If such a modification of `ByteSSlice.Hash`'s behavior were to happen,
// you'd have to modify the following expected values in order to fix the
// tests.. That is, if you're sure about what you're doing.
const (
	TEST_TSKnownName = ".test.TestSchema"
	TEST_DEKnownName = ".test.TestSchema.DepsEntry"

	TEST_TSKnownHashSingle = "PROT-d12b5690ec5877c9c5f6c84b57523b14291b46de50ffaf487a8a3f04d8198bff"
	TEST_DEKnownHashSingle = "PROT-39b441710129f5f7e4c92863b02173db7aa71658799031d6a1dedca260842dfa"

	TEST_TSKnownHashRecurse = "PROT-aae11ece4778cf8da20b7e436958feebcc0a1237807866603d1c197f27a3cb5b"
	TEST_DEKnownHashRecurse = "PROT-d278f5561f05e68f6e68fcbc6b801d29a69b4bf6044bf3e6242ea8fe388ebd6e"
)
