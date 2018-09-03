# Create configuration of new coin

create-coin creates JSON configuration for new coin.
Created JSON configuraion is used by start-coin tool to start coin node 

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
$ create-coin -h
Usage of create-coin:
  -addr int
        number of distribution addresses (default 100)
  -code string
        code of new coin (default "SKY")
  -file string
        file to save configuration of new coin
  -peers int
        number of trusted peers running on localhost (default 3)
  -vol int
        coin volume to send to each of disribution addresses (default 1000000)
```


