package log

import (
	"flag"
	"github.com/lan143/metrology-master/pkg/log/json"
	"github.com/lan143/metrology-master/pkg/log/zaplog"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	DefaultOutput string = "stdout"
	DefaultLevel  string = "info"
)

type Config struct {
	Output string
	Level  string
}

func Export(flags *flag.FlagSet) *Config {
	c := &Config{}

	flags.StringVar(
		&c.Output,
		"output",
		DefaultOutput,
		"",
	)
	flags.StringVar(
		&c.Level,
		"level",
		DefaultLevel,
		"",
	)

	return c
}

func New(config Config) (*zap.Logger, error) {
	w, err := zaplog.ParseOutput(config.Output)
	if err != nil {
		return nil, err
	}

	enabler, err := zaplog.ParseLevel(config.Level)
	if err != nil {
		return nil, err
	}

	core := zapcore.NewCore(
		json.NewEncoder(),
		w,
		enabler,
	)
	log := zap.New(
		core,
		zap.AddCaller(),
	)

	return log, nil
}
