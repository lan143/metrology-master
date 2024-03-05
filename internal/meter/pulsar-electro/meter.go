package pulsar_electro

import (
	"context"
	"fmt"
	"github.com/lan143/metrology-master/internal/meter"
	pulsar_m "github.com/lan143/metrology-master/internal/protocol/pulsar"
	"go.uber.org/zap"
	"math"
	"strings"
)

const (
	Type         string = "pulsar_electro"
	manufacturer string = "Тепловодохран"
	model        string = "Пульсар"
)

const (
	channelT1 int = 0
	channelT2     = 12
)

const (
	paramFrequency uint16 = 0x100
	paramPhase            = 0x10A
	paramModel            = 0x016C
)

const (
	paramVoltageOffset uint16 = iota
	paramCurrentOffset
	paramActivePowerOffset
	paramReactivePowerOffset
	paramFullPowerOffset
	paramCoeffPower
	paramAngle
)

type (
	Config struct {
		Address [4]byte
	}

	pulsarT1 struct {
		config  Config
		service *pulsar_m.Pulsar
		log     *zap.Logger

		params meter.Params
	}
)

func NewPulsarElectro(config Config, service *pulsar_m.Pulsar, log *zap.Logger) meter.ElectricMeter {
	return &pulsarT1{
		config:  config,
		service: service,
		log:     log,
	}
}

func (m *pulsarT1) GetParams() meter.Params {
	return m.params
}

func (m *pulsarT1) Init(ctx context.Context) error {
	return m.buildParams(ctx)
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

	return math.Round((t1+t2)*100) / 100, nil
}

func (m *pulsarT1) GetFrequency(ctx context.Context) (float64, error) {
	resp, err := m.service.ReadParam(
		ctx,
		m.config.Address,
		paramFrequency,
	)
	if err != nil {
		return 0, err
	}

	var freq float64
	freq = float64(uint16(resp[0])|uint16(resp[1])<<8) / 100

	return math.Round(freq*100) / 100, nil
}

func (m *pulsarT1) GetVoltage(ctx context.Context) (float64, error) {
	resp, err := m.service.ReadParam(
		ctx,
		m.config.Address,
		paramPhase+paramVoltageOffset,
	)
	if err != nil {
		return 0, err
	}

	var voltage float64
	voltage = float64(uint16(resp[0])|uint16(resp[1])<<8) / 100

	return math.Round(voltage*100) / 100, nil
}

func (m *pulsarT1) GetCurrent(ctx context.Context) (float64, error) {
	resp, err := m.service.ReadParam(
		ctx,
		m.config.Address,
		paramPhase+paramCurrentOffset,
	)
	if err != nil {
		return 0, err
	}

	var current float64
	current = float64(uint32(resp[0])|uint32(resp[1])<<8|uint32(resp[2])<<16|uint32(resp[3])<<24) / 1000

	return math.Round(current*100) / 100, nil
}

func (m *pulsarT1) GetActivePower(ctx context.Context) (float64, error) {
	resp, err := m.service.ReadParam(
		ctx,
		m.config.Address,
		paramPhase+paramActivePowerOffset,
	)
	if err != nil {
		return 0, err
	}

	var power float64
	power = float64(uint16(resp[0]) | uint16(resp[1])<<8)

	return power, nil
}

func (m *pulsarT1) GetReactivePower(ctx context.Context) (float64, error) {
	resp, err := m.service.ReadParam(
		ctx,
		m.config.Address,
		paramPhase+paramReactivePowerOffset,
	)
	if err != nil {
		return 0, err
	}

	var power float64
	power = float64(uint16(resp[0]) | uint16(resp[1])<<8)

	return power, nil
}

func (m *pulsarT1) GetFullPower(ctx context.Context) (float64, error) {
	resp, err := m.service.ReadParam(
		ctx,
		m.config.Address,
		paramPhase+paramFullPowerOffset,
	)
	if err != nil {
		return 0, err
	}

	var power float64
	power = float64(uint16(resp[0]) | uint16(resp[1])<<8)

	return power, nil
}

func (m *pulsarT1) buildParams(ctx context.Context) error {
	uid := fmt.Sprintf("0x%X%X%X%X", m.config.Address[0], m.config.Address[1], m.config.Address[2], m.config.Address[3])

	version, err := m.service.GetVersion(ctx, m.config.Address)
	if err != nil {
		return err
	}

	data, err := m.service.ReadParam(ctx, m.config.Address, paramModel)
	if err != nil {
		return err
	}

	modelName, err := m.buildModelName(data)
	if err != nil {
		return err
	}

	m.params = meter.Params{
		UID:          uid,
		StateTopic:   fmt.Sprintf("power-meter/%s/state", uid),
		Manufacturer: manufacturer,
		Model:        modelName,
		Name:         manufacturer + " " + modelName,
		HWVersion:    version.HWVersion,
		SWVersion:    version.SWVersion,
		Flags: meter.FlagHasPowerConsumption | meter.FlagHasFrequency | meter.FlagHasVoltage | meter.FlagHasCurrent |
			meter.FlagHasActivePower | meter.FlagHasReactivePower | meter.FlagHasFullPower,
	}

	return nil
}

func (m *pulsarT1) buildModelName(data []byte) (string, error) {
	if len(data) < 8 {
		return "", fmt.Errorf("invalid data length: %d", len(data))
	}

	m.log.Debug(
		"build model name",
		zap.Any("data", data),
	)

	var modelName = []string{model}

	switch data[7] {
	case 0:
		modelName = append(modelName, "Т1")
	case 2:
		modelName = append(modelName, "Т1Т")
	default:
		m.log.Error(
			"build model name",
			zap.Error(fmt.Errorf("unsupported model id %d", data[7])),
		)
	}

	switch data[6] {
	case 0:
		modelName = append(modelName, "1А")
	default:
		m.log.Error(
			"build model name",
			zap.Error(fmt.Errorf("unsupported accuracy class %d", data[7])),
		)
	}

	switch data[5] {
	case 0:
		modelName = append(modelName, "5_60")
	case 1:
		modelName = append(modelName, "5_80")
	case 2:
		modelName = append(modelName, "10_80")
	case 3:
		modelName = append(modelName, "10_100")
	default:
		m.log.Error(
			"build model name",
			zap.Error(fmt.Errorf("unsupported current limits %d", data[7])),
		)
	}

	switch data[4] {
	case 0:
	case 1:
		modelName = append(modelName, "RS-485")
	case 2:
		modelName = append(modelName, "MBus")
	case 3:
		modelName = append(modelName, "IoT")
	case 4:
		modelName = append(modelName, "PLC")
	case 5:
		modelName = append(modelName, "OPTO")
	case 6:
		modelName = append(modelName, "GSM")
	default:
		m.log.Error(
			"build model name",
			zap.Error(fmt.Errorf("unsupported connection type %d", data[7])),
		)
	}

	switch data[1] {
	case 0:
		modelName = append(modelName, "DIN")
	case 1:
		modelName = append(modelName, "UNIVERSAL")
	case 2:
		modelName = append(modelName, "PLANE")
	case 3:
		modelName = append(modelName, "COM")
	default:
		m.log.Error(
			"build model name",
			zap.Error(fmt.Errorf("unsupported case type %d", data[7])),
		)
	}

	return strings.Join(modelName, " "), nil
}
