package cmd

import (
	"context"
	"os"
	"syscall"
)

var (
	DefaultArgs    = os.Args
	DefaultLookup  = os.LookupEnv
	DefaultContext = context.Background()
	DefaultSignals = []os.Signal{
		syscall.SIGTERM,
		syscall.SIGINT,
	}
)

type (
	Option func(*command)
)
