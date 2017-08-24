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
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"runtime"
	"sync"
	"testing"
	"time"

	"go.uber.org/zap"

	"github.com/garyburd/redigo/redis"
	"github.com/gogo/protobuf/proto"
	"github.com/gogo/protobuf/types"
	"github.com/stretchr/testify/assert"
	"github.com/znly/protein/protobuf/test"
	"github.com/znly/protein/protoscan"
)

// -----------------------------------------------------------------------------

var _testDur = time.Second*3 + time.Millisecond*100
var _transcoderTestSchemaXXX = &test.TestSchemaXXX{
	SchemaUID: "uuid-A",
	FQNames:   []string{"fqname-A", "fqname-B"},
	Deps: map[string]*test.TestSchemaXXX_NestedEntry{
		"dep-A": {Key: "dep-A-A", Value: "dep-A-A"},
		"dep-B": {Key: "dep-A-B", Value: "dep-A-B"},
	},
	Ids: map[int32]string{
		1: "id-A",
		2: "id-B",
	},
	Ts: types.Timestamp{
		Seconds: 42,
		Nanos:   43,
	},
	Ots: &test.OtherTestSchemaXXX{
		Ts: &types.Timestamp{
			Seconds: 942,
			Nanos:   943,
		},
	},
	Nss: []test.TestSchemaXXX_NestedEntry{
		{Key: "nss-key-A", Value: "nss-value-A"},
		{Key: "nss-key-B", Value: "nss-value-B"},
	},
	Weathers: []test.TestSchemaXXX_WeatherType{test.TestSchemaXXX_RAIN},
	TSStd:    time.Date(2003, time.April, 10, 10, 10, 10, 10, time.UTC),
	DurStd:   []*time.Duration{&_testDur},
}

var trc *Transcoder

func TestMain(m *testing.M) {
	var err error
	trc, err = NewTranscoder(context.Background(), protoscan.SHA256, "PROT-")
	if err != nil {
		zap.L().Fatal(err.Error())
	}

	os.Exit(m.Run())
}

// -----------------------------------------------------------------------------

