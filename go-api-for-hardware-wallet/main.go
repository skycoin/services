package main

import (
    "fmt"
    "./hardware-wallet"

    messages "./protob"
    "./wire"
    "github.com/golang/protobuf/proto"
)

func main() {
    dev, _ := hardwareWallet.GetTrezorDevice()

    skycoinAddress := &messages.SkycoinAddress{
        Seed:        proto.String("seed"),
        AddressType: messages.SkycoinAddressType_AddressTypeSkycoin.Enum(),
    }
    data, _ := proto.Marshal(skycoinAddress)

    chunks := hardwareWallet.MakeTrezorMessage(data, messages.MessageType_MessageType_SkycoinAddress)

    for _, element := range chunks {
        _, _ = dev.Write(element[:])
    }

    var msg wire.Message
    msg.ReadFrom(dev)

    fmt.Printf("Success %d! Address is: %s\n", msg.Kind, msg.Data) 
}