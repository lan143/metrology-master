package pulsar_t1

import (
	"context"
	"fmt"
	"github.com/lan143/metrology-master/internal/meter"
	pulsar_m "github.com/lan143/metrology-master/internal/protocol/pulsar-m"
)

const (
	Type string = "pulsar_1t"
)

const (
	channelT1 int = 0
	channelT2     = 12
)

type (
	Config struct {
		Address [4]byte
	}

	pulsarT1 struct {
		config  Config
		service *pulsar_m.PulsarM
	}
)

func (m *pulsarT1) GetParams() meter.Params {
	uid := fmt.Sprintf("0x%X%X%X%X", m.config.Address[0], m.config.Address[1], m.config.Address[2], m.config.Address[3])

	return meter.Params{
		UID:          uid,
		StateTopic:   fmt.Sprintf("power-meter/%s/state", uid),
		Manufacturer: "Teplovodohran",
		Model:        "Pulsar 1T",
		Name:         "Teplovodohran Pulsar 1T",
		HWVersion:    "1.0.0",
		SWVersion:    "1.0.0",
		Flags:        meter.FlagHasPowerConsumption,
	}
}

func NewPulsarT1(config Config, service *pulsar_m.PulsarM) meter.ElectricMeter {
	return &pulsarT1{
		config:  config,
		service: service,
	}
}

func (m *pulsarT1) GetPowerConsumption(ctx context.Context) (float64, error) {
	resp, err := m.service.ReadChannels(
		ctx,
		m.config.Address,
		uint32(0xFFFF),
	)
	if err != nil {
		return 0, err
	}

	var t1, t2 float64
	t1 = float64(resp[channelT1]) / 100
	t2 = float64(resp[channelT2]) / 100

	return t1 + t2, nil
}
