# Create configuration of new coin

create-coin creates JSON configuration for new coin.
created JSON configuraion is used by start-coin tool to start coin node 

## Installation

### Prerequisites

* Have go1.9+ installed
* Have `GOPATH` env set

### Install

```bash
$ go get github.com/skycoin/skycoin
$ go get github.com/services/deploy-coin
$ cd $GOPATH/services/deploy-coin/
$ make create-coin
```

## Run

```bash
$ coin-create -h
  Usage of coin-create:
  -c string
        code of new coin (default "SKY")
  -f string
        file to save configuration of new coin
```


