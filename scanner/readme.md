# Scanner for bitcoin wallet

This utility allow you to search transactions of bitcoin addresses by scanning blockchain. It works throw china firewall in contrast to api blockchain.info


## Setup project

### Prerequisites

* Have go1.8+ installed
* Have `GOPATH` env set
* Have btcd started


### Installing btcd

- Run the following commands to obtain btcd, all dependencies, and install it:

```bash
$ go get -u github.com/Masterminds/glide
$ git clone https://github.com/btcsuite/btcd $GOPATH/src/github.com/btcsuite/btcd
$ cd $GOPATH/src/github.com/btcsuite/btcd
$ glide install
$ go install . ./cmd/...
```
- btcd (and utilities) will now be installed in $GOPATH/bin. Go where and

```bash
$ ./btcd
```

### Start scanning

```bash
$ cd services/scanner/
$ go run main.go
```


It will open brower with web interface where you can add addresses to wallet and scan blocks.


