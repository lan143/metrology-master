package serial

import (
	"flag"
)

type Config struct {
	Port     string
	BaudRate int
	DataBits int
	StopBits int
	Parity   int
}

func Export(flag *flag.FlagSet) *Config {
	c := &Config{}

	flag.StringVar(
		&c.Port,
		"port",
		"",
		"",
	)
	flag.IntVar(
		&c.BaudRate,
		"baud-rate",
		9600,
		"",
	)
	flag.IntVar(
		&c.DataBits,
		"data-bits",
		8,
		"",
	)
	flag.IntVar(
		&c.StopBits,
		"stop-bits",
		1,
		"",
	)
	flag.IntVar(
		&c.Parity,
		"parity",
		0,
		"",
	)

	return c
}
