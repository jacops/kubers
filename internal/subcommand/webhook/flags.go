package webhook

import (
	"flag"
	"fmt"

	"github.com/hashicorp/consul/command/flags"
	"github.com/kelseyhightower/envconfig"
)

const (
	//DefaultLogLevel ...
	DefaultLogLevel = "info"

	//DefaultLogFormat ...
	DefaultLogFormat = "standard"
)

//KubersDSpecification ...
type KubersDSpecification struct {
	// Listen is the KUBERS_LISTEN environment variable.
	Listen string `split_words:"true" `

	// LogLevel is the KUBERS_LOG_LEVEL environment variable.
	LogLevel string `split_words:"true"`

	// LogFormat is the KUBERS_LOG_FORMAT environment variable
	LogFormat string `split_words:"true"`

	// TLSAuto is the KUBERS_TLS_AUTO environment variable.
	TLSAuto string `envconfig:"tls_auto"`

	// TLSAutoHosts is the KUBERS_TLS_AUTO_HOSTS environment variable.
	TLSAutoHosts string `envconfig:"tls_auto_hosts"`

	// TLSCertFile is the KUBERS_TLS_CERT_FILE environment variable.
	TLSCertFile string `envconfig:"tls_cert_file"`

	// TLSKeyFile is the KUBERS_TLS_KEY_FILE environment variable.
	TLSKeyFile string `envconfig:"tls_key_file"`

	// AgentImage is the AGENT_IMAGE environment variable.
	AgentImage string `envconfig:"agent_image"`

	// LogLevel is the KUBERS_AGENT_LOG_LEVEL environment variable.
	AgentLogLevel string `split_words:"true"`

	// LogFormat is the KUBERS_AGENT_LOG_FORMAT environment variable
	AgentLogFormat string `split_words:"true"`

	// AgentProvider is the KUBERS_AGENT_PROVIDER environment variable.
	AgentProvider string `envconfig:"agent_provider"`

	// AgentImage is the KUBERS_AGENT_PROVIDER_AWS_REGION environment variable.
	AWSRegion string `envconfig:"provider_aws_region"`
}

//KubersAgentSpecification ...
type KubersAgentSpecification struct {
	// LogLevel is the KUBERS_AGENT_LOG_LEVEL environment variable.
	LogLevel string `split_words:"true"`

	// LogFormat is the KUBERS_AGENT_LOG_FORMAT environment variable
	LogFormat string `split_words:"true"`

	// AgentProvider is the KUBERS_AGENT_PROVIDER environment variable.
	Provider string `envconfig:"provider"`

	// AgentImage is the KUBERS_AGENT_PROVIDER_AWS_REGION environment variable.
	AWSRegion string `envconfig:"provider_aws_region"`
}

func (c *Command) init() {

	defaultAgentImage := "jacops/kubers-agent:" + c.Version

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
	c.flagSet.StringVar(&c.flagAgentImage, "agent-image", defaultAgentImage,
		fmt.Sprintf("Docker image for Agent. Defaults to %q.", defaultAgentImage))
	c.flagSet.StringVar(&c.flagAgentLogLevel, "agent-log-level", DefaultLogLevel, "Agent log verbosity level. Supported values "+
		`(in order of detail) are "trace", "debug", "info", "warn", and "err".`)
	c.flagSet.StringVar(&c.flagAgentLogFormat, "agent-log-format", DefaultLogFormat, "Agent log output format. "+
		`Supported log formats: "standard", "json".`)
	c.flagSet.StringVar(&c.flagAgentProviderAWSRegion, "aws-region", "",
		fmt.Sprintf("AWS region where secret manager is deployed."))
	c.flagSet.StringVar(&c.flagAgentProvider, "provider-name", "",
		fmt.Sprintf("Provider name used to retrieve secrets."))

	c.help = flags.Usage(help, c.flagSet)
}

func (c *Command) parseKubersDEnvs() error {
	var envs KubersDSpecification

	err := envconfig.Process("kubersd", &envs)
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

func (c *Command) parseKubersAgentEnvs() error {
	var envs KubersAgentSpecification

	err := envconfig.Process("kubers_agent", &envs)
	if err != nil {
		return err
	}

	if envs.LogLevel != "" {
		c.flagAgentLogLevel = envs.LogLevel
	}

	if envs.LogFormat != "" {
		c.flagAgentLogFormat = envs.LogFormat
	}

	if envs.Provider != "" {
		c.flagAgentProvider = envs.Provider
	}

	if envs.AWSRegion != "" {
		c.flagAgentProviderAWSRegion = envs.AWSRegion
	}

	return nil
}
