// +build ignore

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

	"github.com/gogo/protobuf/proto"
)

// -----------------------------------------------------------------------------

// Implementers of the Wirer interface expose methods to encode & decode
// protobuf messages.
//
// The default implementation, as seen in wirer/wirer_versioned.go,
// implements a Wirer that integrates with a Bank in order to augment the
// protobuf payloads that it encodes with additional versioning metadata.
// These metadata are then used by the internal decoder of the versioned Wirer
// to determinate how to decode an incoming payload on the wire.
type Wirer interface {
	Encode(o proto.Message) ([]byte, error)
	EncodeWithName(o proto.Message, fqName string) ([]byte, error)

	DecodeStruct(payload []byte) (*reflect.Value, error)
	DecodeMessage(payload []byte, dst proto.Message) error
}
