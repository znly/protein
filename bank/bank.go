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

import (
	"context"

	"github.com/znly/protein/protobuf/schemas"
)

// -----------------------------------------------------------------------------

// Implementers of the Bank interface expose methods to store and retrieve
// ProtobufSchemas.
//
// The default implementation, as seen in bank/tuyau.go, integrates with
// znly/tuyauDB in order to keep its local in-memory cache in sync with a
// TuyauDB store.
type Bank interface {
	Get(ctx context.Context, uid string) (map[string]*schemas.ProtobufSchema, error)
	FQNameToUID(fqName string) []string
	Put(ctx context.Context, ps ...*schemas.ProtobufSchema) error
}
