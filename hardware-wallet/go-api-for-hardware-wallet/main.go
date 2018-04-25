package main

import (
    "fmt"

    "./hardware-wallet"

    messages "./protob"
    "./wire"
    "./usb"
    "github.com/golang/protobuf/proto"
)

func MessageSkycoinAddress() [][64]byte {
    skycoinAddress := &messages.SkycoinAddress{
        AddressN:    proto.Uint32(1),
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
        Signature: proto.String("Bk7jnoMj6W6Zd46AFSqKn5KFfdENKK5nx9qEqHdViWwz6n8RVRXVWnsdPMX5BCze5Lq1HerKTgKHPnzToL3XpHyuh"),
    }

    data, _ := proto.Marshal(skycoinCheckMessageSignature)

    chunks := hardwareWallet.MakeTrezorMessage(data, messages.MessageType_MessageType_SkycoinCheckMessageSignature)
    return chunks
}


func MessageSkycoinSignMessage() [][64]byte {
    skycoinSignMessage := &messages.SkycoinSignMessage{
        SecretKey:   proto.String("Qaj1vWfVPGUvX9dgmTWMRCzqUMcnxzT2M11K5yDMsc"),
        Message:     proto.String("Hello World!"),
    }

    data, _ := proto.Marshal(skycoinSignMessage)

    chunks := hardwareWallet.MakeTrezorMessage(data, messages.MessageType_MessageType_SkycoinSignMessage)
    return chunks
}

func SendToDevice(dev usb.Device, chunks [][64]byte) wire.Message {
    for _, element := range chunks {
        _, _ = dev.Write(element[:])
    }

    var msg wire.Message
    msg.ReadFrom(dev)
    return msg
}

func main() {
    dev, _ := hardwareWallet.GetTrezorDevice()
    var msg wire.Message
    var chunks [][64]byte

    chunks = MessageSkycoinAddress()
    msg = SendToDevice(dev, chunks)
    fmt.Printf("Success %d! Answer is: %s\n", msg.Kind, msg.Data[2:])

    // chunks = MessageSkycoinSignMessage()
    // msg = SendToDevice(dev, chunks)
    // fmt.Printf("Success %d! Answer is: %s\n", msg.Kind, msg.Data[2:])

    // chunks = MessageCheckMessageSignature()
    // msg = SendToDevice(dev, chunks)
    // fmt.Printf("Success %d! Answer is: %s\n", msg.Kind, msg.Data[2:])
}
