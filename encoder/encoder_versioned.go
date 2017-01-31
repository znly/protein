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
	"github.com/pkg/errors"
	"github.com/znly/protein"
	"github.com/znly/protein/bank"
)

// -----------------------------------------------------------------------------

// Versioned implements an Encoder that embeds a Bank in order to augment the
// protobuf payloads that it encodes with additional versioning metadata.
type Versioned struct {
	b bank.Bank
}

func NewVersioned(b bank.Bank) *Versioned { return &Versioned{b: b} }

// -----------------------------------------------------------------------------

// Encode marshals the given protobuf message then wraps it up into a
// ProtobufPayload object that adds additional versioning metadata.
//
// Encode uses the message's fully-qualified name to reverse-lookup its UID.
// Note that a single FQ-name might point to multiple UIDs if multiple versions
// of the associated message are currently availaible in the bank.
// When this happens, the first UID from the returned list will be used (and
// since this list is randomly-ordered, effectively a random UID will be used).
func (v *Versioned) Encode(o proto.Message) ([]byte, error) {
	// marshal the actual protobuf message
	payload, err := proto.Marshal(o)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	// find the first UID associated with the fully-qualified name of `o`
	uids := v.b.FQNameToUID("." + proto.MessageName(o))
	if len(uids) <= 0 {
		return nil, errors.Errorf("`%s`: FQ-name not found in bank")
	}
	// wrap the marshaled payload within a ProtobufPayload message
	pp := &protein.ProtobufPayload{
		UID:     uids[0],
		Payload: payload,
	}
	// marshal the ProtobufPayload
	return proto.Marshal(pp)
}
