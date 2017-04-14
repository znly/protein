<p align="center">
  <img src="resources/pics/protein.png" alt="Protein"/>
</p>

# Protein

Open-source release of _ZProto_.

Do not commit anything Zenly-specific in this repository.

## Performance

Configuration:  
```
Intel(R) Core(TM) i7-4870HQ CPU @ 2.50GHz:  
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

BenchmarkTranscoder_Decode/gogo/protobuf        200000    6785 ns/op
BenchmarkTranscoder_Decode/gogo/protobuf-2      500000    3661 ns/op
BenchmarkTranscoder_Decode/gogo/protobuf-4     1000000    2521 ns/op
BenchmarkTranscoder_Decode/gogo/protobuf-8     1000000    2667 ns/op
BenchmarkTranscoder_Decode/gogo/protobuf-24     500000    3126 ns/op

## znly/protein ##

BenchmarkTranscoder_Decode/znly/protein        200000     6777 ns/op
BenchmarkTranscoder_Decode/znly/protein-2      500000     3986 ns/op
BenchmarkTranscoder_Decode/znly/protein-4      500000     2630 ns/op
BenchmarkTranscoder_Decode/znly/protein-8      500000     2973 ns/op
BenchmarkTranscoder_Decode/znly/protein-24     500000     3037 ns/op
```

## Error handling

*Protein* uses the [`pkg/errors`](https://github.com/pkg/errors) package to handle error propagation throughout the call stack; please take a look at its documentation for more information about how to properly handle these errors.

## Logging

*Protein* rarely logs, but when it does, it uses the global logger from Uber's [*Zap*](https://github.com/uber-go/zap) package.  
You can thus control the behavior of *Protein*'s logger however you like by calling [`zap.ReplaceGlobals`](https://godoc.org/go.uber.org/zap#ReplaceGlobals) at your convenience.

For more information, see *Zap*'s [documentation](https://godoc.org/go.uber.org/zap).

## Monitoring

## Internals overview

layers of hashing

## Contributing

### Running tests

```sh
$ docker-compose -f test/docker-compose.yml up
$ ## wait for datastores to be up & running
$ make test
```

### Running benchmarks

```sh
$ make bench
```

## Authors

See AUTHORS for the list of contributors.

## License ![License](https://img.shields.io/badge/license-Apache2-blue.svg?style=plastic)

The Apache License version 2.0 (Apache2) - see LICENSE for more details

Copyright (c) 2017	Zenly	<hello@zen.ly> [@zenlyapp](https://twitter.com/zenlyapp)
