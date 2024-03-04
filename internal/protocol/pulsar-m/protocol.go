package pulsar_m

import (
	"context"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"io"
	"math/rand"
)

const (
	CommError        byte = 0x00
	CommReadChannels      = 0x01
)

var (
	errNeedMoreBytes = errors.New("need more bytes")
)

type requestResponse struct {
	Address [4]byte
	FN      byte
	Len     byte
	Payload []byte
	ID      uint16
	CRC     uint16

	buffer []byte
}

func (rr *requestResponse) parseBytes(bytes []byte) error {
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

func (rr *requestResponse) generateBytes() []byte {
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

type PulsarM struct {
	port      io.ReadWriteCloser
	responses map[uint16]chan requestResponse

	log *zap.Logger
}

func NewPulsarM(port io.ReadWriteCloser, log *zap.Logger) *PulsarM {
	s := &PulsarM{
		port: port,
		log:  log,
	}
	s.responses = make(map[uint16]chan requestResponse)
	go s.processResponse()

	return s
}

func (s *PulsarM) ReadChannels(ctx context.Context, address [4]byte, mask uint32) ([]uint32, error) {
	respChan, err := s.sendRequest(
		address,
		CommReadChannels,
		[]byte{byte(mask), byte(mask >> 8), byte(mask >> 16), byte(mask >> 24)},
	)
	if err != nil {
		return nil, err
	}

	select {
	case <-ctx.Done():
		return nil, errors.New("device not respond")
	case resp := <-respChan:
		if resp.FN == CommError {
			return nil, errors.New("receive error command")
		}

		payload := make([]uint32, len(resp.Payload)/4)
		for i := 0; i < len(resp.Payload)/4; i++ {
			payload[i] = (uint32(resp.Payload[i+3]) << 24) | (uint32(resp.Payload[i+2]) << 16) |
				(uint32(resp.Payload[i+1]) << 8) | uint32(resp.Payload[i])
		}

		return payload, nil
	}
}

func (s *PulsarM) sendRequest(address [4]byte, command byte, payload []byte) (chan requestResponse, error) {
	req := requestResponse{
		Address: address,
		FN:      command,
		Len:     0,
		Payload: payload,
		ID:      uint16(rand.Uint32()),
		CRC:     0,
	}

	s.responses[req.ID] = make(chan requestResponse, 1)
	bytes := req.generateBytes()

	s.log.Debug(
		"send request",
		zap.Any("request", req),
	)
	s.log.Debug(
		"send bytes",
		zap.Binary("bytes", bytes),
	)

	_, err := s.port.Write(bytes)
	if err != nil {
		return nil, err
	}

	return s.responses[req.ID], nil
}

func (s *PulsarM) processResponse() {
	buff := make([]byte, 255)
	resp := requestResponse{}

	for {
		n, err := s.port.Read(buff)
		if err != nil {
			s.log.Error("read port", zap.Error(err))
		}

		if n == 0 {
			s.log.Debug("read port - eof")
			break
		}

		s.log.Debug("receive bytes", zap.Binary("bytes", buff[:n]))

		err = resp.parseBytes(buff[:n])
		if err != nil {
			if err != errNeedMoreBytes {
				s.log.Error("parse bytes", zap.Error(err))
				resp = requestResponse{}
			}
			continue
		}

		s.log.Debug("receive response", zap.Any("response", resp))

		ch, ok := s.responses[resp.ID]
		if !ok {
			s.log.Error("process response", zap.Error(fmt.Errorf("not found response chan for ID: 0x%X", resp.ID)))
			continue
		}

		ch <- resp
		close(ch)
		delete(s.responses, resp.ID)

		resp = requestResponse{}
	}
}

func calculateCrc(data []byte) uint16 {
	var (
		i, j   int
		result uint16 = 0xFFFF
	)

	for i = 0; i < len(data); i++ {
		result ^= uint16(data[i])

		for j = 0; j < 8; j++ {
			if result&0x1 > 0 {
				result = (result >> 1) ^ 0xA001
			} else {
				result = result >> 1
			}
		}
	}

	return result
}
