package ha

import (
	"encoding/json"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/lan143/metrology-master/internal/ha/entity"
	"github.com/lan143/metrology-master/internal/ha/enum"
	"github.com/lan143/metrology-master/internal/meter"
	"go.uber.org/zap"
	"strings"
)

type DiscoveryMgr struct {
	config Config
	meters []meter.Meter

	mqttClient mqtt.Client
	log        *zap.Logger
}

func NewDiscoveryMgr(mqttClient mqtt.Client, log *zap.Logger) *DiscoveryMgr {
	return &DiscoveryMgr{
		mqttClient: mqttClient,
		log:        log,
	}
}

func (m *DiscoveryMgr) Init(config Config) error {
	m.config = config

	return nil
}

func (m *DiscoveryMgr) AddMeter(mtr meter.Meter) {
	m.log.Debug(
		"add meter to discovery manager",
		zap.Any("meter", mtr),
	)
	m.meters = append(m.meters, mtr)
}

func (m *DiscoveryMgr) Run() error {
	m.log.Debug("discovery manager run")

	if !m.config.AutoDiscovery {
		return nil
	}

	for i := range m.meters {
		err := m.sendDiscovery(m.meters[i])
		if err != nil {
			return err
		}

	}

	return nil
}

func (m *DiscoveryMgr) sendDiscovery(mtr meter.Meter) error {
	params := mtr.GetParams()
	if params.Flags.HasPowerConsumption() {
		data, err := m.buildDiscoveryPowerConsumption(mtr)
		if err != nil {
			return err
		}

		topic := m.buildDiscoveryTopic(
			"sensor",
			"power-consumption",
			params.UID,
		)
		token := m.mqttClient.Publish(
			topic,
			1,
			false,
			data,
		)
		if token.Error() != nil {
			return token.Error()
		}

		m.log.Debug(
			"publish discovery",
			zap.String("type", "PowerConsumption"),
			zap.String("topic", topic),
			zap.String("payload", string(data)),
		)
	}

	if params.Flags.HasFrequency() {
		data, err := m.buildDiscoveryFrequency(mtr)
		if err != nil {
			return err
		}

		topic := m.buildDiscoveryTopic(
			"sensor",
			"frequency",
			params.UID,
		)
		token := m.mqttClient.Publish(
			topic,
			1,
			false,
			data,
		)
		if token.Error() != nil {
			return token.Error()
		}

		m.log.Debug(
			"publish discovery",
			zap.String("type", "Frequency"),
			zap.String("topic", topic),
			zap.String("payload", string(data)),
		)
	}

	if params.Flags.HasVoltage() {
		data, err := m.buildDiscoveryVoltage(mtr)
		if err != nil {
			return err
		}

		topic := m.buildDiscoveryTopic(
			"sensor",
			"voltage",
			params.UID,
		)
		token := m.mqttClient.Publish(
			topic,
			1,
			false,
			data,
		)
		if token.Error() != nil {
			return token.Error()
		}

		m.log.Debug(
			"publish discovery",
			zap.String("type", "Voltage"),
			zap.String("topic", topic),
			zap.String("payload", string(data)),
		)
	}

	if params.Flags.HasCurrent() {
		data, err := m.buildDiscoveryCurrent(mtr)
		if err != nil {
			return err
		}

		topic := m.buildDiscoveryTopic(
			"sensor",
			"current",
			params.UID,
		)
		token := m.mqttClient.Publish(
			topic,
			1,
			false,
			data,
		)
		if token.Error() != nil {
			return token.Error()
		}

		m.log.Debug(
			"publish discovery",
			zap.String("type", "Current"),
			zap.String("topic", topic),
			zap.String("payload", string(data)),
		)
	}

	if params.Flags.HasActivePower() {
		data, err := m.buildDiscoveryActivePower(mtr)
		if err != nil {
			return err
		}

		topic := m.buildDiscoveryTopic(
			"sensor",
			"active_power",
			params.UID,
		)
		token := m.mqttClient.Publish(
			topic,
			1,
			false,
			data,
		)
		if token.Error() != nil {
			return token.Error()
		}

		m.log.Debug(
			"publish discovery",
			zap.String("type", "Active power"),
			zap.String("topic", topic),
			zap.String("payload", string(data)),
		)
	}

	if params.Flags.HasReactivePower() {
		data, err := m.buildDiscoveryReactivePower(mtr)
		if err != nil {
			return err
		}

		topic := m.buildDiscoveryTopic(
			"sensor",
			"reactive_power",
			params.UID,
		)
		token := m.mqttClient.Publish(
			topic,
			1,
			false,
			data,
		)
		if token.Error() != nil {
			return token.Error()
		}

		m.log.Debug(
			"publish discovery",
			zap.String("type", "Reactive power"),
			zap.String("topic", topic),
			zap.String("payload", string(data)),
		)
	}

	if params.Flags.HasFullPower() {
		data, err := m.buildDiscoveryFullPower(mtr)
		if err != nil {
			return err
		}

		topic := m.buildDiscoveryTopic(
			"sensor",
			"full_power",
			params.UID,
		)
		token := m.mqttClient.Publish(
			topic,
			1,
			false,
			data,
		)
		if token.Error() != nil {
			return token.Error()
		}

		m.log.Debug(
			"publish discovery",
			zap.String("type", "Full power"),
			zap.String("topic", topic),
			zap.String("payload", string(data)),
		)
	}

	return nil
}

