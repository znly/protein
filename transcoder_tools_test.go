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
	"testing"

	"github.com/stretchr/testify/assert"
)

// -----------------------------------------------------------------------------

type S struct{ X struct{ Y int } }

func TestFieldByName(t *testing.T) {
	v := &S{X: struct{ Y int }{Y: 42}}

	assert.Equal(t, FieldByName(v, ""), nil)
	assert.Equal(t, FieldByName(v, "X"), struct{ Y int }{Y: 42})
	assert.Equal(t, FieldByName(v, "X."), nil)
	assert.Equal(t, FieldByName(v, "XY"), nil)
	assert.Equal(t, FieldByName(v, "X.Y"), int(42))
	assert.Equal(t, FieldByName(v, "X.Y."), nil)
	assert.Equal(t, FieldByName(v, "X.YZ"), nil)
	assert.Equal(t, FieldByName(v, "X.Y.Z"), nil)
}

func TestSetByName(t *testing.T) {
	var v *S

	v = &S{X: struct{ Y int }{Y: 42}}
	assert.Error(t, SetByName(v, "", 42))
	assert.Equal(t, &S{X: struct{ Y int }{Y: 42}}, v)

	v = &S{X: struct{ Y int }{Y: 42}}
	assert.Error(t, SetByName(v, "X", 66))
	assert.Equal(t, &S{X: struct{ Y int }{Y: 42}}, v)
	assert.NoError(t, SetByName(v, "X", struct{ Y int }{Y: 66}))
	assert.Equal(t, &S{X: struct{ Y int }{Y: 66}}, v)

	v = &S{X: struct{ Y int }{Y: 42}}
	assert.Error(t, SetByName(v, "X.", 66))
	assert.Equal(t, &S{X: struct{ Y int }{Y: 42}}, v)

	v = &S{X: struct{ Y int }{Y: 42}}
	assert.Error(t, SetByName(v, "XY", 66))
	assert.Equal(t, &S{X: struct{ Y int }{Y: 42}}, v)

	v = &S{X: struct{ Y int }{Y: 42}}
	assert.NoError(t, SetByName(v, "X.Y", 66))
	assert.Equal(t, &S{X: struct{ Y int }{Y: 66}}, v)

	v = &S{X: struct{ Y int }{Y: 42}}
	assert.Error(t, SetByName(v, "X.Y.", 66))
	assert.Equal(t, &S{X: struct{ Y int }{Y: 42}}, v)

	v = &S{X: struct{ Y int }{Y: 42}}
	assert.Error(t, SetByName(v, "X.YZ", 66))
	assert.Equal(t, &S{X: struct{ Y int }{Y: 42}}, v)

	v = &S{X: struct{ Y int }{Y: 42}}
	assert.Error(t, SetByName(v, "X.Y.Z", 66))
	assert.Equal(t, &S{X: struct{ Y int }{Y: 42}}, v)
}
