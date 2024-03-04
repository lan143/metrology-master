package meter

import (
	"flag"
	"github.com/lan143/metrology-master/pkg/flag/flagutil"
)

type Config struct {
	Type   string
	UID    string
	Port   string
	Export []string
}

func Export(flags *flag.FlagSet) *Config {
	c := &Config{}

	flags.StringVar(
		&c.Type,
		"type",
		"",
		"",
	)
	flags.StringVar(
		&c.UID,
		"uid",
		"",
		"",
	)
	flags.StringVar(
		&c.Port,
		"port",
		"",
		"",
	)
	flagutil.Func(flags, "export", "", func(name string) error {
		var val string
		flags.StringVar(
			&val,
			name,
			"",
			"",
		)
		c.Export = append(c.Export, val)

		return nil
	})

	return c
}
