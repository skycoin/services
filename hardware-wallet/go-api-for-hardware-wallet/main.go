package main

import (
    "fmt"
    "time"
    "strings"
    wallet "./hardware-wallet"

    messages "./protob"
    "./wire"
    "github.com/golang/protobuf/proto"
)

func MessageInitialize() [][64]byte {
    initialize := &messages.Initialize{}
    data, _ := proto.Marshal(initialize)

    chunks := wallet.MakeTrezorMessage(data, messages.MessageType_MessageType_Initialize)
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
    chunks := emulatorWallet.MakeTrezorMessage(data, messages.MessageType_MessageType_ResetDevice)
    return chunks
}

func MessageWipeDevice() [][64]byte {
    wipeDevice := &messages.WipeDevice{}
    data, err := proto.Marshal(wipeDevice)
    if err != nil {
        fmt.Printf(err.Error())
    }
    chunks := wallet.MakeTrezorMessage(data, messages.MessageType_MessageType_WipeDevice)
    return chunks
}

func MessageButtonAck() [][64]byte{
    buttonAck := &messages.ButtonAck{}
    data, _ := proto.Marshal(buttonAck)
    chunks := wallet.MakeTrezorMessage(data, messages.MessageType_MessageType_ButtonAck)
    return chunks
}

func MessageLoadDevice() [][64]byte {
    loadDevice := &messages.LoadDevice{
        Mnemonic:    proto.String("cloud flower upset remain green metal below cup stem infant art thank"),
    }
    data, _ := proto.Marshal(loadDevice)

    chunks := wallet.MakeTrezorMessage(data, messages.MessageType_MessageType_LoadDevice)
    return chunks
}

func MessageSetMnemonic() [][64]byte {
    setMnemonicMessage := &messages.SetMnemonic{
        Mnemonic:    proto.String("cloud flower upset remain green metal below cup stem infant art thank"),
    }

    data, _ := proto.Marshal(setMnemonicMessage)

    chunks := wallet.MakeTrezorMessage(data, messages.MessageType_MessageType_SetMnemonic)
    return chunks
}

func MessageSkycoinAddress() [][64]byte {
    skycoinAddress := &messages.SkycoinAddress{
        AddressN:    proto.Uint32(1),
        AddressType: messages.SkycoinAddressType_AddressTypeSkycoin.Enum(),
    }
    data, _ := proto.Marshal(skycoinAddress)

    chunks := wallet.MakeTrezorMessage(data, messages.MessageType_MessageType_SkycoinAddress)
    return chunks
}


func MessageCheckMessageSignature() [][64]byte {
    skycoinCheckMessageSignature := &messages.SkycoinCheckMessageSignature{
        Address:   proto.String("2EVNa4CK9SKosT4j1GEn8SuuUUEAXaHAMbM"),
        Message:   proto.String("Hello World!"),
        Signature: proto.String("Bk7jnoMj6W6Zd46AFSqKn5KFfdENKK5nx9qEqHdViWwz6n8RVRXVWnsdPMX5BCze5Lq1HerKTgKHPnzToL3XpHyuh"),
    }

    data, _ := proto.Marshal(skycoinCheckMessageSignature)

    chunks := wallet.MakeTrezorMessage(data, messages.MessageType_MessageType_SkycoinCheckMessageSignature)
    return chunks
}

func MessageSkycoinSignMessage() [][64]byte {
    skycoinSignMessage := &messages.SkycoinSignMessage{
        AddressN:    proto.Uint32(1),
        // SecretKey:   proto.String("Qaj1vWfVPGUvX9dgmTWMRCzqUMcnxzT2M11K5yDMsc"),
        Message:     proto.String("Hello World!"),
    }

    data, _ := proto.Marshal(skycoinSignMessage)

    chunks := wallet.MakeTrezorMessage(data, messages.MessageType_MessageType_SkycoinSignMessage)
    return chunks
}

func SendToDeviceNoAnswer(dev wallet.TrezorDevice, chunks [][64]byte) {
    for _, element := range chunks {
        _, _ = dev.Write(element[:])
    }
}


