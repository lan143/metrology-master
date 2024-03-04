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
	if mtr.GetParams().Flags.HasPowerConsumption() {
		data, err := m.buildDiscoveryPowerConsumption(mtr)
		if err != nil {
			return err
		}

		topic := m.buildDiscoveryTopic(
			"sensor",
			"power-consumption",
			mtr.GetParams().UID,
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
