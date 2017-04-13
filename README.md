<p align="center">
  <img src="resources/pics/protein.png" alt="Protein"/>
</p>

# Protein

Open-source release of _ZProto_.

Do not commit anything Zenly-specific in this repository.

## Internals overview

layers of hashing

## Error handling

github.com/pkg/errors

## Logging

*Protein* rarely logs, but when it does, it uses the global logger from Uber's [*Zap*](https://github.com/uber-go/zap) package.  
You can thus control the behavior of *Protein*'s logger however you like by calling [`zap.ReplaceGlobals`](https://godoc.org/go.uber.org/zap#ReplaceGlobals) at your convenience.

For more information, see *Zap*'s [documentation](https://godoc.org/go.uber.org/zap).

## Monitoring

## Contributing

## Authors

See AUTHORS for the list of contributors.

## License ![License](https://img.shields.io/badge/license-Apache2-blue.svg?style=plastic)

The Apache License version 2.0 (Apache2) - see LICENSE for more details

Copyright (c) 2017	Zenly	<hello@zen.ly> [@zenlyapp](https://twitter.com/zenlyapp)
