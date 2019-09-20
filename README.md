##Tendermint block miss exporter

This prometheus exporter watches a tendermint node and counts the blocks missed by a specific validator.

The `tm_mon_misses` starts from zero and is incremented every time the specified validator misses a block.

In order to track whether the exporter is not stuck it also exposes the last processed height as `tm_mon_height`

###Configuration

The configuration needs to be passed in via the environment.

| Name | Description |
|------|-------------|
| LADDR | Listening address (e.g. `:8080`)
| RPC | RPC address of the tendermint node to use for monitoring |
| ADDRESS | The hex consensus address of the validator that should be monitored |

### Error handling

In case of errors during startup the program will panic. Errors during runtime are printed to the console and might
lead to the exporter not processing blocks. This will be visible in prometheus as `tm_mon_height` will stop increasing.
