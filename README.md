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

## License
