package ha

import "flag"

type Config struct {
	AutoDiscovery bool
	Prefix        string
}

func Export(flags *flag.FlagSet) *Config {
	c := &Config{}

	flags.BoolVar(
		&c.AutoDiscovery,
		"auto-discovery",
		false,
		"",
	)
	flags.StringVar(
		&c.Prefix,
		"prefix",
		"",
		"",
	)

	return c
}
