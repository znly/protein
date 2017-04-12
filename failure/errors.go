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

package failure

import "github.com/pkg/errors"

// -----------------------------------------------------------------------------

// Error is an error returned either by `protein` or any of its sub-packages.
//
// Use `IsProteinError(err)` to know whether or not an error originates from
// this package.
type Error int

const (
	/* Common */
	ErrUnknown            Error = iota // unknown error
	ErrSchemaNotFound     Error = iota // schema's UID not found
	ErrDependencyNotFound Error = iota // dependency's UID not found

	/* Protostruct */
	ErrSchemaNotMessageType   Error = iota // schema is not of messsage type
	ErrFieldTypeNotSupported  Error = iota // field type not supported
	ErrFieldLabelNotSupported Error = iota // field label not supported

	/* Protoscan */
	ErrFDAlreadyInstanciated Error = iota // proto-file-descriptor instanciated multiple times
	ErrFDUnknownType         Error = iota // proto-file-descriptor is of unknown type
	ErrFDMissingDependency   Error = iota // proto-file-descriptor depends on missing schemas
)

// IsProteinError returns true if `err` originates from this package.
func IsProteinError(err error) bool { _, ok := errors.Cause(err).(Error); return ok }

func (e Error) Error() string {
	switch e {
	/* Common */
	case ErrUnknown:
		return "error: unknown"
	case ErrSchemaNotFound:
		return "error: schema not found"
	case ErrDependencyNotFound:
		return "error: dependency not found"

	/* Protostruct */
	case ErrSchemaNotMessageType:
		return "error [protostruct]: schema is not of message type"
	case ErrFieldTypeNotSupported:
		return "error [protostruct]: this field type is not supported"
	case ErrFieldLabelNotSupported:
		return "error [protostruct]: this kind of field label is not supported"

	/* Protoscan */
	case ErrFDAlreadyInstanciated:
		return "error [protoscan]: this protobuf file-descriptor is instanciated multiple times"
	case ErrFDUnknownType:
		return "error [protoscan]: this protobuf file-descriptor is of unknown type"
	case ErrFDMissingDependency:
		return "error [protoscan]: this protobuf file-descriptor depends on missing protobuf schemas"

	default:
		return "error: undefined"
	}
}
