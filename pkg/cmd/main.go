package cmd

import (
	"context"
	"flag"
	"fmt"
	"github.com/lan143/metrology-master/pkg/flag/flagutil"
	log2 "github.com/lan143/metrology-master/pkg/log"
	"go.uber.org/zap"
	"os"
	"path/filepath"
)

type command struct {
	Args   []string
	Lookup flagutil.LookupFunc

	Context context.Context
	Signals []os.Signal

	Command Command
	config  config
}

func (c *command) Main() {
	l, err := c.init()
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer l.Sync()

	l.Info("starting app")
	defer l.Info("app stopped")

	err = c.run(l)
	if err != nil {
		l.Fatal("app aborted", zap.Error(err))
	}
}

func (c *command) init() (*zap.Logger, error) {
	err := c.parse()
	if err != nil {
		return nil, err
	}

	return log2.New(*c.config.Log)
}

func (c *command) parse() error {
	name := filepath.Base(c.Args[0])
	f := flag.NewFlagSet(name, flag.ExitOnError)

	var cnf string
	f.StringVar(&cnf, "config", "", "")

	c.config.Export(f)
	Setup(c.Command, f)

	err := f.Parse(c.Args[1:])
	if err != nil {
		return err
	}

	if cnf != "" {
		err = flagutil.ParseFile(f, cnf)
	}

	if err == nil {
		err = flagutil.ParseEnv(f, c.Lookup)
	}

	return err
}

func (c *command) run(log *zap.Logger) error {
	err := Init(c.Command, log)
	if err != nil {
		return err
	}

	ctx, cancel := NewContext(c.Context, log)
	defer cancel()

	ctx.Listen(c.Signals, c.config.Grace)
	ctx.Start(func() error {
		return c.Command.Run(ctx)
	})

	return ctx.Wait()
}
