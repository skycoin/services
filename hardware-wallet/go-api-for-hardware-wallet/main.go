package main

import(
    "fmt"
    "log"
    deviceWallet "./device-wallet"
	messages "./device-wallet/protob"
)


func main() {
	var deviceType deviceWallet.DeviceType
	if deviceWallet.DeviceConnected(deviceWallet.DeviceTypeEmulator) {
		deviceType = deviceWallet.DeviceTypeEmulator
	} else if deviceWallet.DeviceConnected(deviceWallet.DeviceTypeUsb) {
		deviceType = deviceWallet.DeviceTypeUsb
	} else {
		return
    }

    deviceWallet.WipeDevice(deviceType)
    var pinEnc string
    kind, _ := deviceWallet.DeviceChangePin(deviceType)
    for kind == uint16(messages.MessageType_MessageType_PinMatrixRequest) {
        log.Printf("PinMatrixRequest response: ")
        fmt.Scanln(&pinEnc)
        kind, _ = deviceWallet.DevicePinMatrixAck(deviceType, pinEnc)
    }

    // come on one-more time
    // testing what happen when we try to change an existing pin code
    kind, _ = deviceWallet.DeviceChangePin(deviceType)
    for kind == uint16(messages.MessageType_MessageType_PinMatrixRequest) {
        log.Printf("PinMatrixRequest response: ")
        fmt.Scanln(&pinEnc)
        kind, _ = deviceWallet.DevicePinMatrixAck(deviceType, pinEnc)
    }
}