// This example demonstrates the use the Protein package in order to:
// initialize a `Transcoder`,
// sniff the local protobuf schemas from memory,
// synchronize the local schema-database with a remote datastore (here `redis`),
// use a `Transcoder` to encode & decode protobuf payloads using an already
// known schema,
// use a `Transcoder` to decode protobuf payloads without any prior knowledge
// of their schema.
func Example() {
	// A local `redis` server must be up & running for this example to work:
	//
	//   $ docker run -p 6379:6379 --name my-redis --rm redis:3.2 redis-server

	// open up a new `redis` connection pool
	p := &redis.Pool{
		Dial: func() (redis.Conn, error) {
			return redis.DialURL("redis://localhost:6379/0")
		},
	}
	defer p.Close()

	/* INITIALIZATION */

	// initialize a `Transcoder` that is transparently kept in-sync with
	// a `redis` datatore.
	trc, err := NewTranscoder(
		// this context defines the timeout & deadline policies when pushing
		// schemas to the local `redis`; i.e. it is forwarded to the
		// `TranscoderSetter` function that is passed below
		context.Background(),
		// the schemas found in memory will be versioned using a MD5 hash
		// algorithm, prefixed by the 'PROT-' string
		protoscan.MD5, "PROT-",
		// configure the `Transcoder` to push every protobuf schema it can find
		// in memory into the specified `redis` connection pool
		TranscoderOptSetter(NewTranscoderSetterRedis(p)),
		// configure the `Transcoder` to query the given `redis` connection pool
		// when it cannot find a specific protobuf schema in its local cache
		TranscoderOptGetter(NewTranscoderGetterRedis(p)),
	)

	// At this point, the local `redis` datastore should contain all the
	// protobuf schemas known to the `Transcoder`, as defined by their respective
	// versioning hashes:
	//
	//   $ docker run -it --link my-redis:redis --rm redis:3.2 redis-cli -h redis -p 6379 -c KEYS '*'
	//
	//   1) "PROT-31c64ad1c6476720f3afee6881e6f257"
	//   2) "PROT-56b347c6c212d3176392ab9bf5bb21ee"
	//   3) "PROT-c2dbc910081a372f31594db2dc2adf72"
	//   4) "PROT-09595a7e58d28b081d967b69cb00e722"
	//   5) "PROT-05dc5bd440d980600ecc3f1c4a8e315d"
	//   6) "PROT-8cbb4e79fdeadd5f0ff0971bbf7de31e"
	//   ... etc ...

	/* ENCODING */

	// create a simple object to be serialized
	ts, _ := types.TimestampProto(time.Now())
	obj := &test.TestSchemaXXX{
		Ids: map[int32]string{
			42:  "the-answer",
			666: "the-devil",
		},
		Ts: *ts,
	}

	// wrap the object and its versioning metadata within a `ProtobufPayload`
	// object, then serialize the bundle as a protobuf binary blob
	payload, err := trc.Encode(obj)
	if err != nil {
		log.Fatal(err)
	}

	/* DECODING */

	var myObj test.TestSchemaXXX

	// /!\ This will fail in cryptic ways since vanilla protobuf is unaware
	// of how Protein bundles the versioning metadata within the payload...
	// Don't do this!
	_ = proto.Unmarshal(payload, &myObj)

	// this will properly unbundle the data from the metadata before
	// unmarshalling the payload
	err = trc.DecodeAs(payload, &myObj)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("A:", myObj.Ids[42]) // prints the answer!

	/* RUNTIME-DECODING */

	// empty the content of the local cache of protobuf schemas in order to
	// make sure that the `Transcoder` will have to lazily fetch the schema
	// and its dependencies from the `redis` datastore during the decoding
	trc.sm = NewSchemaMap()

	// the `Transcoder` will do a lot of stuff behind the scenes so it can
	// successfully decode the payload:
	// 1. the versioning metadata is extracted from the payload
	// 2. the corresponding schema as well as its dependencies are lazily
	//    fetched from the `redis` datastore (using the `TranscoderGetter` that
	//    was passed to the constructor)
	// 3. a structure-type definition is created from these schemas using Go's
	//    reflection APIs, with the right protobuf tags & hints for the protobuf
	//    deserializer to do its thing
	// 4. an instance of this structure is created, then the payload is
	//    unmarshalled into it
	myRuntimeObj, err := trc.Decode(context.Background(), payload)
	if err != nil {
		log.Fatal(err)
	}
	myRuntimeIDs := myRuntimeObj.Elem().FieldByName("IDs")
	fmt.Println("B:", myRuntimeIDs.MapIndex(reflect.ValueOf(int32(666)))) // prints the devil!

	// Output:
	// A: the-answer
	// B: the-devil
}

// -----------------------------------------------------------------------------

// NOTE: Simple cache coherence test, not really necessary anymore as the rest
// of the tests encompass this and much more but, well, it's free anyway.
func TestTranscoder_localCache(t *testing.T) {
	var expectedUID string
	var revUID string

	// `.test.TestSchema` should be in there
	expectedUID = "PROT-aae11ece4778cf8da20b7e436958feebcc0a1237807866603d1c197f27a3cb5b"
	revUID = trc.sm.GetByFQName(".test.TestSchema").SchemaUID
	assert.NotEmpty(t, revUID)
	assert.Equal(t, expectedUID, revUID)
	schemas, err := trc.GetAndUpsert(context.Background(), revUID)
	assert.Nil(t, err)
	assert.NotEmpty(t, schemas)
	assert.Equal(t, 2, len(schemas)) // `.test.TestSchema` + nested `DepsEntry`

	// `.test.TestSchema.DepsEntry` should be in there
	expectedUID = "PROT-d278f5561f05e68f6e68fcbc6b801d29a69b4bf6044bf3e6242ea8fe388ebd6e"
	revUID = trc.sm.GetByFQName(".test.TestSchema.DepsEntry").SchemaUID
	assert.NotEmpty(t, revUID)
	assert.Equal(t, expectedUID, revUID)
	schemas, err = trc.GetAndUpsert(context.Background(), revUID)
	assert.Nil(t, err)
	assert.NotEmpty(t, schemas)
	assert.Equal(t, 1, len(schemas)) // `.test.TestSchema.DepsEntry` only
}

// -----------------------------------------------------------------------------

