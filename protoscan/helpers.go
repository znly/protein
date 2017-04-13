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

import (
	"bytes"
	"compress/gzip"
	"io/ioutil"

	"github.com/gogo/protobuf/proto"
	"github.com/gogo/protobuf/protoc-gen-gogo/descriptor"
	"github.com/pkg/errors"
)

// -----------------------------------------------------------------------------

// UnzipAndUnmarshal g-unzips the given binary blob then unmarshals it into
// a `FileDescriptorProto`.
//
// It is typically used to decode the binary blobs generated by the protobuf
// compiler.
func UnzipAndUnmarshal(b []byte) (*descriptor.FileDescriptorProto, error) {
	r, err := gzip.NewReader(bytes.NewReader(b))
	if err != nil {
		return nil, errors.WithStack(err)
	}
	bu, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	var fdp descriptor.FileDescriptorProto
	err = proto.Unmarshal(bu, &fdp)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return &fdp, nil
}
