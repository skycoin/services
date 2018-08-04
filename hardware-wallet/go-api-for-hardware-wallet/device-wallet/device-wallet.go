package deviceWallet

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"log"
	"net"
	"time"

	"./usb"
	"./wire"

	proto "github.com/golang/protobuf/proto"
	messages "./protob"
)

// DeviceType type of device: emulated or usb
type DeviceType int32

const (
	// DeviceTypeEmulator use emulator
	DeviceTypeEmulator DeviceType = 1
	// DeviceTypeUsb use usb
	DeviceTypeUsb DeviceType = 2
)

func getEmulatorDevice() (net.Conn, error) {
	return net.Dial("udp", "127.0.0.1:21324")
}

func getUsbDevice() (usb.Device, error) {
	w, err := usb.InitWebUSB()
	if err != nil {
		log.Panicf("webusb: %s", err)
		return nil, err
	}
	h, err := usb.InitHIDAPI()
	if err != nil {
		log.Panicf("hidapi: %s", err)
		return nil, err
	}
	b := usb.Init(w, h)

	var infos []usb.Info
	infos, err = b.Enumerate()
	if len(infos) <= 0 {
		return nil, err
	}
	tries := 0
	dev, err := b.Connect(infos[0].Path)
	if err != nil {
		log.Panicf(err.Error())
		if tries < 3 {
			tries++
			time.Sleep(100 * time.Millisecond)
		} else {
			return nil, err
		}
	}
	return dev, err
}

func sendToDeviceNoAnswer(dev io.ReadWriteCloser, chunks [][64]byte) error {
	for _, element := range chunks {
		_, err := dev.Write(element[:])
		if err != nil {
			return err
		}
	}
	return nil
}
func sendToDevice(dev io.ReadWriteCloser, chunks [][64]byte) (wire.Message, error) {
	var msg wire.Message
	for _, element := range chunks {
		_, err := dev.Write(element[:])
		if err != nil {
			return msg, err
		}
	}
	_, err := msg.ReadFrom(dev)
	return msg, err
}

func makeTrezorHeader(data []byte, msgID messages.MessageType) []byte {
	header := new(bytes.Buffer)
	binary.Write(header, binary.BigEndian, []byte("?##"))
	binary.Write(header, binary.BigEndian, uint16(msgID))
	binary.Write(header, binary.BigEndian, uint32(len(data)))
	binary.Write(header, binary.BigEndian, []byte("\n"))
	return header.Bytes()
}

func makeTrezorMessage(data []byte, msgID messages.MessageType) [][64]byte {
	message := new(bytes.Buffer)
	binary.Write(message, binary.BigEndian, []byte("##"))
	binary.Write(message, binary.BigEndian, uint16(msgID))
	binary.Write(message, binary.BigEndian, uint32(len(data)))
	binary.Write(message, binary.BigEndian, []byte("\n"))
	if len(data) > 0 {
		binary.Write(message, binary.BigEndian, data[1:])
	}

	messageLen := message.Len()
	var chunks [][64]byte
	i := 0
	for messageLen > 0 {
		var chunk [64]byte
		chunk[0] = '?'
		copy(chunk[1:], message.Bytes()[63*i:63*(i+1)])
		chunks = append(chunks, chunk)
		messageLen -= 63
		i = i + 1
	}
	return chunks
}

func getDevice(deviceType DeviceType) (io.ReadWriteCloser, error) {
	var dev io.ReadWriteCloser
	var err error
	switch deviceType {
	case DeviceTypeEmulator:
		dev, err = getEmulatorDevice()
		break
	case DeviceTypeUsb:
		dev, err = getUsbDevice()
		break
	}
	if (dev == nil && err == nil) {
		err = errors.New("No device connected")
	}
	return dev, err
}

