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
	"context"

	"github.com/pkg/errors"
	"github.com/rainycape/memcache"
	"github.com/znly/protein/failure"
)

// -----------------------------------------------------------------------------

//type TranscoderGetter func(ctx context.Context, schemaUID string) ([]byte, error)
//type TranscoderSetter func(ctx context.Context, schemaUID string, payload []byte) error

// -----------------------------------------------------------------------------

/* Memcached */

// CreateTranscoderGetterMemcached returns a TranscoderGetter suitable for
// querying a binary blob from a memcached-compatible store.
//
// The specified context will be ignored.
func CreateTranscoderGetterMemcached(c *memcache.Client) TranscoderGetter {
	return func(ctx context.Context, schemaUID string) ([]byte, error) {
		item, err := c.Get(schemaUID)
		if err != nil {
			if err == memcache.ErrCacheMiss {
				return nil, errors.WithStack(failure.ErrSchemaNotFound)
			}
			return nil, errors.WithStack(err)
		}
		return item.Value, nil
	}
}

// CreateTranscoderSetterMemcached returns a TranscoderSetter suitable for
// setting a binary blob into a memcached-compatible store.
//
// The specified context will be ignored.
func CreateTranscoderSetterMemcached(c *memcache.Client) TranscoderSetter {
	return func(ctx context.Context, schemaUID string, payload []byte) error {
		return c.Set(&memcache.Item{
			Key:   schemaUID,
			Value: payload,
		})
	}
}

// -----------------------------------------------------------------------------

/* Redis */

// -----------------------------------------------------------------------------

/* Cassandra */
