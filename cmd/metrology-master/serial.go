package main

import (
	"errors"
	"github.com/lan143/metrology-master/pkg/serial"
	serial2 "go.bug.st/serial"
	"io"
)

type Serial struct {
	ports map[string]io.ReadWriteCloser
}

func (c *Command) InitSerial(configs map[string]*serial.Config) error {
	c.serial.ports = make(map[string]io.ReadWriteCloser, len(configs))

	for name, config := range configs {
		var stopBits serial2.StopBits
		switch config.StopBits {
		case 1:
			stopBits = serial2.OneStopBit
		case 2:
			stopBits = serial2.OnePointFiveStopBits
		case 3:
			stopBits = serial2.TwoStopBits
		default:
			return errors.New("unsupported stop bits setting")
		}

		var parity serial2.Parity
		switch config.Parity {
		case 0:
			parity = serial2.NoParity
		case 1:
			parity = serial2.OddParity
		case 2:
			parity = serial2.EvenParity
		case 3:
			parity = serial2.MarkParity
		case 4:
			parity = serial2.SpaceParity
		default:
			return errors.New("unsupported parity setting")
		}

		mode := &serial2.Mode{
			BaudRate: config.BaudRate,
			DataBits: config.DataBits,
			StopBits: stopBits,
			Parity:   parity,
		}

		port, err := serial2.Open(config.Port, mode)
		if err != nil {
			return err
		}

		c.serial.ports[name] = port
	}

	return nil
}
