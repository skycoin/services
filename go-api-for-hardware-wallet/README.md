# Before use

    # generate protobuf files
    protoc -I ./protob  --go_out=./protob protob/messages.proto protob/types.proto 
    go build
    go run main.go #code example in main.go
