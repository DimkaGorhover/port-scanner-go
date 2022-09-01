# Port-Scanner-GO

[![Tool Category](https://badgen.net/badge/Tool/Port%20Scanner/black)](https://github.com/nxenon/port-scanner-go)
[![APP Version](https://badgen.net/badge/Version/v1.2/red)](https://github.com/nxenon/port-scanner-go)
[![Go Version](https://badgen.net/badge/Go/1.19/blue)](https://golang.org/doc/go1.19)
[![License](https://badgen.net/badge/License/GPLv2/purple)](https://github.com/nxenon/port-scanner-go/blob/master/LICENSE)

Simple TCP port scanner in golang.

## Installation & Build

You have to have GO version 1.19

```shell
go build
```

## Build Docker Image

```shell
docker build --force-rm --rm --target release .
```

## Run

```
help:
    ./port-scanner-go --help
    ./port-scanner-go help
version:
    ./port-scanner-go --version
ports 1-32000:
    ./port-scanner-go --host 192.168.1.1
single port:
    ./port-scanner-go --host 192.168.1.1 --port 80
ports range:
    ./port-scanner-go --host 192.168.1.1 --port 1-1024
specific ports:
    ./port-scanner-go --host 192.168.1.1 --port 80,443,22
debug mode:
    ./port-scanner-go --host 192.168.1.1 --port 80 --debug
```

## Help

```
NAME:
   port-scanner-go - A new cli application

USAGE:
   port-scanner-go [global options] command [command options] [arguments...]

VERSION:
   docker

AUTHORS:
   M Amin Nasiri <Khodexenon@gmail.com>

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --concurrency value  How many concurrency job run (default: CPU count)
   --debug              Enable Debug Logs (default: false)
   --help, -h           show help (default: false)
   --host value         Target Host
   --port value         Ports Range e.g 80 or 1-1024 or 80,22,23 (default: 1-32000)
   --timeout value      TCP Timeout in Millisecond (default: "500")
   --version, -v        print the version (default: false)
```
