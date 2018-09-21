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

* `qqtool host:port` - connect to TCP port and `host` (analog `telent host port` or `nc host port`)
* `qqtool -t host:port` - connect to TCP port and `host` with TLS
* `qqtool -l :80` - run TCP server on port 80 (analog `nc -l -p 80` but not one-off)
* `qqtool -l :443 -t -f host:port` - run TCP server on port 443 and forward incoming connection to `host` and `port`
* `qqtool -s /tmp host:port` - connect to TCP port and host and save dump to file in `/tmp` folder
* `qqtool -r host:port` - connect and retry connection after disconnecting
* `qqtool -l :80 -F /tmp` - run simple file server on port 80 and in `/tmp` dir
* `qqtool -l :80 -e /bin/bash` - run server on 80 port and run `/bin/bash` for incomming connection (analog `nc -l -p 80 -e /bin/bash` but not one-off)

## Man

`qqtool -h`

```
Usage of ./qqtool:
  -F string
        Run file server on path (with -l)
  -e string
        Run command for incomming connections (with -l)
  -f string
        Forward incoming connection to host. (with -l)
  -l string
        Listening TCP port. Example: :80 or 0.0.0.0:80
  -r    Repeat the connection after disconnecting (simple host connection)
  -s string
        Path to folder for save data from incoming connections (with -l or simple host connection)
  -t    Connect with TLS. (With -f or simple host connection)
```