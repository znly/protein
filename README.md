<p align="center">
  <img src="resources/pics/protein.png" alt="Protein"/>
</p>

# Protein

*Protein* is an encoding/decoding library for [*Protobuf*](https://developers.google.com/protocol-buffers/) that comes with schema-versioning and runtime-decoding capabilities.

It has diverse use-cases, including but not limited to:
- setting up schema registries
- decoding Protobuf payloads without the need to know their schema at compile-time
- identifying & preventing applicative bugs and data corruption issues
- creating custom-made container formats for on-disk storage
- ...and more!

**Table of Contents:**  
<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->


- [Performance](#performance)
- [Error handling](#error-handling)
- [Logging](#logging)
- [Monitoring](#monitoring)
- [Contributing](#contributing)
  - [Running tests](#running-tests)
  - [Running benchmarks](#running-benchmarks)
- [Authors](#authors)
- [See also](#see-also)
- [License](#license-)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## Performance

Configuration:
```
MacBook Pro (Retina, 15-inch, Mid 2015)
Intel(R) Core(TM) i7-4870HQ CPU @ 2.50GHz
```

Encoding:
```
## gogo/protobuf ##

BenchmarkTranscoder_Encode/gogo/protobuf        300000    4287 ns/op
BenchmarkTranscoder_Encode/gogo/protobuf-2     1000000    2195 ns/op
BenchmarkTranscoder_Encode/gogo/protobuf-4     1000000    1268 ns/op
BenchmarkTranscoder_Encode/gogo/protobuf-8     1000000    1258 ns/op
BenchmarkTranscoder_Encode/gogo/protobuf-24    1000000    1536 ns/op

## znly/protein ##

BenchmarkTranscoder_Encode/znly/protein         300000    5556 ns/op
BenchmarkTranscoder_Encode/znly/protein-2       500000    2680 ns/op
BenchmarkTranscoder_Encode/znly/protein-4      1000000    1638 ns/op
BenchmarkTranscoder_Encode/znly/protein-8      1000000    1798 ns/op
BenchmarkTranscoder_Encode/znly/protein-24     1000000    2288 ns/op
```

Decoding:
```
## gogo/protobuf ##

BenchmarkTranscoder_DecodeAs/gogo/protobuf      300000    5970 ns/op
BenchmarkTranscoder_DecodeAs/gogo/protobuf-2    500000    3226 ns/op
BenchmarkTranscoder_DecodeAs/gogo/protobuf-4   1000000    2125 ns/op
BenchmarkTranscoder_DecodeAs/gogo/protobuf-8   1000000    2015 ns/op
BenchmarkTranscoder_DecodeAs/gogo/protobuf-24  1000000    2380 ns/op

## znly/protein ##

BenchmarkTranscoder_Decode/znly/protein        200000     6777 ns/op
BenchmarkTranscoder_Decode/znly/protein-2      500000     3986 ns/op
BenchmarkTranscoder_Decode/znly/protein-4      500000     2630 ns/op
BenchmarkTranscoder_Decode/znly/protein-8      500000     2973 ns/op
BenchmarkTranscoder_Decode/znly/protein-24     500000     3037 ns/op
```

See [*transcoder_test.go*](./transcoder_test.go) for the actual benchmarking code.

## Error handling

*Protein* uses the [`pkg/errors`](https://github.com/pkg/errors) package to handle error propagation throughout the call stack; please take a look at the related documentation for more information on how to properly handle these errors.

## Logging

*Protein* rarely logs, but when it does, it uses the global logger from Uber's [*Zap*](https://github.com/uber-go/zap) package.  
You can thus control the behavior of *Protein*'s logger however you like by calling [`zap.ReplaceGlobals`](https://godoc.org/go.uber.org/zap#ReplaceGlobals) at your convenience.

For more information, see *Zap*'s [documentation](https://godoc.org/go.uber.org/zap).

## Monitoring

*Protein* does not offer any kind of monitoring hooks, yet.

## Contributing

Contributions of any kind are welcome.

### Running tests

```sh
$ docker-compose -f test/docker-compose.yml up
$ ## wait for the datastores to be up & running
$ make test
```

### Running benchmarks

```sh
$ make bench
```

<!--## Internals overview-->

<!--See the blog post.-->

## Authors

See AUTHORS for the list of contributors.

## See also

- [*Docker-Protobuf*](https://github.com/znly/docker-protobuf): an all inclusive `protoc` suite for *Protobuf* and *gRPC*.

## License ![License](https://img.shields.io/badge/license-Apache2-blue.svg?style=plastic)

The Apache License version 2.0 (Apache2) - see LICENSE for more details.

Copyright (c) 2017	Zenly	<hello@zen.ly> [@zenlyapp](https://twitter.com/zenlyapp)
