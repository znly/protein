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

import (
	"github.com/gogo/protobuf/proto"
	"github.com/znly/tuyauDB/client"
)

// -----------------------------------------------------------------------------

type Optor func(d *DecoderVersioned)

func SetClient(c client.Client) Optor {
	return func(d *DecoderVersioned) { d.client = c }
}
func SetCollector(c collector.Collector) Optor {
	return func(d *DecoderVersioned) { d.collector = collector }
}

// -----------------------------------------------------------------------------

// DecoderVersioned implements a Decoder that integrates with with znly/tuyauDB
// in order to add support for versioned protobuf objects via a central registry.
type DecoderVersioned struct {
	client    client.Client
	collector collector.Collector
}

func New(opts ...Optor) Decoder {
	d := &DecoderVersioned{}
	for _, o := range opts {
		o(d) // apply option
	}
	return d
}

// -----------------------------------------------------------------------------

func (c *DecoderVersioned) Decode(payload []byte) (map[string]interface{}, error) {
	return nil, nil
}
func (c *DecoderVersioned) DecodeAs(payload []byte, dst proto.Message) error {
	return nil
}
