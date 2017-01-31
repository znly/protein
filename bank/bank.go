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

package bank

import "github.com/znly/protein"

// -----------------------------------------------------------------------------

// Implementers of the Bank interface expose methods to store and retrieve
// ProtobufSchemas.
//
// The default implementation, as seen in bank/tuyau.go, integrates with
// znly/tuyauDB in order to keep its local in-memory cache in sync with a
// TuyauDB store.
type Bank interface {
	Get(uid string) (map[string]*protein.ProtobufSchema, error)
	FQNameToUID(fqName string) []string
	Put(ps ...*protein.ProtobufSchema) error
}
