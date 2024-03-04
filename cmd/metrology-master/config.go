package main

import (
	"flag"
	"github.com/lan143/metrology-master/internal/ha"
	"github.com/lan143/metrology-master/internal/meter"
	"github.com/lan143/metrology-master/pkg/flag/flagutil"
	"github.com/lan143/metrology-master/pkg/mqtt"
	"github.com/lan143/metrology-master/pkg/serial"
)

type Config struct {
	MQTT   *mqtt.Config
	HA     *ha.Config
	Serial map[string]*serial.Config
	Meters map[string]*meter.Config
}

func (c *Config) Export(flags *flag.FlagSet) {
	c.Meters = make(map[string]*meter.Config)

	flagutil.Subset(flags, "mqtt", func(sub *flag.FlagSet) {
		c.MQTT = mqtt.Export(sub)
	})
	flagutil.Subset(flags, "home-assistant", func(set *flag.FlagSet) {
		c.HA = ha.Export(set)
	})
	flagutil.Subset(flags, "serial", func(sub *flag.FlagSet) {
		c.Serial = make(map[string]*serial.Config)

		flagutil.Func(sub, "include", "", func(name string) error {
			flagutil.Subset(sub, name, func(sub *flag.FlagSet) {
				c.Serial[name] = serial.Export(sub)
			})

			return nil
		})
	})
	flagutil.Subset(flags, "meters", func(sub *flag.FlagSet) {
		c.Meters = make(map[string]*meter.Config)

		flagutil.Func(sub, "include", "", func(name string) error {
			flagutil.Subset(sub, name, func(sub *flag.FlagSet) {
				c.Meters[name] = meter.Export(sub)
			})

			return nil
		})
	})
}
