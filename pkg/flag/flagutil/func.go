package flagutil

import (
	"flag"
)

type exporter struct {
	flag *flag.FlagSet
	set  func(string) error

	origin *flag.FlagSet
	prefix string
}

func (e *exporter) String() string {
	return ""
}

func (e *exporter) Set(value string) error {
	if e.origin == nil {
		return e.set(value)
	}

	defined := make(map[string]bool)
	e.flag.VisitAll(func(f *flag.Flag) {
		defined[f.Name] = true
	})

	err := e.set(value)
	if err != nil {
		return err
	}

	e.flag.VisitAll(func(f *flag.Flag) {
		if defined[f.Name] {
			delete(defined, f.Name)
			return
		}

		defined[f.Name] = true
	})

	for name := range defined {
		f := e.flag.Lookup(name)

		e.origin.Var(
			f.Value,
			e.prefix+name,
			f.Usage,
		)
	}

	return nil
}

func (e *exporter) Trace(super *flag.FlagSet, prefix string) {
	if n := len(prefix); n > 0 {
		if prefix[n-1] != '.' {
			prefix += "."
		}
	}

	e.prefix = prefix + e.prefix
	e.origin = super
}

func Func(flag *flag.FlagSet, name, usage string, fn func(string) error) {
	value := exporter{flag: flag, set: fn}

	flag.Var(&value, name, usage)
}