// DeviceCheckMessageSignature Check a message signature matches the given address.
func DeviceCheckMessageSignature(deviceType DeviceType, message string, signature string, address string) (uint16, []byte) {

	dev, err := getDevice(deviceType)
	if err != nil {
		log.Panicf(err.Error())
		return 0, make([]byte, 0)
	}
	defer dev.Close()

	// Send CheckMessageSignature

	skycoinCheckMessageSignature := &messages.SkycoinCheckMessageSignature{
		Address:   proto.String(address),
		Message:   proto.String(message),
		Signature: proto.String(signature),
	}

	data, _ := proto.Marshal(skycoinCheckMessageSignature)
	chunks := makeTrezorMessage(data, messages.MessageType_MessageType_SkycoinCheckMessageSignature)
	msg, err := sendToDevice(dev, chunks)
	if err != nil {
		log.Printf(err.Error())
		return msg.Kind, msg.Data
	}
	log.Printf("Success %d! address that issued the signature is: %s\n", msg.Kind, msg.Data)
	return msg.Kind, msg.Data
}

// MessageButtonAck send this message (before user action) when the device expects the user to push a button
func MessageButtonAck() [][64]byte {
	buttonAck := &messages.ButtonAck{}
	data, _ := proto.Marshal(buttonAck)
	chunks := makeTrezorMessage(data, messages.MessageType_MessageType_ButtonAck)
	return chunks
}

// DeviceSetMnemonic Configure the device with a mnemonic.
func DeviceSetMnemonic(deviceType DeviceType, mnemonic string) {

	dev, err := getDevice(deviceType)
	if err != nil {
		log.Panicf(err.Error())
		return
	}
	defer dev.Close()

	// Send SetMnemonic

	skycoinSetMnemonic := &messages.SetMnemonic{
		Mnemonic: proto.String(mnemonic),
	}

	data, _ := proto.Marshal(skycoinSetMnemonic)
	chunks := makeTrezorMessage(data, messages.MessageType_MessageType_SetMnemonic)

	msg, err := sendToDevice(dev, chunks)
	if err != nil {
		log.Panicf(err.Error())
		return
	}

	log.Printf("Success %d! Mnemonic %s\n", msg.Kind, msg.Data)

	// Send ButtonAck
	chunks = MessageButtonAck()
	err = sendToDeviceNoAnswer(dev, chunks)
	if err != nil {
		log.Panicf(err.Error())
		return
	}

	time.Sleep(1 * time.Second)
	_, err = msg.ReadFrom(dev)
	if err != nil {
		log.Panicf(err.Error())
		return
	}

	log.Printf("MessageButtonAck Answer is: %d / %s\n", msg.Kind, msg.Data)
}

// DecodeFailMsg convert byte data into string containing the failure returned by the device
func DecodeFailMsg(kind uint16, data []byte) (uint16, string) {
	if kind == uint16(messages.MessageType_MessageType_Failure) {
		failure := &messages.Failure{}
		err := proto.Unmarshal(data, failure)
		if err != nil {
			log.Panicf("unmarshaling error: %s\n", err.Error())
			return kind, ""
		}
		return kind, failure.GetMessage()
	}
	log.Panic("Calling DecodeFailMsg with wrong message type")
	return kind, ""
}

// DecodeResponseSkycoinAddress convert byte data into list of addresses, meant to be used after DevicePinMatrixAck
func DecodeResponseSkycoinAddress(kind uint16, data []byte) (uint16, []string) {
	if kind == uint16(messages.MessageType_MessageType_ResponseSkycoinAddress) {
		responseSkycoinAddress := &messages.ResponseSkycoinAddress{}
		err := proto.Unmarshal(data, responseSkycoinAddress)
		if err != nil {
			log.Panicf("unmarshaling error: %s\n", err.Error())
			return kind, make([]string, 0)
		}
		return kind, responseSkycoinAddress.GetAddresses()
	}
	log.Panic("Calling DecodeResponseSkycoinAddress with wrong message type")
	return kind, make([]string, 0)
}

// DecodeResponseSkycoinSignMessage convert byte data into signed message, meant to be used after DevicePinMatrixAck
func DecodeResponseSkycoinSignMessage(kind uint16, data []byte) (uint16, string) {
	if kind == uint16(messages.MessageType_MessageType_ResponseSkycoinSignMessage) {
		responseSkycoinSignMessage := &messages.ResponseSkycoinSignMessage{}
		err := proto.Unmarshal(data, responseSkycoinSignMessage)
		if err != nil {
			log.Panicf("unmarshaling error: %s\n", err.Error())
			return kind, ""
		}
		return kind, responseSkycoinSignMessage.GetSignedMessage()
	}
	log.Panic("Calling DecodeResponseeSkycoinSignMessage with wrong message type")
	return kind, ""
}

