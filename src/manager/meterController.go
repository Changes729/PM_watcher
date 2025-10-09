package manager

import (
	"encoding/hex"
	"errors"
	"log"
)

type DLT_645_2007 struct {
	ID          []byte
	ControlCode uint8
	DataLength  uint8
	Data        []byte
}

/**
 * - CC for Command Code
 * - RC for Receive Code
 */
const (
	EMPTY_ADDRESS = "AAAAAAAAAAAA"

	CC_GET_IP   uint8 = 0x13
	CC_GET_DATA uint8 = 0x11

	RC_GET_DATA   uint8 = 0x91
	RC_GET_DATA_N uint8 = 0xB1
	RC_GET_DATA_E uint8 = 0xD1
)

const (
	_START_BYTE          uint8 = 0x68
	_END_BYTE            uint8 = 0x16
	_SENDER_PROCESS_BYTE uint8 = 0x33
)

var (
	FC_COMBINED_ENERGY      = []byte{0x00, 0x00, 0x00, 0x00}
	FC_NEUTRAL_LINE_CURRENT = []byte{0x02, 0x80, 0x00, 0x01}
)

func ReverseArray(arr []byte) []byte {
	left := 0
	right := len(arr) - 1
	for left < right {
		arr[left], arr[right] = arr[right], arr[left]
		left++
		right--
	}

	return arr
}

func GenCommand(address string, control_code uint8, data []byte) (cmd []byte) {
	if len(address) != 12 {
		log.Fatal("Address length error")
		return
	}

	WATTMETER_ID, err := hex.DecodeString(address)
	if err != nil {
		log.Fatal("Address hex decode error")
		return
	}
	WATTMETER_ID = ReverseArray(WATTMETER_ID)

	data_length := len(data)
	cmd = append(cmd, _START_BYTE)
	cmd = append(cmd, WATTMETER_ID...)
	cmd = append(cmd, _START_BYTE)
	cmd = append(cmd, control_code)
	cmd = append(cmd, byte(data_length))
	for _, d := range data {
		cmd = append(cmd, d+_SENDER_PROCESS_BYTE)
	}

	check_sum := byte(0)
	for _, d := range cmd {
		check_sum += d
	}
	cmd = append(cmd, check_sum&0xFF)
	cmd = append(cmd, byte(_END_BYTE))

	return
}

func _MarshalPackage(bytes []byte) (data DLT_645_2007, err error) {
	log.Printf("Marshal packages: %X", bytes)
	WEAKUP_BYTE := byte(0xFE)
	err = errors.New("No valid package found")

	for i, b := range bytes {
		if b != WEAKUP_BYTE {
			if bytes[i] == byte(_START_BYTE) && bytes[len(bytes)-1] == byte(_END_BYTE) {
				sum := 0
				for _, d := range bytes[i : len(bytes)-2] {
					sum += int(d)
				}

				if (sum & 0xFF) != int(bytes[len(bytes)-2]) {
					err = errors.New("Checksum error")
					break
				}

				data = DLT_645_2007{
					ID:          bytes[i+1 : i+1+6],
					ControlCode: bytes[i+8],
					DataLength:  bytes[i+9],
					Data:        bytes[i+10 : i+10+int(bytes[i+9])],
				}

				for i, b := range data.Data {
					data.Data[i] = b - _SENDER_PROCESS_BYTE
				}
				err = nil
			} else {
				err = errors.New("Package not end")
			}
			break
		}
	}

	return
}

func MarshalDeviceAddress(bytes []byte) (address string, err error) {
	address = EMPTY_ADDRESS
	data, err := _MarshalPackage(bytes)
	if err != nil {
		return
	}

	if data.ControlCode != 0x93 {
		err = errors.New("Control code error")
		return
	}

	if data.DataLength != 6 {
		err = errors.New("Data length error")
		return
	}

	for i := 0; i < int(data.DataLength); i++ {
		if data.Data[i] != data.ID[i] {
			err = errors.New("Data mismatch error")
			return
		}
	}

	address = hex.EncodeToString(ReverseArray(data.ID))
	return
}

func MarshalData(bytes []byte) (data DLT_645_2007, err error) {
	data, err = _MarshalPackage(bytes)
	return
}

func CheckPackIntegrity(bytes []byte) (integrity bool, beginIndex int) {
	integrity = false
	beginIndex = 0

	for i, b := range bytes {
		if b == byte(_START_BYTE) {
			integrity =
				len(bytes) >= i+9 &&
					len(bytes) >= (i+9+2+int(bytes[i+9])) &&
					bytes[i+7] == byte(_START_BYTE) &&
					bytes[i+9+2+int(bytes[i+9])] == byte(_END_BYTE)
		}

		if integrity {
			break
		}
	}

	log.Printf("Receive: %X, Integrity: %v", bytes, integrity)
	return
}
