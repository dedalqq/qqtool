# qqtool

`qqtool` is the network tool for testing connection and run server to another case

## Build

`make`

or

```
go get github.com/dedalqq/qqtool
go build qqtool
```

## Use

* `qqtool host:port` - connect to TCP port and `host` (analog `telent host port` or `nc host:port`)
* `qqtool -t host:port` - connect to TCP port and `host` with TLS
* `qqtool -l :80` - run TCP server on port 80 (analog `nc -l -p 80` but not one-off)
* `qqtool -l :443 -f host:port` - run TCP server on port 443 and forward incoming connection to `host` and `port`

## Man

`qqtool -h`

```
Usage of ./qqtool:
  -f string
    	Forward incoming connection to host. (with -l)
  -l string
    	Listening TCP port. Example: :80 or 0.0.0.0:80
  -r	Repeat the connection after disconnecting (simple host connection)
  -s string
    	Path to folder for save data from incoming connections (with -l or simple host connection)
  -t	Connect with TLS. (With -f or simple host connection)
```