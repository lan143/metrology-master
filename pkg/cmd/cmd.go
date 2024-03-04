package cmd

import (
	"flag"
	"go.uber.org/zap"
)

type (
	Command interface {
		Run(Context) error
	}

	Setuper interface {
		Setup(set *flag.FlagSet)
	}

	Initer interface {
		Init(logger *zap.Logger) error
	}
)

func Main(cmd Command, opts ...Option) {
	c := &command{
		Args:    DefaultArgs,
		Lookup:  DefaultLookup,
		Context: DefaultContext,
		Signals: DefaultSignals,
		Command: cmd,
	}

	for i := range opts {
		opts[i](c)
	}

	c.Main()
}

func Setup(cmd Command, flags *flag.FlagSet) {
	if x, ok := cmd.(Setuper); ok {
		x.Setup(flags)
	}
}

func Init(cmd Command, log *zap.Logger) error {
	if x, ok := cmd.(Initer); ok {
		return x.Init(log)
	}

	return nil
}
