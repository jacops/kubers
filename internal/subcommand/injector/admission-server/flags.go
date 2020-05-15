package admissionserver

import (
	"flag"
	"fmt"

	"github.com/hashicorp/consul/command/flags"
	"github.com/jacops/kubers/internal/injector"
	"github.com/kelseyhightower/envconfig"
)

const (
	DefaultLogLevel  = "info"
	DefaultLogFormat = "standard"
)

type Specification struct {
	// Listen is the AGENT_INJECT_LISTEN environment variable.
	Listen string `split_words:"true" `

	// LogLevel is the AGENT_INJECT_LOG_LEVEL environment variable.
	LogLevel string `split_words:"true"`

	// LogFormat is the AGENT_INJECT_LOG_FORMAT environment variable
	LogFormat string `split_words:"true"`

	// TLSAuto is the AGENT_INJECT_TLS_AUTO environment variable.
	TLSAuto string `envconfig:"tls_auto"`

	// TLSAutoHosts is the AGENT_INJECT_TLS_AUTO_HOSTS environment variable.
	TLSAutoHosts string `envconfig:"tls_auto_hosts"`

	// TLSCertFile is the AGENT_INJECT_TLS_CERT_FILE environment variable.
	TLSCertFile string `envconfig:"tls_cert_file"`

	// TLSKeyFile is the AGENT_INJECT_TLS_KEY_FILE environment variable.
	TLSKeyFile string `envconfig:"tls_key_file"`

	// AgentImage is the AGENT_IMAGE environment variable.
	AgentImage string `envconfig:"agent_image"`
}

func (c *Command) init() {
	c.flagSet = flag.NewFlagSet("", flag.ContinueOnError)
	c.flagSet.StringVar(&c.flagListen, "listen", ":8080", "Address to bind listener to.")
	c.flagSet.StringVar(&c.flagLogLevel, "log-level", DefaultLogLevel, "Log verbosity level. Supported values "+
		`(in order of detail) are "trace", "debug", "info", "warn", and "err".`)
	c.flagSet.StringVar(&c.flagLogFormat, "log-format", DefaultLogFormat, "Log output format. "+
		`Supported log formats: "standard", "json".`)
	c.flagSet.StringVar(&c.flagAutoName, "tls-auto", "",
		"MutatingWebhookConfiguration name. If specified, will auto generate cert bundle.")
	c.flagSet.StringVar(&c.flagAutoHosts, "tls-auto-hosts", "",
		"Comma-separated hosts for auto-generated TLS cert. If specified, will auto generate cert bundle.")
	c.flagSet.StringVar(&c.flagCertFile, "tls-cert-file", "",
		"PEM-encoded TLS certificate to serve. If blank, will generate random cert.")
	c.flagSet.StringVar(&c.flagKeyFile, "tls-key-file", "",
		"PEM-encoded TLS private key to serve. If blank, will generate random cert.")
	c.flagSet.StringVar(&c.flagAgentImage, "agent-image", injector.DefaultVaultImage,
		fmt.Sprintf("Docker image for Agent. Defaults to %q.", injector.DefaultVaultImage))

	c.help = flags.Usage(help, c.flagSet)
}

func (c *Command) parseEnvs() error {
	var envs Specification

	err := envconfig.Process("agent_inject", &envs)
	if err != nil {
		return err
	}

	if envs.Listen != "" {
		c.flagListen = envs.Listen
	}

	if envs.LogLevel != "" {
		c.flagLogLevel = envs.LogLevel
	}

	if envs.LogFormat != "" {
		c.flagLogFormat = envs.LogFormat
	}

	if envs.TLSAuto != "" {
		c.flagAutoName = envs.TLSAuto
	}

	if envs.TLSAutoHosts != "" {
		c.flagAutoHosts = envs.TLSAutoHosts
	}

	if envs.TLSCertFile != "" {
		c.flagCertFile = envs.TLSCertFile
	}

	if envs.TLSKeyFile != "" {
		c.flagKeyFile = envs.TLSKeyFile
	}

	if envs.AgentImage != "" {
		c.flagAgentImage = envs.AgentImage
	}

	return nil
}
