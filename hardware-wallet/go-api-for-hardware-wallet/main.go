package main

import (
    "fmt"
    "time"

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
    return chunks
}

func MessageWipeDevice() [][64]byte {
    wipeDevice := &messages.WipeDevice{}
    data, err := proto.Marshal(wipeDevice)
    if err != nil {
        fmt.Printf(err.Error())
    }
    chunks := hardwareWallet.MakeTrezorMessage(data, messages.MessageType_MessageType_WipeDevice)
    return chunks
}

func MessageButtonAck() [][64]byte{
    buttonAck := &messages.ButtonAck{}
    data, _ := proto.Marshal(buttonAck)
    chunks := hardwareWallet.MakeTrezorMessage(data, messages.MessageType_MessageType_ButtonAck)
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

func MessageSetMnemonic() [][64]byte {
    setMnemonicMessage := &messages.SetMnemonic{
        Mnemonic:    proto.String("cloud flower upset remain green metal below cup stem infant art thank"),
    }

    data, _ := proto.Marshal(setMnemonicMessage)

    chunks := hardwareWallet.MakeTrezorMessage(data, messages.MessageType_MessageType_SetMnemonic)
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

func Initialize(dev usb.Device) {
    var msg wire.Message
    var chunks [][64]byte

    chunks = MessageInitialize()
    msg = SendToDevice(dev, chunks)
    initMsg := &messages.Initialize{}
    proto.Unmarshal(msg.Data, initMsg)
    fmt.Printf("Init success Answer is: %s\n", initMsg.State)
}

func WipeDevice(dev usb.Device) {
    var msg wire.Message
    var chunks [][64]byte
    var err error
    
    Initialize(dev)

    chunks = MessageWipeDevice()
    msg = SendToDevice(dev, chunks)
    fmt.Printf("buttonRequest %d! Answer is: %x\n", msg.Kind, msg.Data)

    chunks = MessageButtonAck()
    SendToDeviceNoAnswer(dev, chunks)

    _, err = msg.ReadFrom(dev)
    time.Sleep(3*time.Second)
	if err != nil {
        fmt.Printf(err.Error())
		return
    }
    fmt.Printf("MessageButtonAck Answer is: %d / %s\n", msg.Kind, msg.Data)

    Initialize(dev)
}

func LoadDevice(dev usb.Device) {
    var msg wire.Message
    var chunks [][64]byte
    var err error

    Initialize(dev)

    chunks = MessageLoadDevice()
    msg = SendToDevice(dev, chunks)
    fmt.Printf("LoadDevice %d! Answer is: %s\n", msg.Kind, msg.Data[2:])

    chunks = MessageButtonAck()
    SendToDeviceNoAnswer(dev, chunks)

    _, err = msg.ReadFrom(dev)
    time.Sleep(3*time.Second)
	if err != nil {
        fmt.Printf(err.Error())
		return
    }
    fmt.Printf("MessageButtonAck Answer is: %d / %s\n", msg.Kind, msg.Data)

    Initialize(dev)
}

func SetMnemonic(dev usb.Device) {

    var msg wire.Message
    var chunks [][64]byte
    var err error

    chunks = MessageSetMnemonic()
    msg = SendToDevice(dev, chunks)
    fmt.Printf("SetMnemonic %d! Answer is: %s\n", msg.Kind, msg.Data[2:])

    chunks = MessageButtonAck()
    SendToDeviceNoAnswer(dev, chunks)

    _, err = msg.ReadFrom(dev)
    time.Sleep(1*time.Second)
	if err != nil {
        fmt.Printf(err.Error())
		return
    }
    fmt.Printf("MessageButtonAck Answer is: %d / %s\n", msg.Kind, msg.Data)
}

func main() {
    dev, _ := hardwareWallet.GetTrezorDevice()
    var msg wire.Message
    var chunks [][64]byte

    // WipeDevice(dev)

    // LoadDevice(dev)

    SetMnemonic(dev)

    chunks = MessageSkycoinAddress()
    msg = SendToDevice(dev, chunks)
    fmt.Printf("Success %d! Answer is: %s\n", msg.Kind, msg.Data[2:])

    chunks = MessageSkycoinSignMessage()
    msg = SendToDevice(dev, chunks)
    fmt.Printf("Success %d! Answer is: %s\n", msg.Kind, msg.Data[2:])

    chunks = MessageCheckMessageSignature()
    msg = SendToDevice(dev, chunks)
    fmt.Printf("Success %d! Answer is: %s\n", msg.Kind, msg.Data[2:])

    Initialize(dev)
}
