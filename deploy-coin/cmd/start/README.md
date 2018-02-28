# Start new coin node with predefined configuration

starts-coin launches new coin node from JSON configuraion created by create-coin tool.
If path to file is not provided start-coin reads JSON configuration from standart input. 
If node is started as master it distributes inital coin volume among specified addresses.

## Installation

### Prerequisites

* Have go1.9+ installed
* Have `GOPATH` env set

### Install

```bash
$ go get github.com/skycoin/skycoin
$ go get github.com/services/deploy-coin
$ cd $GOPATH/services/deploy-coin/
$ make start-coin
```

## Run

```bash
$ start-coin -h
Usage of start-coin:
  -config string
        path to JSON configuration file for coin
  -gui
        lanuch web GUI for node in browser
  -guiPort int
        override guiPort from config
  -master
        run node as master and distribute initial coin volume
  -port int
        override port from config
  -rpc
        run web RPC service
  -rpcPort int
        override rpcPort from config
```


