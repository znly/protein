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
	"strings"

	"github.com/pkg/errors"
)

// -----------------------------------------------------------------------------

// FieldByName is equivalent to `reflect.FieldByName` but allows for recursive
// names such as 'Field1.Field13.Field138'.
func FieldByName(st interface{}, nameR string) interface{} {
	return fieldByNameR(reflect.ValueOf(st), strings.Split(nameR, "."))
}
func fieldByNameR(st reflect.Value, nameParts []string) interface{} {
	if len(nameParts) <= 0 { // end of recursion
		if st.IsValid() {
			return st.Interface()
		}
		return nil
	}
	for st.Kind() == reflect.Ptr {
		st = st.Elem()
	}
	if st.Kind() != reflect.Struct {
		return nil
	}
	return fieldByNameR(st.FieldByName(nameParts[0]), nameParts[1:])
}

// SetByName is equivalent to `FieldByName` but allows for modifying the field
// pointed to by `nameR`.
func SetByName(st interface{}, nameR string, v interface{}) error {
	vt := reflect.ValueOf(st)
	if vt.Kind() != reflect.Ptr {
		return errors.New("'st' must be a pointer")
	}
	for vt.Kind() == reflect.Ptr {
		vt = vt.Elem()
	}
	if vt.Kind() != reflect.Struct {
		return errors.New("inner (non-leaf) fields must be of strut type")
	}
	return setByNameR(vt, strings.Split(nameR, "."), v)
}
func setByNameR(st reflect.Value, nameParts []string, v interface{}) error {
	if len(nameParts) <= 0 { // end of recursion
		if st.CanSet() {
			if vv := reflect.ValueOf(v); vv.Type() == st.Type() {
				st.Set(reflect.ValueOf(v))
				return nil
			}
			return errors.New("incompatible types")
		}
		return errors.New("unassignable field")
	}
	for st.Kind() == reflect.Ptr {
		st = st.Elem()
	}
	if st.Kind() != reflect.Struct {
		return errors.New("inner (non-leaf) fields must be of strut type")
	}
	return setByNameR(st.FieldByName(nameParts[0]), nameParts[1:], v)
}
