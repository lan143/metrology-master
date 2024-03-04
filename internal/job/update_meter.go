package job

import (
	"context"
	"encoding/json"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/lan143/metrology-master/internal/meter"
	mqtt2 "github.com/lan143/metrology-master/internal/mqtt"
	"go.uber.org/zap"
)

type UpdateElectricMeterJob struct {
	meter      meter.ElectricMeter
	mqttClient mqtt.Client
	log        *zap.Logger
}

func NewUpdateMeterJob(
	meter meter.ElectricMeter,
	mqttClient mqtt.Client,
	log *zap.Logger,
) *UpdateElectricMeterJob {
	return &UpdateElectricMeterJob{
		meter:      meter,
		mqttClient: mqttClient,
		log:        log,
	}
}

func (j *UpdateElectricMeterJob) Execute(ctx context.Context) error {
	powerConsumption, err := j.meter.GetPowerConsumption(ctx)
	if err != nil {
		return err
	}

	state := mqtt2.State{PowerConsumption: powerConsumption}
	data, err := json.Marshal(state)
	if err != nil {
		return err
	}

	token := j.mqttClient.Publish(j.meter.GetParams().StateTopic, 1, false, data)
	if token.Error() != nil {
		return token.Error()
	}

	j.log.Debug(
		"publish state",
		zap.String("topic", j.meter.GetParams().StateTopic),
		zap.String("payload", string(data)),
	)

	return nil
}
