package parse

import (
	"gopkg.in/yaml.v3"
	"os"
)

func (fs *FlagSet) ParseFile(name string) error {
	data, err := os.ReadFile(name)
	if err != nil {
		return fs.fail(err)
	}

	err = fs.parseYAML(data)
	if err != nil {
		return fs.failf("parse %s: %w", name, err)
	}

	fs.update()

	return nil
}

func (fs *FlagSet) parseYAML(data []byte) error {
	var m map[string]any
	err := yaml.Unmarshal(data, &m)
	if err != nil {
		return err
	}

	return Setup(m, fs.SetUnspecified)
}
