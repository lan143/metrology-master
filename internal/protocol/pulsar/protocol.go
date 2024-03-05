package pulsar

import (
	"context"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"io"
	"math/rand"
	"time"
)

const (
	commError        byte = 0x00
	commReadChannels      = 0x01
	commReadParam         = 0x0A
)

const (
	errCodeUnknownFunction    byte = 0x01
	errCodeBitMask                 = 0x02
	errCodeRequestLength           = 0x03
	errCodeParam                   = 0x04
	errCodeNeedAuth                = 0x05
	errCodeWriteParamRange         = 0x06
	errCodeUnknownArchiveType      = 0x07
	errCodeMaxArchive              = 0x08
)

const (
	paramUID     uint16 = 0x0000
	paramAddress        = 0x0001
	paramVersion        = 0x0002
)

var (
	ErrDeviceNotResponding = errors.New("the device is not responding")
)

type Version struct {
	SWVersion string
	HWVersion string
}

type Pulsar struct {
	port      io.ReadWriteCloser
	responses map[uint16]chan frame

	log *zap.Logger
}

func NewPulsar(port io.ReadWriteCloser, log *zap.Logger) *Pulsar {
	s := &Pulsar{
		port: port,
		log:  log,
	}
	s.responses = make(map[uint16]chan frame)
	go s.processResponse()

	return s
}

func (s *Pulsar) ReadChannels(ctx context.Context, address [4]byte, mask uint32) ([]uint32, error) {
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	respChan, err := s.sendRequest(
		address,
		commReadChannels,
		[]byte{byte(mask), byte(mask >> 8), byte(mask >> 16), byte(mask >> 24)},
	)
	if err != nil {
		return nil, err
	}

	select {
	case <-ctx.Done():
		return nil, ErrDeviceNotResponding
	case resp := <-respChan:
		if resp.FN == commError {
			return nil, s.mapErrorCode(resp.Payload[0])
		}

		payload := make([]uint32, len(resp.Payload)/4)
		for i := 0; i < len(resp.Payload)/4; i++ {
			payload[i] = (uint32(resp.Payload[i+3]) << 24) | (uint32(resp.Payload[i+2]) << 16) |
				(uint32(resp.Payload[i+1]) << 8) | uint32(resp.Payload[i])
		}

		s.log.Debug(
			"ReadChannels",
			zap.Any("uint32", payload),
		)

		return payload, nil
	}
}

func (s *Pulsar) ReadParam(ctx context.Context, address [4]byte, index uint16) ([]byte, error) {
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	respChan, err := s.sendRequest(
		address,
		commReadParam,
		[]byte{byte(index), byte(index >> 8)},
	)
	if err != nil {
		return nil, err
	}

	select {
	case <-ctx.Done():
		return nil, ErrDeviceNotResponding
	case resp := <-respChan:
		if resp.FN == commError {
			return nil, s.mapErrorCode(resp.Payload[0])
		}

		s.log.Debug(
			"ReadParam",
			zap.Any("Data", resp.Payload),
		)

		return resp.Payload, nil
	}
}

func (s *Pulsar) GetVersion(ctx context.Context, address [4]byte) (Version, error) {
	payload, err := s.ReadParam(ctx, address, paramVersion)
	if err != nil {
		return Version{}, err
	}

	version := Version{
		SWVersion: fmt.Sprintf("%d.%d.%d.%d", payload[3], payload[2], payload[7], payload[6]),
		HWVersion: fmt.Sprintf("%d.%d.%d-%d", payload[5], payload[4], payload[1], payload[0]),
	}

	s.log.Debug(
		"GetVersion",
		zap.Any("version", version),
	)

	return version, nil
}

func (s *Pulsar) sendRequest(address [4]byte, command byte, payload []byte) (chan frame, error) {
	req := frame{
		Address: address,
		FN:      command,
		Len:     0,
		Payload: payload,
		ID:      uint16(rand.Uint32()),
		CRC:     0,
	}

	s.responses[req.ID] = make(chan frame, 1)
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

func (s *Pulsar) processResponse() {
	buff := make([]byte, 255)
	resp := frame{}

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
				resp = frame{}
			}
			continue
		}

		s.log.Debug("receive response", zap.Any("response", resp))

		ch, ok := s.responses[resp.ID]
		if !ok {
			s.log.Error(
				"process response",
				zap.Error(fmt.Errorf("not found response chan for ID: 0x%X", resp.ID)),
			)
			continue
		}

		ch <- resp
		close(ch)
		delete(s.responses, resp.ID)

		resp = frame{}
	}
}

func (s *Pulsar) mapErrorCode(code byte) error {
	switch code {
	case errCodeUnknownFunction:
		return errors.New("the requested function code is unknown")
	case errCodeBitMask:
		return errors.New("error in request bitmask")
	case errCodeRequestLength:
		return errors.New("invalid request length")
	case errCodeParam:
		return errors.New("missing parameter")
	case errCodeNeedAuth:
		return errors.New("the entry is blocked, authorization is required")
	case errCodeWriteParamRange:
		return errors.New("the parameter written is outside the specified range")
	case errCodeUnknownArchiveType:
		return errors.New("the requested archive type is missing")
	case errCodeMaxArchive:
		return errors.New("exceeding the maximum number of archived values per package")
	default:
		return fmt.Errorf("unknown error, code: %d", code)
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
