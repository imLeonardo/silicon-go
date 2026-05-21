package flag

import (
	"flag"
)

type Flag struct {
	registry.NopHook

	ConfigFilePath string
}

func NewFlag() (*Flag, error) {
	configFilePath := flag.String("f", "conf.d/default.yaml", "config file")
	flag.Parse()
	return &Flag{
		ConfigFilePath: *configFilePath,
	}, nil
}
