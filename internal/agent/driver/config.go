package driver

import "github.com/hashicorp/go-hclog"

type Config struct {
	logger hclog.Logger
}

func NewConfig(logger hclog.Logger) Config {
	return Config{
		logger: logger,
	}
}
