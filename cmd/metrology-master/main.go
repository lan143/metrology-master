package main

import (
	"flag"
	"github.com/lan143/metrology-master/internal/ha"
	"github.com/lan143/metrology-master/internal/scheduler"
	"github.com/lan143/metrology-master/pkg/cmd"
	"go.uber.org/zap"
)

type Command struct {
	config Config

	mqtt         MQTT
	serial       Serial
	meters       Meters
	scheduler    *scheduler.Scheduler
	discoveryMgr *ha.DiscoveryMgr

	log *zap.Logger
}

func (c *Command) Setup(flags *flag.FlagSet) {
	c.config.Export(flags)
}

func (c *Command) Init(log *zap.Logger) error {
	zap.ReplaceGlobals(log)
	c.log = log

	err := c.InitMQTT(*c.config.MQTT)
	if err != nil {
		return err
	}

	err = c.InitSerial(c.config.Serial)
	if err != nil {
		return err
	}

	c.scheduler = scheduler.NewScheduler(c.log)

	c.discoveryMgr = ha.NewDiscoveryMgr(c.mqtt.client, c.log)
	err = c.discoveryMgr.Init(*c.config.HA)
	if err != nil {
		return err
	}

	err = c.InitMeters(c.config.Meters)
	if err != nil {
		return err
	}

	return nil
}

func (c *Command) Run(ctx cmd.Context) error {
	err := c.DialMQTT()
	if err != nil {
		return err
	}
	defer c.CloseMQTT()

	err = c.discoveryMgr.Run()
	if err != nil {
		return err
	}

	ctx.Serve(c.scheduler)

	<-ctx.Shutdown()
	<-ctx.Done()

	return nil
}

func main() {
	cmd.Main(&Command{})
}
