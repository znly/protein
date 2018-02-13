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
// schemas and, as such, both their respective single & recursive hashes can be
// known in advance and shouldn't ever change.
//
// If any modification to either of these schemas were to happen, you'd have to
// modify the following expected values in order to fix the tests.. That is, if
// you're sure about what you're doing.
const (
	TEST_TSKnownName = ".test.TestSchema"
	TEST_DEKnownName = ".test.TestSchema.DepsEntry"
	TEST_GTKnownName = ".test.TestSchema.GhostType"

	TEST_TSKnownHashSingle = "PROT-00ede76a8940ef0f5d9022ecbca679d9"
	TEST_DEKnownHashSingle = "PROT-624796f94565bcdd2e785ef24a037ebb"
	TEST_GTKnownHashSingle = "PROT-f4d9460136ed7169a701ac2bab5a642b"

	TEST_TSKnownHashRecurse = "PROT-8b244a1a35e88f1e1aad8915dd603021"
	TEST_DEKnownHashRecurse = "PROT-4f6928d2737ba44dac0e3df123f80284"
	TEST_GTKnownHashRecurse = "PROT-3fecf73710581dfb3f46718988b9316e"
)
