package config

import "time"

type Config struct {
	Search struct {
		Timeout time.Duration `yaml:"timeout"`
		Actions []action      `yaml:"actions"`
	} `yaml:"search"`
}

// Action for research.
type action struct {
	Type   string   `yaml:"type"`
	Params []string `yaml:"params"`
}
