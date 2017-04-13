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
	"fmt"

	"go.uber.org/zap"

	"github.com/garyburd/redigo/redis"
	"github.com/gocql/gocql"
	"github.com/pkg/errors"
	"github.com/rainycape/memcache"
	"github.com/znly/protein/failure"
)

// -----------------------------------------------------------------------------

/* Memcached */

// NewTranscoderGetterMemcached returns a `TranscoderGetter` suitable for
// querying a binary blob from a memcached-compatible store.
//
// The specified context will be ignored.
func NewTranscoderGetterMemcached(c *memcache.Client) TranscoderGetter {
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

// NewTranscoderSetterMemcached returns a `TranscoderSetter` suitable for
// setting a binary blob into a memcached-compatible store.
//
// The specified context will be ignored.
func NewTranscoderSetterMemcached(c *memcache.Client) TranscoderSetter {
	return func(ctx context.Context, schemaUID string, payload []byte) error {
		return c.Set(&memcache.Item{
			Key:   schemaUID,
			Value: payload,
		})
	}
}

// -----------------------------------------------------------------------------

/* Redis */

// NewTranscoderGetterRedis returns a `TranscoderGetter` suitable for
// querying a binary blob from a redis-compatible store.
//
// The specified context will be ignored.
func NewTranscoderGetterRedis(p *redis.Pool) TranscoderGetter {
	return func(_ context.Context, schemaUID string) ([]byte, error) {
		c := p.Get() // avoid defer()
		b, err := redis.Bytes(c.Do("GET", schemaUID))
		if err := c.Close(); err != nil {
			zap.L().Error(err.Error())
		}
		if err != nil {
			if err == redis.ErrNil {
				return nil, errors.WithStack(failure.ErrSchemaNotFound)
			}
			return nil, errors.WithStack(err)
		}
		return b, nil
	}
}

// NewTranscoderSetterRedis returns a `TranscoderSetter` suitable for
// setting a binary blob into a redis-compatible store.
//
// The specified context will be ignored.
func NewTranscoderSetterRedis(p *redis.Pool) TranscoderSetter {
	return func(_ context.Context, schemaUID string, payload []byte) error {
		c := p.Get() // avoid defer()
		_, err := c.Do("SET", schemaUID, payload)
		if err := c.Close(); err != nil {
			zap.L().Error(err.Error())
		}
		return errors.WithStack(err)
	}
}

// -----------------------------------------------------------------------------

/* Cassandra */

// NewTranscoderGetterCassandra returns a `TranscoderGetter` suitable for
// querying a binary blob from a cassandra-compatible store.
//
// The <table> column-family is expected to have (at least) the following
// columns:
//   TABLE (
//   	<keyCol> ascii,
//   	<dataCol> blob,
//   PRIMARY KEY (<keyCol>))
//
// The given context is forwarded to `gocql`.
func NewTranscoderGetterCassandra(s *gocql.Session,
	table, keyCol, dataCol string,
) TranscoderGetter {
	queryGet := fmt.Sprintf(
		`SELECT %s FROM %s WHERE %s = ?;`, dataCol, table, keyCol,
	)
	return func(ctx context.Context, schemaUID string) ([]byte, error) {
		var payload []byte
		if err := s.Query(queryGet, schemaUID).
			WithContext(ctx).
			Scan(&payload); err != nil {
			if err == gocql.ErrNotFound {
				return nil, errors.WithStack(failure.ErrSchemaNotFound)
			}
			return nil, errors.WithStack(err)
		}
		return payload, nil
	}
}

// NewTranscoderSetterCassandra returns a `TranscoderSetter` suitable for
// setting a binary blob into a redis-compatible store.
//
// The <table> column-family is expected to have (at least) the following
// columns:
//   TABLE (
//   	<keyCol> ascii,
//   	<dataCol> blob,
//   PRIMARY KEY (<keyCol>))
//
// The given context is forwarded to `gocql`.
func NewTranscoderSetterCassandra(s *gocql.Session,
	table, keyCol, dataCol string,
) TranscoderSetter {
	querySet := fmt.Sprintf(
		`INSERT INTO %s (%s, %s) VALUES (?, ?);`, table, keyCol, dataCol,
	)
	return func(ctx context.Context, schemaUID string, payload []byte) error {
		return errors.WithStack(s.Query(querySet, schemaUID, payload).Exec())
	}
}