func SendToDevice(dev wallet.TrezorDevice, chunks [][64]byte) wire.Message {
    for _, element := range chunks {
        _, _ = dev.Write(element[:])
    }

    var msg wire.Message
    msg.ReadFrom(dev)
    return msg
}

func Initialize(dev wallet.TrezorDevice) {
    var msg wire.Message
    var chunks [][64]byte

    chunks = MessageInitialize()
    msg = SendToDevice(dev, chunks)
    initMsg := &messages.Initialize{}
    proto.Unmarshal(msg.Data, initMsg)
    fmt.Printf("Init success Answer is: %s\n", initMsg.State)
}

func WipeDevice(dev wallet.TrezorDevice) {
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

func LoadDevice(dev wallet.TrezorDevice) {
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

func SetMnemonic(dev wallet.TrezorDevice) {

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

func MessageRecoveryDevice(words uint32) [][64]byte {
    msg := &messages.RecoveryDevice{
        WordCount: proto.Uint32(words),
        Type: proto.Uint32(0),
    }
    data, _ := proto.Marshal(msg)

    chunks := wallet.MakeTrezorMessage(data, messages.MessageType_MessageType_RecoveryDevice)
    return chunks
}

func MessageWordAck(word string) [][64]byte {
    msg := &messages.WordAck{
        Word: proto.String(word),
    }
    data, _ := proto.Marshal(msg)

    chunks := wallet.MakeTrezorMessage(data, messages.MessageType_MessageType_WordAck)
    return chunks
} 

func DeviceConnected(dev wallet.TrezorDevice) bool {
    if dev == nil {
        return false
    }
    msgRaw := &messages.Ping{}
    data, err := proto.Marshal(msgRaw)
    chunks := wallet.MakeTrezorMessage(data, messages.MessageType_MessageType_Ping)
    for _, element := range chunks {
        _, err = dev.Write(element[:])
        if err != nil {
            return false
        }
    }
    var msg wire.Message
    _, err = msg.ReadFrom(dev)
    if err != nil {
        return false
    }
    return msg.Kind == uint16(messages.MessageType_MessageType_Success)
}

func main() {
    dev, _ := wallet.GetTrezorDevice()
    var msg wire.Message
    var chunks [][64]byte
    var inputWord string
    
    if DeviceConnected(dev) {
        fmt.Printf("Connected\n")
    } else {
        fmt.Printf("Not Connected\n")
        return
    }
    
    WipeDevice(dev)
    
    chunks = MessageRecoveryDevice(12)
    msg = SendToDevice(dev, chunks)
    if msg.Kind == uint16(messages.MessageType_MessageType_ButtonRequest) {
        chunks = MessageButtonAck()
        msg = SendToDevice(dev, chunks)
    }
    for msg.Kind == uint16(messages.MessageType_MessageType_WordRequest) {
        fmt.Print("Word request: ")
        fmt.Scanln(&inputWord)
        chunks = MessageWordAck(strings.TrimSpace(inputWord))
        msg = SendToDevice(dev, chunks)
    }
    fmt.Printf("Response: %s\n", messages.MessageType_name[int32(msg.Kind)])
    if msg.Kind == uint16(messages.MessageType_MessageType_Failure) {
        failMsg := &messages.Failure{}
        proto.Unmarshal(msg.Data, failMsg)
        fmt.Printf("Code: %d\nMessage: %s\n", failMsg.GetCode(), failMsg.GetMessage());
    }


    chunks = MessageSkycoinAddress()
    msg = SendToDevice(dev, chunks)
    responseSkycoinAddress := &messages.ResponseSkycoinAddress{}
    err = proto.Unmarshal(msg.Data, responseSkycoinAddress)
    if err != nil {
        fmt.Printf("unmarshaling error: %s\n", err.Error())
    }

    fmt.Printf("Success %d! Answer is: %s\n", msg.Kind, responseSkycoinAddress.String())
}
