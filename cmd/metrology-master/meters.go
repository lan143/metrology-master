package main

import (
	"context"
	"fmt"
	"github.com/lan143/metrology-master/internal/job"
	"github.com/lan143/metrology-master/internal/meter"
	pulsar_t1 "github.com/lan143/metrology-master/internal/meter/pulsar-electro"
	pulsar_m "github.com/lan143/metrology-master/internal/protocol/pulsar"
)

type Meters struct {
	electricMeters map[string]meter.ElectricMeter
}

func (c *Command) InitMeters(configs map[string]*meter.Config) error {
	c.meters.electricMeters = make(map[string]meter.ElectricMeter)

	for name, config := range configs {
		switch config.Type {
		case pulsar_t1.Type:
			port, ok := c.serial.ports[config.Port]
			if !ok {
				return fmt.Errorf("port \"%s\" not found in ports list", config.Port)
			}

			protocol := pulsar_m.NewPulsar(port, c.log)

			address, err := pulsar_m.ParseAddress(config.UID)
			if err != nil {
				return fmt.Errorf("parse address \"%s\": %s", config.UID, err.Error())
			}
			m := pulsar_t1.NewPulsarElectro(
				pulsar_t1.Config{Address: address},
				protocol,
				c.log,
			)
			err = m.Init(context.Background())
			if err != nil {
				return err
			}

			c.meters.electricMeters[name] = m
			c.scheduler.AddJob(
				job.NewUpdateMeterJob(
					m,
					c.mqtt.client,
					c.log,
				),
			)
			c.discoveryMgr.AddMeter(m)
		default:
			return fmt.Errorf("unsupported meter type \"%s\"", name)
		}
	}

	return nil
}
