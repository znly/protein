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

package encoder

import (
	"github.com/gogo/protobuf/proto"
	"github.com/znly/tuyauDB/client"
)

// -----------------------------------------------------------------------------

// Versioned implements an Encoder that integrates with with znly/tuyauDB
// in order to add support for versioned protobuf objects via a central registry.
type Versioned struct {
	client client.Client
}

func New(c client.Client) *Versioned {
	return nil
}

// -----------------------------------------------------------------------------

func (c *Versioned) Encode(o proto.Message) ([]byte, error) {
	return nil, nil
}
func (c *Versioned) EncodeAs(o proto.Message, hash []byte) ([]byte, error) {
	return nil, nil
}
