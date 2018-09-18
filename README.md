# qqtool

`qqtool` is the network tool for testing connection and run server to another case

## Build

`make`

## Use

* `qqtool host:port` - connect to TCP port and `host` (analog `telent host port` or `nc host:port`)
* `qqtool -t host:port` - connect to TCP port and `host` with TLS
* `qqtool -l :80` - run TCP server on port 80 (analog `nc -l -p 80` but not one-off)
* `qqtool -l :443 -f host:port` - run TCP server on port 443 and forward incoming connection to `host` and `port`