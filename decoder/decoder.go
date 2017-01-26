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

package decoder

import "github.com/gogo/protobuf/proto"

// -----------------------------------------------------------------------------

// Implementers of the Decoder interface expose methods to decode protobuf
// messages.
//
// The default implementation, as seen in decoder/decoder_versioned.go,
// integrates with znly/tuyauDB in order to add support for versioned protobuf
// objects via a central registry.
type Decoder interface {
	Decode(payload []byte) (map[string]interface{}, error)
	DecodeAs(payload []byte, dst proto.Message) error
}