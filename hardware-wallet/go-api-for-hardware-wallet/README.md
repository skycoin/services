# Before use

# Install golang:

    https://github.com/golang/go/wiki/Ubuntu

# Install google protobuf:

    sudo apt-get install protobuf-compiler python-protobuf golang-goprotobuf-dev
    go get -u github.com/golang/protobuf/protoc-gen-go

# Compile the protobuf project dependencies

    make -C vendor/nanopb/generator/proto/
    make -C protob/


The you can generate protobuf files (only once each time the messages change)

    protoc -I../tiny-firmware/vendor/nanopb/generator/proto/ -I ./protob  --go_out=./protob protob/messages.proto protob/types.proto 

# Usage

    go run main.go #code example in main.go
