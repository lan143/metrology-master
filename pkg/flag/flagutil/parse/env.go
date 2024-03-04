package parse

import (
	"flag"
	"strings"
)

type LookupFunc func(string) (string, bool)

func (fs *FlagSet) ParseEnv(lookup LookupFunc) (err error) {
	replacer := strings.NewReplacer(".", "_", "-", "_")

	fs.VisitUnspecified(func(f *flag.Flag) {
		if err != nil {
			return
		}

		key := strings.ToUpper(f.Name)
		key = replacer.Replace(key)

		v, found := lookup(key)
		if found {
			err = fs.Set(f.Name, v)
		}
	})

	if err != nil {
		return fs.fail(err)
	}

	fs.update()

	return nil
}
