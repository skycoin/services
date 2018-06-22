package hardwareWallet

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"time"
    "errors"
	"../usb"

	messages "../protob"
)
type TrezorDevice = usb.Device

func GetTrezorDevice() (TrezorDevice, error) {
	w, err := usb.InitWebUSB()
	if err != nil {
		return nil, err
	}
	h, err := usb.InitHIDAPI()
	if err != nil {
		return nil, err
	}
	b := usb.Init(w, h)

	var infos []usb.Info
	infos, _ = b.Enumerate()
    if len(infos) < 1 {
        return nil, errors.New("No USB devices connected.")
    }
	tries := 0
	dev, err := b.Connect(infos[0].Path)
	if err != nil {
		fmt.Printf(err.Error())
		if tries < 3 {
			tries++
			time.Sleep(100 * time.Millisecond)
		}
	}
	return dev, err
}
func MakeTrezorHeader(data []byte, msgID messages.MessageType) []byte {
	header := new(bytes.Buffer)
	binary.Write(header, binary.BigEndian, []byte("?##"))
	binary.Write(header, binary.BigEndian, uint16(msgID))
	binary.Write(header, binary.BigEndian, uint32(len(data)))
	binary.Write(header, binary.BigEndian, []byte("\n"))
	return header.Bytes()
}

func MakeTrezorMessage(data []byte, msgID messages.MessageType) [][64]byte {
	message := new(bytes.Buffer)
	binary.Write(message, binary.BigEndian, []byte("##"))
	binary.Write(message, binary.BigEndian, uint16(msgID))
	binary.Write(message, binary.BigEndian, uint32(len(data)))
	binary.Write(message, binary.BigEndian, []byte("\n"))
	if (len(data) >= 1){
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
