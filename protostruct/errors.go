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

package protostruct

// -----------------------------------------------------------------------------

type Error int

const (
	ErrUnknown                Error = iota // unknown error
	ErrSchemaNotFound         Error = iota // schema's UID not found
	ErrSchemaNotMessageType   Error = iota // schema is not of messsage type
	ErrFieldTypeNotSupported  Error = iota // field type not supported
	ErrFieldLabelNotSupported Error = iota // field label not supported
)

// IsProtostructError returns true if `e` is a protostruct.Error.
func IsProtostructError(e error) bool { _, ok := e.(Error); return ok }

func (e Error) Error() string {
	switch e {
	case ErrUnknown:
		return "error: unknown"
	case ErrSchemaNotFound:
		return "error: no such schema UID"
	case ErrSchemaNotMessageType:
		return "error: schema is not of message type"
	case ErrFieldTypeNotSupported:
		return "error: this field type is not supported"
	case ErrFieldLabelNotSupported:
		return "error: this kind of field label is not supported"
	default:
		return "error: wat"
	}
}
