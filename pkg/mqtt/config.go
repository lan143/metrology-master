package mqtt

import (
	"flag"
)

type Config struct {
	Host     string
	Port     int
	ClientID string
	UserName string
	Password string
}

func Export(sub *flag.FlagSet) *Config {
	c := &Config{}

	sub.StringVar(
		&c.Host,
		"host",
		"",
		"",
	)
	sub.IntVar(
		&c.Port,
		"port",
		1883,
		"",
	)
	sub.StringVar(
		&c.ClientID,
		"client-id",
		"metrology-master",
		"",
	)
	sub.StringVar(
		&c.UserName,
		"username",
		"",
		"",
	)
	sub.StringVar(
		&c.Password,
		"password",
		"",
		"",
	)

	return c
}
