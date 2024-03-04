package cmd

import (
	"flag"
	"github.com/lan143/metrology-master/pkg/flag/flagutil"
	"github.com/lan143/metrology-master/pkg/log"
	"time"
)

const (
	DefaultGracePeriod = 10
)

type config struct {
	Log   *log.Config
	Grace time.Duration
}

func (c *config) Export(flags *flag.FlagSet) {
	flagutil.Subset(flags, "log", func(set *flag.FlagSet) {
		c.Log = log.Export(set)
	})

	flags.DurationVar(
		&c.Grace,
		"grace.period",
		DefaultGracePeriod,
		"",
	)
}
