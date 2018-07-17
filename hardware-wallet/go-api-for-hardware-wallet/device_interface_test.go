package deviceWallet

import (
	"testing"

	messages "./protob"
	"github.com/stretchr/testify/require"
)

func TestMain(t *testing.T) { 

	var deviceType DeviceType
	if DeviceConnected(DeviceTypeEmulator) {
		deviceType = DeviceTypeEmulator
	} else if DeviceConnected(DeviceTypeUsb) {
		deviceType = DeviceTypeUsb
	} else {
		t.Skip("TestMain do not work if nor Emulator and Usb device is connected")
		return
	}

    // var msg wire.Message
    // var chunks [][64]byte
    // var inputWord string
    // var pinEnc string
    // var err error

    
    WipeDevice(deviceType)

    // Send ChangePin message
    // chunks = MessageChangePin()
    // msg = SendToDevice(dev, chunks)
    // // Acknowledge that a button has been pressed
    // if msg.Kind == uint16(messages.MessageType_MessageType_ButtonRequest) {
    //     chunks = MessageButtonAck()
    //     msg = SendToDevice(dev, chunks)
    // }
    /*
        The message that is sent contains an encoded form of the PIN.
        The digits of the PIN are displayed in a 3x3 matrix on the Trezor,
        and the message that is sent back is a string containing the positions
        of the digits on that matrix. Below is the mapping between positions
        and characters to be sent:
        
        7 8 9
        4 5 6
        1 2 3
        
        For example, if the numbers are laid out in this way on the Trezor,
        
        3 1 5
        7 8 4
        9 6 2
        
        To set the PIN "12345", the positions are:
        
        top, bottom-right, top-left, right, top-right
        
        so you must send "83769".
    */
    // for msg.Kind == uint16(messages.MessageType_MessageType_PinMatrixRequest) {
    //     fmt.Print("PinMatrixRequest response: ")
    //     fmt.Scanln(&pinEnc)
    //     chunks = MessagePinMatrixAck(pinEnc)
    //     msg = SendToDevice(dev, chunks)
    // }
    // fmt.Printf("Response: %s\n", messages.MessageType_name[int32(msg.Kind)])
    // if msg.Kind == uint16(messages.MessageType_MessageType_Failure) {
    //     failMsg := &messages.Failure{}
    //     proto.Unmarshal(msg.Data, failMsg)
    //     fmt.Printf("Code: %d\nMessage: %s\n", failMsg.GetCode(), failMsg.GetMessage());
    // }
    
    // chunks = MessageRecoveryDevice(12)
    // msg = SendToDevice(dev, chunks)
    // if msg.Kind == uint16(messages.MessageType_MessageType_ButtonRequest) {
    //     chunks = MessageButtonAck()
    //     msg = SendToDevice(dev, chunks)
    // }
    // for msg.Kind == uint16(messages.MessageType_MessageType_WordRequest) {
    //     fmt.Print("Word request: ")
    //     fmt.Scanln(&inputWord)
    //     chunks = MessageWordAck(strings.TrimSpace(inputWord))
    //     msg = SendToDevice(dev, chunks)
    // }
    // fmt.Printf("Response: %s\n", messages.MessageType_name[int32(msg.Kind)])
    // if msg.Kind == uint16(messages.MessageType_MessageType_Failure) {
    //     failMsg := &messages.Failure{}
    //     proto.Unmarshal(msg.Data, failMsg)
    //     fmt.Printf("Code: %d\nMessage: %s\n", failMsg.GetCode(), failMsg.GetMessage());
    // }


    DeviceSetMnemonic(deviceType, "cloud flower upset remain green metal below cup stem infant art thank")

	kind, addresses := DeviceAddressGen(deviceType, 9, 15)
	logger.Info(addresses)
	require.Equal(t, kind, uint16(messages.MessageType_MessageType_ResponseSkycoinAddress))
	i := 0
	require.Equal(t, 9, len(addresses))
	require.Equal(t, addresses[i], "3NpgZ6g1UWZc5f5B7gC3hU6NhyEWxznohG")
	i++
	require.Equal(t, addresses[i], "Wr6wE5bHwBpg4kTs3EF4xi2cLs2dEWy1BN")
	i++
	require.Equal(t, addresses[i], "2DpKC15mSBhNMptvLgudRim6ScY4df1TwLd")
	i++
	require.Equal(t, addresses[i], "ZdaQWbWers3qYpKKSoBNq237CXQhGmHwX9")
	i++
	require.Equal(t, addresses[i], "9mTMfX1v6TnCYCK8frzSKAL4m2Lx1uu7Kq")
	i++
	require.Equal(t, addresses[i], "2cKu9tZz3eGqo6ny7D447o4RpMFNEk8KyXr")
	i++
	require.Equal(t, addresses[i], "2mqM8j7Zqq5MiWLEgJyAzTAPQ9sd575nh9X")
	i++
	require.Equal(t, addresses[i], "29pYKsirWo21ZPhEsdNmcCVExgAeK5ShpMF")
	i++
	require.Equal(t, addresses[i], "n6ou5D4hSGCXsAiVCJX6y6jc454xvcoSet")
    // chunks = MessageBackupDevice()
    // msg = SendToDevice(dev, chunks)
    // for msg.Kind == uint16(messages.MessageType_MessageType_ButtonRequest) {
    //     chunks = MessageButtonAck()
    //     msg = SendToDevice(dev, chunks)
    // }
	// fmt.Printf("Success %d! Answer is: %s\n", msg.Kind, msg.Data[2:])
	

	kind, addresses = DeviceAddressGen(deviceType, 1, 1)
	logger.Info(addresses)
	require.Equal(t, uint16(messages.MessageType_MessageType_ResponseSkycoinAddress), kind)
	require.Equal(t, len(addresses), 1)
	require.Equal(t, addresses[0], "zC8GAQGQBfwk7vtTxVoRG7iMperHNuyYPs")

	message:= "Hello World!"
	kind, signature := DeviceSignMessage(deviceType, 1, message)
	logger.Info(string(signature[1:]))
	require.Equal(t, uint16(messages.MessageType_MessageType_Success), kind) //Success message
	require.Equal(t, 89, len(signature[2:]))
    // chunks = MessageSkycoinSignMessage()
    // msg = SendToDevice(dev, chunks)
    // fmt.Printf("Success %d! Answer is: %s\n", msg.Kind, msg.Data[2:])

	kind, data := DeviceCheckMessageSignature(deviceType, message, string(signature[2:]), addresses[0])
	require.Equal(t, uint16(messages.MessageType_MessageType_Success), kind) //Success message
	require.Equal(t, "zC8GAQGQBfwk7vtTxVoRG7iMperHNuyYPs", string(data[2:]))
    // chunks = MessageCheckMessageSignature()
    // msg = SendToDevice(dev, chunks)
    // fmt.Printf("Success %d! Answer is: %s\n", msg.Kind, msg.Data[2:])
}
func TestGetAddressUsb(t *testing.T) {
	if DeviceConnected(DeviceTypeUsb) == false {
		t.Skip("TestGetAddressUsb do not work if Usb device is not connected")
		return
	}

	require.True(t, DeviceConnected(DeviceTypeUsb))
	WipeDevice(DeviceTypeUsb)
	// need to connect the usb device
	DeviceSetMnemonic(DeviceTypeUsb, "cloud flower upset remain green metal below cup stem infant art thank")
	kind, address := DeviceAddressGen(DeviceTypeUsb, 2, 0)
	logger.Info(address)
	require.Equal(t, kind, uint16(messages.MessageType_MessageType_ResponseSkycoinAddress)) //Success message
	require.Equal(t, address[0], "2EU3JbveHdkxW6z5tdhbbB2kRAWvXC2pLzw")
	require.Equal(t, address[1], "zC8GAQGQBfwk7vtTxVoRG7iMperHNuyYPs")
}

func TestGetAddressEmulator(t *testing.T) {
	if DeviceConnected(DeviceTypeEmulator) == false {
		t.Skip("TestGetAddressEmulator do not work if Emulator device is not running")
		return
	}

	require.True(t, DeviceConnected(DeviceTypeEmulator))
	WipeDevice(DeviceTypeEmulator)
	DeviceSetMnemonic(DeviceTypeEmulator, "cloud flower upset remain green metal below cup stem infant art thank")
	kind, address := DeviceAddressGen(DeviceTypeEmulator, 2, 0)
	logger.Info(address)
	require.Equal(t, kind, uint16(messages.MessageType_MessageType_ResponseSkycoinAddress)) //Success message
	require.Equal(t, address[0], "2EU3JbveHdkxW6z5tdhbbB2kRAWvXC2pLzw")
	require.Equal(t, address[1], "zC8GAQGQBfwk7vtTxVoRG7iMperHNuyYPs")
}
