package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

func Read(path string) (cfg Config, err error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return
	}

	err = yaml.Unmarshal(data, &cfg)
	return
}
