<p align="center">
  <img src="resources/pics/protein.png" alt="Protein"/>
</p>

# Protein

Open-source release of _ZProto_.

Do not commit anything Zenly-specific in this repository.

## Performance

*Protein* has basically no performance overhead compared to standard protobuf.

Intel(R) Core(TM) i7-4870HQ CPU @ 2.50GHz:  
```
BenchmarkTranscoder_DecodeAs/gogo/protobuf            	  200000	      6629 ns/op
BenchmarkTranscoder_DecodeAs/gogo/protobuf-2          	  500000	      3515 ns/op
BenchmarkTranscoder_DecodeAs/gogo/protobuf-4          	 1000000	      2292 ns/op
BenchmarkTranscoder_DecodeAs/gogo/protobuf-8          	 1000000	      2334 ns/op
BenchmarkTranscoder_DecodeAs/gogo/protobuf-24         	 1000000	      2527 ns/op
BenchmarkTranscoder_DecodeAs/znly/protein             	  200000	      6335 ns/op
BenchmarkTranscoder_DecodeAs/znly/protein-2           	  500000	      3391 ns/op
BenchmarkTranscoder_DecodeAs/znly/protein-4           	  500000	      2123 ns/op
BenchmarkTranscoder_DecodeAs/znly/protein-8           	 1000000	      2428 ns/op
BenchmarkTranscoder_DecodeAs/znly/protein-24          	  500000	      2760 ns/op
BenchmarkTranscoder_Decode/znly/protein               	  200000	      6355 ns/op
BenchmarkTranscoder_Decode/znly/protein-2             	  500000	      3658 ns/op
BenchmarkTranscoder_Decode/znly/protein-4             	 1000000	      2607 ns/op
BenchmarkTranscoder_Decode/znly/protein-8             	  500000	      2608 ns/op
BenchmarkTranscoder_Decode/znly/protein-24            	  500000	      3316 ns/op
```

## Error handling

github.com/pkg/errors

## Logging

*Protein* rarely logs, but when it does, it uses the global logger from Uber's [*Zap*](https://github.com/uber-go/zap) package.  
You can thus control the behavior of *Protein*'s logger however you like by calling [`zap.ReplaceGlobals`](https://godoc.org/go.uber.org/zap#ReplaceGlobals) at your convenience.

For more information, see *Zap*'s [documentation](https://godoc.org/go.uber.org/zap).

## Monitoring

## Internals overview

layers of hashing

## Contributing

## Authors

See AUTHORS for the list of contributors.

## License ![License](https://img.shields.io/badge/license-Apache2-blue.svg?style=plastic)

The Apache License version 2.0 (Apache2) - see LICENSE for more details

Copyright (c) 2017	Zenly	<hello@zen.ly> [@zenlyapp](https://twitter.com/zenlyapp)