func TestTranscoder_Encode(t *testing.T) {
	tsExpected := &test.TestSchema{
		Uid:    "test-uuid",
		FqName: "test-schema",
		Deps: map[string]string{
			"test": "schema",
		},
	}
	payload, err := trc.Encode(tsExpected)
	assert.Nil(t, err)
	assert.NotNil(t, payload)

	var pp ProtobufPayload
	var ts test.TestSchema
	assert.Nil(t, proto.Unmarshal(payload, &pp))
	uidExpected := "PROT-aae11ece4778cf8da20b7e436958feebcc0a1237807866603d1c197f27a3cb5b"
	assert.Equal(t, uidExpected, pp.GetSchemaUID())
	assert.NotEmpty(t, pp.GetPayload())
	assert.Nil(t, proto.Unmarshal(pp.GetPayload(), &ts))
	assert.Equal(t, tsExpected, &ts)
}

func BenchmarkTranscoder_Encode(b *testing.B) {
	b.Run("gogo/protobuf", func(b *testing.B) {
		b.SetParallelism(3)
		b.RunParallel(func(pb *testing.PB) {
			var payload []byte
			var err error
			for pb.Next() {
				payload, err = proto.Marshal(_transcoderTestSchemaXXX)
				assert.Nil(b, err)
				assert.NotNil(b, payload)
			}
		})
	})
	b.Run("znly/protein/implicit-fqname", func(b *testing.B) {
		b.SetParallelism(3)
		b.RunParallel(func(pb *testing.PB) {
			var payload []byte
			var err error
			for pb.Next() {
				payload, err = trc.Encode(_transcoderTestSchemaXXX)
				assert.Nil(b, err)
				assert.NotNil(b, payload)
			}
		})
	})
	b.Run("znly/protein/explicit-fqname", func(b *testing.B) {
		b.SetParallelism(3)
		b.RunParallel(func(pb *testing.PB) {
			var payload []byte
			var err error
			for pb.Next() {
				payload, err = trc.Encode(
					_transcoderTestSchemaXXX, "test.TestSchemaXXX",
				)
				assert.Nil(b, err)
				assert.NotNil(b, payload)
			}
		})
	})
}

// -----------------------------------------------------------------------------

func TestTranscoder_DecodeAs(t *testing.T) {
	payload, err := trc.Encode(_transcoderTestSchemaXXX)
	assert.Nil(t, err)
	assert.NotNil(t, payload)

	var ts test.TestSchemaXXX
	assert.Nil(t, trc.DecodeAs(payload, &ts))
	assert.Equal(t, _transcoderTestSchemaXXX, &ts)
}

func BenchmarkTranscoder_DecodeAs(b *testing.B) {
	payload, err := trc.Encode(_transcoderTestSchemaXXX)
	assert.Nil(b, err)
	assert.NotNil(b, payload)
	payloadRaw, err := proto.Marshal(_transcoderTestSchemaXXX)
	assert.Nil(b, err)
	assert.NotNil(b, payloadRaw)
	b.Run("gogo/protobuf", func(b *testing.B) {
		b.SetParallelism(3)
		b.RunParallel(func(pb *testing.PB) {
			var ts test.TestSchemaXXX
			var err error
			for pb.Next() {
				err = proto.Unmarshal(payloadRaw, &ts)
				assert.Nil(b, err)
				assert.Equal(b, _transcoderTestSchemaXXX.FQNames, ts.FQNames)
			}
		})
	})
	b.Run("znly/protein", func(b *testing.B) {
		b.SetParallelism(3)
		b.RunParallel(func(pb *testing.PB) {
			var ts test.TestSchemaXXX
			var err error
			for pb.Next() {
				err = trc.DecodeAs(payload, &ts)
				assert.Nil(b, err)
				assert.Equal(b, _transcoderTestSchemaXXX.FQNames, ts.FQNames)
			}
		})
	})
}

func TestTranscoder_Decode(t *testing.T) {
	payload, err := trc.Encode(_transcoderTestSchemaXXX)
	assert.Nil(t, err)
	assert.NotNil(t, payload)

	v, err := trc.Decode(context.Background(), payload)
	assert.Nil(t, err)
	assertFieldValues(t, reflect.ValueOf(_transcoderTestSchemaXXX), v)
}

