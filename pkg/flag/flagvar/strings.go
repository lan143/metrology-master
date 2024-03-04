package flagvar

import "flag"

func Strings(flag *flag.FlagSet, name, usage string) *[]string {
	p := new([]string)
	StringsVar(flag, p, name, usage)

	return p
}

func StringsVar(flag *flag.FlagSet, p *[]string, name, usage string) {
	flag.Func(name, usage, func(s string) error {
		*p = append(*p, s)
		return nil
	})
}
