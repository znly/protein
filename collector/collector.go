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

import "github.com/gogo/protobuf/proto"

// -----------------------------------------------------------------------------

// Implementers of the Collector interface expose methods to collect, keep
// track of, and publish available protobuf schemas.
//
// The default implementation, as seen in collector/collector_tuyau.go,
// integrates with znly/tuyauDB in order to Publish() the content of its
// internal in-memory cache into a pipe.Pipe.
type Collector interface {
	Collect() error
	Hash(s proto.Message) string
	Add(hash []byte, d proto.Descriptor) error
	Publish() error
}