// example: homeassistant/sensor/0x08833976/power-consumption/config
func (m *DiscoveryMgr) buildDiscoveryTopic(oType string, name string, uid string) string {
	return fmt.Sprintf("%s/%s/%s/%s/config", m.config.Prefix, oType, name, uid)
}

func (m *DiscoveryMgr) buildDiscoveryPowerConsumption(mtr meter.Meter) ([]byte, error) {
	objectID := strings.ToLower(
		strings.ReplaceAll(mtr.GetParams().Name, " ", "_"),
	)
	uniqueID := mtr.GetParams().UID + "_" + objectID

	obj := entity.Sensor{
		StateTopic:        mtr.GetParams().StateTopic,
		ValueTemplate:     "{{ value_json.powerConsumption }}",
		UnitOfMeasurement: "kWh",
		DeviceClass:       enum.DeviceClassEnergy,
		StateClass:        enum.StateClassTotal,
		Base: entity.Base{
			Device: entity.Device{
				Identifiers:  []string{mtr.GetParams().UID},
				HWVersion:    mtr.GetParams().HWVersion,
				Manufacturer: mtr.GetParams().Manufacturer,
				Model:        mtr.GetParams().Model,
				Name:         mtr.GetParams().Name,
				SWVersion:    mtr.GetParams().SWVersion,
			},
			ObjectID:    objectID,
			UniqueID:    uniqueID,
			ForceUpdate: true,
		},
	}
	data, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (m *DiscoveryMgr) buildDiscoveryFrequency(mtr meter.Meter) ([]byte, error) {
	objectID := strings.ToLower(
		strings.ReplaceAll(mtr.GetParams().Name, " ", "_"),
	)
	uniqueID := mtr.GetParams().UID + "_" + objectID + "_frequency"

	obj := entity.Sensor{
		StateTopic:        mtr.GetParams().StateTopic,
		ValueTemplate:     "{{ value_json.frequency }}",
		UnitOfMeasurement: "Hz",
		DeviceClass:       enum.DeviceClassFrequency,
		StateClass:        enum.StateClassMeasurement,
		Base: entity.Base{
			Device: entity.Device{
				Identifiers:  []string{mtr.GetParams().UID},
				HWVersion:    mtr.GetParams().HWVersion,
				Manufacturer: mtr.GetParams().Manufacturer,
				Model:        mtr.GetParams().Model,
				Name:         mtr.GetParams().Name,
				SWVersion:    mtr.GetParams().SWVersion,
			},
			ObjectID:    objectID,
			UniqueID:    uniqueID,
			ForceUpdate: true,
		},
	}
	data, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (m *DiscoveryMgr) buildDiscoveryVoltage(mtr meter.Meter) ([]byte, error) {
	objectID := strings.ToLower(
		strings.ReplaceAll(mtr.GetParams().Name, " ", "_"),
	)
	uniqueID := mtr.GetParams().UID + "_" + objectID + "_voltage"

	obj := entity.Sensor{
		StateTopic:        mtr.GetParams().StateTopic,
		ValueTemplate:     "{{ value_json.voltage }}",
		UnitOfMeasurement: "V",
		DeviceClass:       enum.DeviceClassVoltage,
		StateClass:        enum.StateClassMeasurement,
		Base: entity.Base{
			Device: entity.Device{
				Identifiers:  []string{mtr.GetParams().UID},
				HWVersion:    mtr.GetParams().HWVersion,
				Manufacturer: mtr.GetParams().Manufacturer,
				Model:        mtr.GetParams().Model,
				Name:         mtr.GetParams().Name,
				SWVersion:    mtr.GetParams().SWVersion,
			},
			ObjectID:    objectID,
			UniqueID:    uniqueID,
			ForceUpdate: true,
		},
	}
	data, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (m *DiscoveryMgr) buildDiscoveryCurrent(mtr meter.Meter) ([]byte, error) {
	objectID := strings.ToLower(
		strings.ReplaceAll(mtr.GetParams().Name, " ", "_"),
	)
	uniqueID := mtr.GetParams().UID + "_" + objectID + "_current"

	obj := entity.Sensor{
		StateTopic:        mtr.GetParams().StateTopic,
		ValueTemplate:     "{{ value_json.current }}",
		UnitOfMeasurement: "A",
		DeviceClass:       enum.DeviceClassCurrent,
		StateClass:        enum.StateClassMeasurement,
		Base: entity.Base{
			Device: entity.Device{
				Identifiers:  []string{mtr.GetParams().UID},
				HWVersion:    mtr.GetParams().HWVersion,
				Manufacturer: mtr.GetParams().Manufacturer,
				Model:        mtr.GetParams().Model,
				Name:         mtr.GetParams().Name,
				SWVersion:    mtr.GetParams().SWVersion,
			},
			ObjectID:    objectID,
			UniqueID:    uniqueID,
			ForceUpdate: true,
		},
	}
	data, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (m *DiscoveryMgr) buildDiscoveryActivePower(mtr meter.Meter) ([]byte, error) {
	objectID := strings.ToLower(
		strings.ReplaceAll(mtr.GetParams().Name, " ", "_"),
	)
	uniqueID := mtr.GetParams().UID + "_" + objectID + "_active_power"

	obj := entity.Sensor{
		StateTopic:        mtr.GetParams().StateTopic,
		ValueTemplate:     "{{ value_json.activePower }}",
		UnitOfMeasurement: "W",
		DeviceClass:       enum.DeviceClassPower,
		StateClass:        enum.StateClassMeasurement,
		Base: entity.Base{
			Name: "Active power",
			Device: entity.Device{
				Identifiers:  []string{mtr.GetParams().UID},
				HWVersion:    mtr.GetParams().HWVersion,
				Manufacturer: mtr.GetParams().Manufacturer,
				Model:        mtr.GetParams().Model,
				Name:         mtr.GetParams().Name,
				SWVersion:    mtr.GetParams().SWVersion,
			},
			ObjectID:    objectID,
			UniqueID:    uniqueID,
			ForceUpdate: true,
		},
	}
	data, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (m *DiscoveryMgr) buildDiscoveryReactivePower(mtr meter.Meter) ([]byte, error) {
	objectID := strings.ToLower(
		strings.ReplaceAll(mtr.GetParams().Name, " ", "_"),
	)
	uniqueID := mtr.GetParams().UID + "_" + objectID + "_reactive_power"

	obj := entity.Sensor{
		StateTopic:        mtr.GetParams().StateTopic,
		ValueTemplate:     "{{ value_json.reactivePower }}",
		UnitOfMeasurement: "VAr",
		DeviceClass:       enum.DeviceClassPower,
		StateClass:        enum.StateClassMeasurement,
		Base: entity.Base{
			Name: "Reactive power",
			Device: entity.Device{
				Identifiers:  []string{mtr.GetParams().UID},
				HWVersion:    mtr.GetParams().HWVersion,
				Manufacturer: mtr.GetParams().Manufacturer,
				Model:        mtr.GetParams().Model,
				Name:         mtr.GetParams().Name,
				SWVersion:    mtr.GetParams().SWVersion,
			},
			ObjectID:    objectID,
			UniqueID:    uniqueID,
			ForceUpdate: true,
		},
	}
	data, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (m *DiscoveryMgr) buildDiscoveryFullPower(mtr meter.Meter) ([]byte, error) {
	objectID := strings.ToLower(
		strings.ReplaceAll(mtr.GetParams().Name, " ", "_"),
	)
	uniqueID := mtr.GetParams().UID + "_" + objectID + "_full_power"

	obj := entity.Sensor{
		StateTopic:        mtr.GetParams().StateTopic,
		ValueTemplate:     "{{ value_json.fullPower }}",
		UnitOfMeasurement: "VA",
		DeviceClass:       enum.DeviceClassPower,
		StateClass:        enum.StateClassMeasurement,
		Base: entity.Base{
			Name: "Full power",
			Device: entity.Device{
				Identifiers:  []string{mtr.GetParams().UID},
				HWVersion:    mtr.GetParams().HWVersion,
				Manufacturer: mtr.GetParams().Manufacturer,
				Model:        mtr.GetParams().Model,
				Name:         mtr.GetParams().Name,
				SWVersion:    mtr.GetParams().SWVersion,
			},
			ObjectID:    objectID,
			UniqueID:    uniqueID,
			ForceUpdate: true,
		},
	}
	data, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}

	return data, nil
}
