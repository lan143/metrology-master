package main

import (
	"fmt"
	mqtt2 "github.com/eclipse/paho.mqtt.golang"
	"github.com/lan143/metrology-master/pkg/mqtt"
	"time"
)

type MQTT struct {
	client mqtt2.Client
}

func (c *Command) InitMQTT(config mqtt.Config) error {
	opts := mqtt2.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", config.Host, config.Port))
	opts.SetClientID(config.ClientID)

	if config.UserName != "" {
		opts.SetUsername(config.UserName)
	}

	if config.Password != "" {
		opts.SetPassword(config.Password)
	}

	c.mqtt.client = mqtt2.NewClient(opts)

	return nil
}

func (c *Command) DialMQTT() error {
	c.log.Debug("dial mqtt")

	if token := c.mqtt.client.Connect(); token.Wait() && token.Error() != nil {
		return token.Error()
	}

	c.log.Debug("dial mqtt - complete")

	return nil
}

func (c *Command) CloseMQTT() {
	c.log.Debug("close mqtt")

	c.mqtt.client.Disconnect(uint((10 * time.Second).Milliseconds()))
}