// DeviceAddressGen Ask the device to generate an address
func DeviceAddressGen(deviceType DeviceType, addressN int, startIndex int) (uint16, []string) {

	dev, err := getDevice(deviceType)
	if err != nil {
		log.Panicf(err.Error())
		return 0, make([]string, 0)
	}
	defer dev.Close()
	skycoinAddress := &messages.SkycoinAddress{
		AddressN:   proto.Uint32(uint32(addressN)),
		StartIndex: proto.Uint32(uint32(startIndex)),
	}
	data, _ := proto.Marshal(skycoinAddress)

	chunks := makeTrezorMessage(data, messages.MessageType_MessageType_SkycoinAddress)

	msg, err := sendToDevice(dev, chunks)
	if err != nil {
		log.Panicf("sendToDevice error: %s\n", err.Error())
	}
	if msg.Kind == uint16(messages.MessageType_MessageType_ResponseSkycoinAddress) {
		return DecodeResponseSkycoinAddress(msg.Kind, msg.Data)
	} else if msg.Kind == uint16(messages.MessageType_MessageType_PinMatrixRequest) {
		log.Println("This operation requires a PIN code")
		return msg.Kind, make([]string, 0)
	}
	failureMsg := &messages.Failure{}
	err = proto.Unmarshal(msg.Data, failureMsg)
	if err != nil {
		log.Panicf("unmarshaling error: %s\n", err.Error())
	}
	log.Printf("Failure %d! Answer is: %s\n", failureMsg.GetCode(), failureMsg.GetMessage())
	return msg.Kind, make([]string, 0)
}

// DeviceSignMessage Ask the device to sign a message using the secret key at given index.
func DeviceSignMessage(deviceType DeviceType, addressN int, message string) (uint16, string) {

	dev, err := getDevice(deviceType)
	if err != nil {
		log.Panicf(err.Error())
		return 0, ""
	}
	defer dev.Close()

	skycoinSignMessage := &messages.SkycoinSignMessage{
		AddressN: proto.Uint32(uint32(addressN)),
		Message:  proto.String(message),
	}

	data, _ := proto.Marshal(skycoinSignMessage)

	chunks := makeTrezorMessage(data, messages.MessageType_MessageType_SkycoinSignMessage)

	msg, err := sendToDevice(dev, chunks)
	if err != nil {
		log.Panicf(err.Error())
	}

	if msg.Kind == uint16(messages.MessageType_MessageType_ResponseSkycoinSignMessage) {
		return DecodeResponseSkycoinSignMessage(msg.Kind, msg.Data)
	} else if msg.Kind == uint16(messages.MessageType_MessageType_PinMatrixRequest) {
		log.Println("This operation requires a PIN code")
		return msg.Kind, ""
	}
	log.Printf("DeviceSignMessage answer kind %d", msg.Kind)
	failureMsg := &messages.Failure{}
	err = proto.Unmarshal(msg.Data, failureMsg)
	if err != nil {
		log.Panicf("unmarshaling error: %s\n", err.Error())
	}
	log.Printf("Failure %d! Answer is: %s\n", failureMsg.GetCode(), failureMsg.GetMessage())
	return msg.Kind, ""
}

