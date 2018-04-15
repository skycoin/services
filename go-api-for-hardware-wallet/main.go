package main

import (
    "fmt"

    "./hardware-wallet"

    messages "./protob"
    "./wire"
    "github.com/golang/protobuf/proto"
)

func MessageSkycoinAddress() [][64]byte {
    skycoinAddress := &messages.SkycoinAddress{
        Seed:        proto.String("seed"),
        AddressType: messages.SkycoinAddressType_AddressTypeSkycoin.Enum(),
    }
    data, _ := proto.Marshal(skycoinAddress)

    chunks := hardwareWallet.MakeTrezorMessage(data, messages.MessageType_MessageType_SkycoinAddress)
    return chunks
}


func MessageCheckMessageSignature() [][64]byte {
    skycoinCheckMessageSignature := &messages.SkycoinCheckMessageSignature{
		Address:   proto.String("2EVNa4CK9SKosT4j1GEn8SuuUUEAXaHAMbM"),
		Message:   proto.String("Hello World!"),
		Signature: proto.String("GA82nXSwVEPV5soMjCiQkJb4oLEAo6FMK8CAE2n2YBTm7xjhAknUxtZrhs3RPVMfQsEoLwkJCEgvGj8a2vzthBQ1M"),
	}

    data, _ := proto.Marshal(skycoinCheckMessageSignature)

    chunks := hardwareWallet.MakeTrezorMessage(data, messages.MessageType_MessageType_SkycoinCheckMessageSignature)
    return chunks
}

func main() {
    dev, _ := hardwareWallet.GetTrezorDevice()

    chunks := MessageCheckMessageSignature()

    for _, element := range chunks {
        _, _ = dev.Write(element[:])
    }

    var msg wire.Message
    msg.ReadFrom(dev)

    fmt.Printf("Success %d! Address is: %s\n", msg.Kind, msg.Data)
}
