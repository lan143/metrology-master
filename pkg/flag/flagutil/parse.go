package flagutil

import (
	"flag"
	"github.com/lan143/metrology-master/pkg/flag/flagutil/parse"
)

type LookupFunc = parse.LookupFunc

func ParseEnv(set *flag.FlagSet, lookup LookupFunc) error {
	fs := parse.NewFlagSet(set)

	return fs.ParseEnv(lookup)
}

func ParseFile(set *flag.FlagSet, name string) error {
	fs := parse.NewFlagSet(set)

	return fs.ParseFile(name)
}