// DeviceConnected check if a device is connected
func DeviceConnected(deviceType DeviceType) bool {
	dev, err := getDevice(deviceType)
	if dev == nil {
		return false
	}
	defer dev.Close()
	if err != nil {
		return false
	}
	msgRaw := &messages.Ping{}
	data, err := proto.Marshal(msgRaw)
	chunks := makeTrezorMessage(data, messages.MessageType_MessageType_Ping)
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

// Initialize send an init request to the device
func initialize(dev io.ReadWriteCloser) {
	var chunks [][64]byte

	initialize := &messages.Initialize{}
	data, _ := proto.Marshal(initialize)
	chunks = makeTrezorMessage(data, messages.MessageType_MessageType_Initialize)
	_, err := sendToDevice(dev, chunks)
	if err != nil {
		log.Panicf(err.Error())
		return
	}
}

// WipeDevice wipes out device configuration
func WipeDevice(deviceType DeviceType) {
	dev, err := getDevice(deviceType)
	if err != nil {
		log.Panicf(err.Error())
		return
	}
	defer dev.Close()
	var msg wire.Message
	var chunks [][64]byte

	initialize(dev)

	wipeDevice := &messages.WipeDevice{}
	data, _ := proto.Marshal(wipeDevice)
	chunks = makeTrezorMessage(data, messages.MessageType_MessageType_WipeDevice)
	msg, err = sendToDevice(dev, chunks)
	if err != nil {
		log.Panicf(err.Error())
		return
	}
	log.Printf("Wipe device %d! Answer is: %x\n", msg.Kind, msg.Data)

	// Send ButtonAck
	chunks = MessageButtonAck()
	err = sendToDeviceNoAnswer(dev, chunks)
	if err != nil {
		log.Panicf(err.Error())
		return
	}

	_, err = msg.ReadFrom(dev)
	time.Sleep(1 * time.Second)
	if err != nil {
		log.Panicf(err.Error())
		return
	}
	log.Printf("MessageButtonAck Answer is: %d / %s\n", msg.Kind, msg.Data)

	initialize(dev)
}


// DeviceChangePin changes device's PIN code
// The message that is sent contains an encoded form of the PIN.
// The digits of the PIN are displayed in a 3x3 matrix on the Trezor,
// and the message that is sent back is a string containing the positions
// of the digits on that matrix. Below is the mapping between positions
// and characters to be sent:
// 7 8 9
// 4 5 6
// 1 2 3
// For example, if the numbers are laid out in this way on the Trezor,
// 3 1 5
// 7 8 4
// 9 6 2
// To set the PIN "12345", the positions are:
// top, bottom-right, top-left, right, top-right
// so you must send "83769".
func DeviceChangePin(deviceType DeviceType) (uint16, []byte) {
	dev, err := getDevice(deviceType)
	if err != nil {
		log.Panicf(err.Error())
		return 0, make([]byte, 0)
	}
	defer dev.Close()

    changePin := &messages.ChangePin{}
    data, _ := proto.Marshal(changePin)
	chunks := makeTrezorMessage(data, messages.MessageType_MessageType_ChangePin)
	msg, err := sendToDevice(dev, chunks)
	if err != nil {
		log.Panicf(err.Error())
		return msg.Kind, msg.Data
	}
    // Acknowledge that a button has been pressed
    if msg.Kind == uint16(messages.MessageType_MessageType_ButtonRequest) {
		chunks = MessageButtonAck()
		err = sendToDeviceNoAnswer(dev, chunks)
		if err != nil {
			log.Panicf(err.Error())
			return msg.Kind, msg.Data
		}
	
		_, err = msg.ReadFrom(dev)
		time.Sleep(1 * time.Second)
		log.Printf("MessageButtonAck Answer is: %d / %s\n", msg.Kind, msg.Data)
	}
	return msg.Kind, msg.Data
}

// DevicePinMatrixAck during PIN code setting use this message to send user input to device
func DevicePinMatrixAck(deviceType DeviceType, p string) (uint16, []byte) {
	time.Sleep(1 * time.Second)
	dev, err := getDevice(deviceType)
	if err != nil {
		log.Panicf(err.Error())
		return 0, make([]byte, 0)
	}
	defer dev.Close()
	var msg wire.Message
	log.Printf("Setting pin: %s\n", p)
    pinAck := &messages.PinMatrixAck{
        Pin: proto.String(p),
    }
    data, _ := proto.Marshal(pinAck)

    chunks := makeTrezorMessage(data, messages.MessageType_MessageType_PinMatrixAck)
	msg, err = sendToDevice(dev, chunks)
	if err != nil {
		log.Panicf(err.Error())
		return msg.Kind, msg.Data
	}
	log.Printf("MessagePinMatrixAck Answer is: %d / %s\n", msg.Kind, msg.Data)
	return msg.Kind, msg.Data
}