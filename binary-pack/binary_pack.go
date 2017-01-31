package binary_pack

import (
	"strings"
	"strconv"
	"errors"
	"encoding/binary"
	"bytes"
	"fmt"
)

type BinaryPack struct {}

func (bp *BinaryPack) Pack(format []string, msg []interface{}) ([]byte, error) {
	if len(format) > len(msg) {
		return nil, errors.New("Format is longer than values to pack")
	}

	res := []byte{}

	for i, f := range format {
		switch f {
		case "?":
			res = append(res, boolToBytes(msg[i].(bool))...)
		case "h", "H":
			res = append(res, intToBytes(msg[i].(int), 2)...)
		case "i", "I", "l", "L":
			res = append(res, intToBytes(msg[i].(int), 4)...)
		case "q", "Q":
			res = append(res, intToBytes(msg[i].(int), 8)...)
		case "f":
			res = append(res, float32ToBytes(msg[i].(float32), 4)...)
		case "d":
			res = append(res, float64ToBytes(msg[i].(float64), 8)...)
		default:
			if strings.Contains(f, "s") {
				n, _ := strconv.Atoi(strings.TrimRight(f, "s"))
				res = append(res, []byte(fmt.Sprintf("%s%s",
					msg[i].(string), strings.Repeat("\x00", n - len(msg[i].(string)))))...)
			} else {
				return nil, errors.New("Unexpected format token: '" + f + "'")
			}
		}
	}

	return res, nil
}

func (bp *BinaryPack) UnPack(format []string, msg []byte) ([]interface{}, error) {
	expected_size, err := bp.CalcSize(format)

	if err != nil {
		return nil, err
	}

	if expected_size > len(msg) {
		return nil, errors.New("Expected size is bigger than actual size of message")
	}

	res := []interface{}{}

	for _, f := range format {
		switch f {
		case "?":
			res = append(res, bytesToBool(msg[:1]))
			msg = msg[1:]
		case "h", "H":
			res = append(res, bytesToInt(msg[:2]))
			msg = msg[2:]
		case "i", "I", "l", "L":
			res = append(res, bytesToInt(msg[:4]))
			msg = msg[4:]
		case "q", "Q":
			res = append(res, bytesToInt(msg[:8]))
			msg = msg[8:]
		case "f":
			res = append(res, bytesToFloat32(msg[:4]))
			msg = msg[4:]
		case "d":
			res = append(res, bytesToFloat64(msg[:8]))
			msg = msg[8:]
		default:
			if strings.Contains(f, "s") {
				n, _ := strconv.Atoi(strings.TrimRight(f, "s"))
				res = append(res, string(msg[:n]))
				msg = msg[n:]
			} else {
				return nil, errors.New("Unexpected format token: '" + f + "'")
			}
		}
	}

	return res, nil
}

func (bp *BinaryPack) CalcSize(format []string) (int, error) {
	var size int

	for _, f := range format {
		switch f {
		case "?":
			size = size + 1
		case "h", "H":
			size = size + 2
		case "i", "I", "l", "L", "f":
			size = size + 4
		case "q", "Q", "d":
			size = size + 8
		default:
			if strings.Contains(f, "s") {
				n, _ := strconv.Atoi(strings.TrimRight(f, "s"))
				size = size + n
			} else {
				return 0, errors.New("Unexpected format token: '" + f + "'")
			}
		}
	}

	return size, nil
}

func boolToBytes(x bool) []byte {
	if x {
		return intToBytes(1, 1)
	}
	return intToBytes(0, 1)
}

func bytesToBool(b []byte) bool {
	return bytesToInt(b) > 0
}

func intToBytes(n int, size int) []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, int64(n))
	return buf.Bytes()[buf.Len()-size:]
}

func bytesToInt(b []byte) int {
	buf := bytes.NewBuffer(b)

	switch len(b) {
	case 1:
		var x int8
		binary.Read(buf, binary.BigEndian, &x)
		return int(x)
	case 2:
		var x int16
		binary.Read(buf, binary.BigEndian, &x)
		return int(x)
	case 4:
		var x int32
		binary.Read(buf, binary.BigEndian, &x)
		return int(x)
	default:
		var x int64
		binary.Read(buf, binary.BigEndian, &x)
		return int(x)
	}
}

func float32ToBytes(n float32, size int) []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, n)
	return buf.Bytes()[buf.Len()-size:]
}

func bytesToFloat32(b []byte) float32 {
	var x float32
	buf := bytes.NewBuffer(b)
	binary.Read(buf, binary.BigEndian, &x)
	return x
}

func float64ToBytes(n float64, size int) []byte {
	buf := bytes.NewBuffer([]byte{})
	binary.Write(buf, binary.BigEndian, n)
	return buf.Bytes()[buf.Len()-size:]
}

func bytesToFloat64(b []byte) float64 {
	var x float64
	buf := bytes.NewBuffer(b)
	binary.Read(buf, binary.BigEndian, &x)
	return x
}
