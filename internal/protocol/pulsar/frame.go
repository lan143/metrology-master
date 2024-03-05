package pulsar

import (
	"errors"
	"fmt"
)

var (
	errNeedMoreBytes = errors.New("need more bytes")
)

type frame struct {
	Address [4]byte
	FN      byte
	Len     byte
	Payload []byte
	ID      uint16
	CRC     uint16

	buffer []byte
}

func (rr *frame) parseBytes(bytes []byte) error {
	var i int

	rr.buffer = append(rr.buffer, bytes...)

	if len(rr.buffer) < 10 {
		return errors.New("frame length less that 10")
	}

	for i = 0; i < 4; i++ {
		rr.Address[i] = rr.buffer[i]
	}

	rr.FN = rr.buffer[i]
	i++

	rr.Len = rr.buffer[i]
	i++

	if int(rr.Len) > len(rr.buffer) {
		return errNeedMoreBytes
	}

	if rr.Len > 10 {
		j := 0
		rr.Payload = make([]byte, rr.Len-4)

		for i < int(rr.Len)-4 {
			rr.Payload[j] = rr.buffer[i]
			i++
			j++
		}
	}

	rr.ID = uint16(rr.buffer[i])
	i++
	rr.ID |= uint16(rr.buffer[i]) << 8
	i++

	crc := calculateCrc(rr.buffer[:i])

	rr.CRC = uint16(rr.buffer[i])
	i++
	rr.CRC |= uint16(rr.buffer[i]) << 8
	i++

	if crc != rr.CRC {
		return fmt.Errorf("invalid crc. Excepted: 0x%X. Actual: 0x%X", rr.CRC, crc)
	}

	return nil
}

func (rr *frame) generateBytes() []byte {
	rr.Len = byte(10 + len(rr.Payload))

	var (
		i     int
		bytes = make([]byte, rr.Len)
	)

	for i = 0; i < 4; i++ {
		bytes[i] = rr.Address[i]
	}

	bytes[i] = rr.FN
	i++
	bytes[i] = byte(10 + len(rr.Payload))
	i++

	for j := 0; j < len(rr.Payload); j++ {
		bytes[i] = rr.Payload[j]
		i++
	}

	bytes[i] = byte(rr.ID)
	i++
	bytes[i] = byte(rr.ID >> 8)
	i++

	rr.CRC = calculateCrc(bytes[:rr.Len-2])

	bytes[i] = byte(rr.CRC)
	i++
	bytes[i] = byte(rr.CRC >> 8)
	i++

	return bytes
}
