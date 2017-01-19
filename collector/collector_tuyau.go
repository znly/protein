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

package collector

import (
	"github.com/gogo/protobuf/proto"
	"github.com/znly/tuyauDB/pipe"
)

// -----------------------------------------------------------------------------

type Optor func(ct *CollectorTuyau)

func SetPipeDriver(p pipe.Pipe) Optor {
	return func(ct *CollectorTuyau) { ct.pipe = p }
}

// -----------------------------------------------------------------------------

// CollectorTuyau implements a Collector that stores collected protobuf schemas
// in an internal in-memory cache.
//
// It integrates with with znly/tuyauDB in order to Publish() the content of its
// internal cache into a pipe.Pipe.
type CollectorTuyau struct {
	pipe pipe.Pipe
}

func New(opts ...Optor) Collector {
	ct := &CollectorTuyau{}
	for _, o := range opts {
		o(ct) // apply option
	}
	return ct
}

// -----------------------------------------------------------------------------

func (c *CollectorTuyau) Collect() error              { return nil }
func (c *CollectorTuyau) Hash(s proto.Message) string { return "" }
func (c *CollectorTuyau) Add(hash []byte, d proto.Descriptor) error {
	return nil
}
func (c *CollectorTuyau) Publish() error {
	return nil
}
