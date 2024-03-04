package parse

import (
	"flag"
	"fmt"
	"os"
)

type FlagSet struct {
	origin    *flag.FlagSet
	specified map[string]bool
}

func NewFlagSet(origin *flag.FlagSet) *FlagSet {
	fs := &FlagSet{
		origin:    origin,
		specified: make(map[string]bool),
	}
	fs.update()

	return fs
}

func (fs *FlagSet) update() {
	fs.origin.Visit(func(f *flag.Flag) {
		fs.specified[f.Name] = true
	})
}

func (fs *FlagSet) VisitUnspecified(fn func(*flag.Flag)) {
	fs.origin.VisitAll(func(f *flag.Flag) {
		if fs.specified[f.Name] {
			return
		}

		fn(f)
	})
}

func (fs *FlagSet) SetUnspecified(name, value string) error {
	if fs.specified[name] {
		return nil
	}

	return fs.Set(name, value)
}

func (fs *FlagSet) Set(name, value string) error {
	err := fs.origin.Set(name, value)
	if err != nil {
		return fmt.Errorf("invalid value %q for %s: %w", value, name, err)
	}

	return nil
}

func (fs *FlagSet) failf(format string, a ...any) error {
	return fs.fail(fmt.Errorf(format, a...))
}

func (fs *FlagSet) fail(err error) error {
	_, _ = fmt.Fprintln(fs.origin.Output(), err)

	switch fs.origin.ErrorHandling() {
	case flag.ExitOnError:
		os.Exit(2)
	case flag.PanicOnError:
		panic(err)
	case flag.ContinueOnError:
	}

	return err
}