func assertFieldValues(t *testing.T, expected, actual reflect.Value) {
	switch expected.Kind() {
	case reflect.Ptr:
		assertFieldValues(t, expected.Elem(), actual.Elem())
	case reflect.Map:
		for _, ek := range expected.MapKeys() {
			assertFieldValues(t, expected.MapIndex(ek), actual.MapIndex(ek))
		}
	case reflect.Slice, reflect.Array:
		for i := 0; i < expected.Len(); i++ {
			assertFieldValues(t, expected.Index(i), actual.Index(i))
		}
	case reflect.Struct:
		for i := 0; i < expected.NumField(); i++ {
			f := expected.Type().Field(i)
			if f.PkgPath == "time" {
				return // let's not go there, few have come back
			}
			name := f.Name
			switch name { // due to Protein's camel-casing features
			case "Ids":
				name = "IDs"
			case "Ts":
				name = "TS"
			}
			assertFieldValues(t, expected.Field(i), actual.FieldByName(name))
		}
	default:
		// NOTE: Enums are not fully supported yet, we handle them as simple
		// scalar types (i.e. 'int32').
		if e, ok := expected.Interface().(test.TestSchemaXXX_WeatherType); ok {
			assert.True(t, reflect.DeepEqual(int32(e), actual.Interface()))
		} else {
			assert.True(t,
				reflect.DeepEqual(expected.Interface(), actual.Interface()),
			)
		}
	}
}

func BenchmarkTranscoder_Decode(b *testing.B) {
	payload, err := trc.Encode(_transcoderTestSchemaXXX)
	assert.Nil(b, err)
	assert.NotNil(b, payload)
	b.Run("znly/protein", func(b *testing.B) {
		b.SetParallelism(3)
		b.RunParallel(func(pb *testing.PB) {
			var ts reflect.Value
			var err error
			for pb.Next() {
				ts, err = trc.Decode(context.Background(), payload)
				assert.Nil(b, err)
				assert.NotEqual(b, reflect.ValueOf(nil), ts)
			}
		})
	})
}

// -----------------------------------------------------------------------------

func TestTranscoder_Parallelism(t *testing.T) {
	payload, err := trc.Encode(_transcoderTestSchemaXXX)
	assert.Nil(t, err)
	assert.NotNil(t, payload)

	nbRoutines := runtime.GOMAXPROCS(0) * 10
	nbOps := 1000

	t.Run("Encode", func(t *testing.T) {
		t.Parallel()
		wg := &sync.WaitGroup{}
		wg.Add(nbRoutines)
		for p := 0; p < nbRoutines; p++ {
			go func() {
				var payload []byte
				var err error
				for i := 0; i < nbOps; i++ {
					payload, err = trc.Encode(_transcoderTestSchemaXXX)
					assert.Nil(t, err)
					assert.NotNil(t, payload)
				}
				wg.Done()
			}()
		}
		wg.Wait()
	})
	t.Run("Decode", func(t *testing.T) {
		t.Parallel()
		wg := &sync.WaitGroup{}
		wg.Add(nbRoutines)
		for p := 0; p < nbRoutines; p++ {
			go func() {
				var ts reflect.Value
				var err error
				for i := 0; i < nbOps; i++ {
					ts, err = trc.Decode(context.Background(), payload)
					assert.Nil(t, err)
					assert.NotEqual(t, reflect.ValueOf(nil), ts)
				}
				wg.Done()
			}()
		}
		wg.Wait()
	})
	t.Run("DecodeAs", func(t *testing.T) {
		t.Parallel()
		wg := &sync.WaitGroup{}
		wg.Add(nbRoutines)
		for p := 0; p < nbRoutines; p++ {
			go func() {
				var ts test.TestSchemaXXX
				var err error
				for i := 0; i < nbOps; i++ {
					err = trc.DecodeAs(payload, &ts)
					assert.Nil(t, err)
					assert.Equal(t, _transcoderTestSchemaXXX.FQNames, ts.FQNames)
				}
				wg.Done()
			}()
		}
		wg.Wait()
	})
}

// -----------------------------------------------------------------------------

func TestTranscoder_SaveState_LoadState(t *testing.T) {
	f, err := ioutil.TempFile("", "")
	assert.NoError(t, err)
	path := f.Name()
	assert.NoError(t, f.Close())

	assert.NoError(t, trc.SaveState(path))

	trcEmpty, err := NewTranscoderFromSchemaMap(context.Background(),
		NewSchemaMap())
	assert.NoError(t, err)
	assert.NoError(t, trcEmpty.LoadState(path))

	for schemaUID, schema := range trc.sm.schemaMap {
		assert.Equal(t, schemaUID, trcEmpty.sm.schemaMap[schemaUID].SchemaUID)
		assert.Equal(t,
			schema.String(), trcEmpty.sm.schemaMap[schemaUID].String())
	}
}
