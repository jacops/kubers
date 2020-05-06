package fetchsecrets

import (
	"flag"

	"github.com/hashicorp/consul/command/flags"
	"github.com/kelseyhightower/envconfig"
)

const (
	DefaultLogLevel  = "info"
	DefaultLogFormat = "standard"
)

type Specification struct {
	// LogLevel is the AGENT_INJECT_LOG_LEVEL environment variable.
	LogLevel string `split_words:"true"`

	// LogFormat is the AGENT_INJECT_LOG_FORMAT environment variable
	LogFormat string `split_words:"true"`

	Config string `split_words:"true"`
}

func (c *Command) init() {
	c.flagSet = flag.NewFlagSet("", flag.ContinueOnError)
	c.flagSet.StringVar(&c.flagLogLevel, "log-level", DefaultLogLevel, "Log verbosity level. Supported values "+
		`(in order of detail) are "trace", "debug", "info", "warn", and "err".`)
	c.flagSet.StringVar(&c.flagLogFormat, "log-format", DefaultLogFormat, "Log output format. "+
		`Supported log formats: "standard", "json".`)
	c.flagSet.StringVar(&c.flagAgentConfig, "agent-config", "", "base64 Agent config")

	c.help = flags.Usage(help, c.flagSet)
}

func (c *Command) parseEnvs() error {
	var envs Specification

	err := envconfig.Process("agent", &envs)
	if err != nil {
		return err
	}

	if envs.LogLevel != "" {
		c.flagLogLevel = envs.LogLevel
	}

	if envs.LogFormat != "" {
		c.flagLogFormat = envs.LogFormat
	}

	if envs.Config != "" {
		c.flagAgentConfig = envs.Config
	}

	return nil
}
