package main

import (
    "fmt"
    // "time"

    "./hardware-wallet"

    messages "./protob"
    "./wire"
    "./usb"
    "github.com/golang/protobuf/proto"
)

func MessageInitialize() [][64]byte {
    initialize := &messages.Initialize{}
    data, _ := proto.Marshal(initialize)

    chunks := hardwareWallet.MakeTrezorMessage(data, messages.MessageType_MessageType_Initialize)
    fmt.Printf("chunks: %s\n",chunks)
    return chunks
}


func MessageResetDevice() [][64]byte {
    resetDevice := &messages.ResetDevice{
        Strength:    proto.Uint32(256),
        U2FCounter:    proto.Uint32(0),
        Language:   proto.String("english"),
        SkipBackup:     proto.Bool(false),
        PassphraseProtection:     proto.Bool(false),
        PinProtection:     proto.Bool(false),
        DisplayRandom:     proto.Bool(false),
    }
    data, _ := proto.Marshal(resetDevice)
    chunks := hardwareWallet.MakeTrezorMessage(data, messages.MessageType_MessageType_ResetDevice)
    return chunks
}

func MessageWipeDevice() [][64]byte {
    wipeDevice := &messages.WipeDevice{}
    data, err := proto.Marshal(wipeDevice)
    if err != nil {
        fmt.Printf(err.Error())
    }
    fmt.Printf("data: %x\n",data)
    chunks := hardwareWallet.MakeTrezorMessage(data, messages.MessageType_MessageType_WipeDevice)

    fmt.Printf("chunks: %s\n",chunks)
    return chunks
}

func MessageButtonAckWipeDevice() [][64]byte{
    buttonRequest := &messages.ButtonRequest{}
    data, _ := proto.Marshal(buttonRequest)
    chunks := hardwareWallet.MakeTrezorMessage(data, messages.MessageType_MessageType_ButtonRequest)
    return chunks
}

func MessageLoadDevice() [][64]byte {
    loadDevice := &messages.LoadDevice{
        Mnemonic:    proto.String("cloud flower upset remain green metal below cup stem infant art thank"),
    }
    data, _ := proto.Marshal(loadDevice)

    chunks := hardwareWallet.MakeTrezorMessage(data, messages.MessageType_MessageType_LoadDevice)
    return chunks
}

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
        AddressN:    proto.Uint32(1),
        // SecretKey:   proto.String("Qaj1vWfVPGUvX9dgmTWMRCzqUMcnxzT2M11K5yDMsc"),
        Message:     proto.String("Hello World!"),
    }

    data, _ := proto.Marshal(skycoinSignMessage)

    chunks := hardwareWallet.MakeTrezorMessage(data, messages.MessageType_MessageType_SkycoinSignMessage)
    return chunks
}

func SendToDeviceNoAnswer(dev usb.Device, chunks [][64]byte) {
    for _, element := range chunks {
        _, _ = dev.Write(element[:])
    }
}


func SendToDevice(dev usb.Device, chunks [][64]byte) wire.Message {
    for _, element := range chunks {
        _, _ = dev.Write(element[:])
    }

    var msg wire.Message
    msg.ReadFrom(dev)
    return msg
}

func WipeDevice(dev usb.Device) {
    var msg wire.Message
    var chunks [][64]byte
    var err error

    chunks = MessageInitialize()
    msg = SendToDevice(dev, chunks)
    fmt.Printf("Success %d! Answer is: %s\n", msg.Kind, msg.Data[2:])

    chunks = MessageWipeDevice()
    msg = SendToDevice(dev, chunks)
    fmt.Printf("Success %d! Answer is: %s\n", msg.Kind, msg.Data[2:])

    chunks = MessageButtonAckWipeDevice()
    SendToDeviceNoAnswer(dev, chunks)

    _, err = msg.ReadFrom(dev)
	if err != nil {
        fmt.Printf(err.Error())
		return
    }
    fmt.Printf("WipeDevice Answer is: %d / %s\n", msg.Kind, msg.Data)

    chunks = MessageInitialize()
    msg = SendToDevice(dev, chunks)
    fmt.Printf("Success %d! Answer is: %s\n", msg.Kind, msg.Data[2:])
}

func main() {
    dev, _ := hardwareWallet.GetTrezorDevice()
    var msg wire.Message
    var chunks [][64]byte

    WipeDevice(dev)

    // chunks = MessageResetDevice()
    // msg = SendToDevice(dev, chunks)
    // fmt.Printf("Success %d! Answer is: %s\n", msg.Kind, msg.Data[2:])

    // chunks = MessageLoadDevice()
    // msg = SendToDevice(dev, chunks)
    // fmt.Printf("Success %d! Answer is: %s\n", msg.Kind, msg.Data[2:])

    // chunks = MessageSkycoinAddress()
    // msg = SendToDevice(dev, chunks)
    // fmt.Printf("Success %d! Answer is: %s\n", msg.Kind, msg.Data[2:])

    // chunks = MessageSkycoinSignMessage()
    // msg = SendToDevice(dev, chunks)
    // fmt.Printf("Success %d! Answer is: %s\n", msg.Kind, msg.Data[2:])

    // chunks = MessageCheckMessageSignature()
    // msg = SendToDevice(dev, chunks)
    // fmt.Printf("Success %d! Answer is: %s\n", msg.Kind, msg.Data[2:])
}
