# Go lang api for skycoin hardware wallet

## Installation

### Install golang

    https://github.com/golang/go/wiki/Ubuntu

### Install google protobuf

    sudo apt-get install protobuf-compiler python-protobuf golang-goprotobuf-dev
    go get -u github.com/golang/protobuf/proto/proto
    go get -u github.com/stretchr/testify/require

## Compile the protobuf project dependencies

    make -C vendor/nanopb/generator/proto/
    make -C protob/

## Usage

### Generate protobuf files

Only once each time the messages change:

    cd device-wallet/ 
    protoc -I../../tiny-firmware/vendor/nanopb/generator/proto/ -I ./protob  --go_out=./protob protob/messages.proto protob/types.proto

### Run

    go test -run TestMain
