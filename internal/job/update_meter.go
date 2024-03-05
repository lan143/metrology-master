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
	var (
		err    error
		state  = mqtt2.State{}
		params = j.meter.GetParams()
	)

	if params.Flags.HasPowerConsumption() {
		state.PowerConsumption, err = j.meter.GetPowerConsumption(ctx)
		if err != nil {
			return err
		}
	}

	if params.Flags.HasFrequency() {
		state.Frequency, err = j.meter.GetFrequency(ctx)
		if err != nil {
			return err
		}
	}

	if params.Flags.HasVoltage() {
		state.Voltage, err = j.meter.GetVoltage(ctx)
		if err != nil {
			return err
		}
	}

	if params.Flags.HasCurrent() {
		state.Current, err = j.meter.GetCurrent(ctx)
		if err != nil {
			return err
		}
	}

	if params.Flags.HasActivePower() {
		state.ActivePower, err = j.meter.GetActivePower(ctx)
		if err != nil {
			return err
		}
	}

	if params.Flags.HasReactivePower() {
		state.ReactivePower, err = j.meter.GetReactivePower(ctx)
		if err != nil {
			return err
		}
	}

	if params.Flags.HasFullPower() {
		state.FullPower, err = j.meter.GetFullPower(ctx)
		if err != nil {
			return err
		}
	}

	data, err := json.Marshal(state)
	if err != nil {
		return err
	}

	token := j.mqttClient.Publish(params.StateTopic, 1, false, data)
	if token.Error() != nil {
		return token.Error()
	}

	j.log.Debug(
		"publish state",
		zap.String("topic", params.StateTopic),
		zap.String("payload", string(data)),
	)

	return nil
}
