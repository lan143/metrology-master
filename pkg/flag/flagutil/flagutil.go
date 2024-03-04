package flagutil

import (
	"flag"
)

type Tracer interface {
	Trace(super *flag.FlagSet, prefix string)
}

func Subset(super *flag.FlagSet, prefix string, fn func(set *flag.FlagSet)) {
	if n := len(prefix); n > 0 {
		if prefix[n-1] != '.' {
			prefix += "."
		}
	}

	sub := flag.NewFlagSet(super.Name(), 0)
	fn(sub)

	sub.VisitAll(func(f *flag.Flag) {
		if t, ok := f.Value.(Tracer); ok {
			t.Trace(super, prefix)
		}

		super.Var(f.Value, prefix+f.Name, f.Usage)
	})
}
