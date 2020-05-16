package fetchsecrets

import (
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"sync"

	"github.com/hashicorp/go-hclog"
	"github.com/jacops/kubers/internal/agent"
	"github.com/jacops/kubers/internal/common"
	"github.com/mitchellh/cli"
)

type Command struct {
	UI cli.Ui

	flagLogLevel    string // Log verbosity
	flagLogFormat   string // Log format
	flagAgentConfig string // Agent config

	flagSet *flag.FlagSet

	once sync.Once
	help string
}

func (c *Command) Run(args []string) int {

	c.once.Do(c.init)
	if err := c.flagSet.Parse(args); err != nil {
		return 1
	}

	if err := c.parseEnvs(); err != nil {
		c.UI.Error(fmt.Sprintf("Error parsing environment variables: %s", err))
		return 1
	}

	var agentConfig *agent.Config
	agentConfigJSON, err := base64.StdEncoding.DecodeString(c.flagAgentConfig)
	if err != nil {
		c.UI.Error(fmt.Sprintf("Could not decode agent config: %s", err))
	}

	if err := json.Unmarshal(agentConfigJSON, &agentConfig); err != nil {
		c.UI.Error(fmt.Sprintf("Error parsing environment variables: %s", err))
		return 1
	}

	level, err := common.GetLogLevel(c.flagLogLevel)
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error setting log level: %s", err))
		return 1
	}

	logger := hclog.New(&hclog.LoggerOptions{
		Name:       "handler",
		Level:      level,
		JSONFormat: (c.flagLogFormat == "json")})

	logger.Info("Starting the agent...")

	agent, err := agent.New(agentConfig, logger)
	if err != nil {
		c.UI.Error(fmt.Sprintf("Error creating an agent: %s", err))
		return 1
	}

	if err := agent.Retrieve(); err != nil {
		c.UI.Error(fmt.Sprintf("Error occurred while retrieving secrets: %s", err))
	}

	return 0
}

func (c *Command) Synopsis() string { return synopsis }
func (c *Command) Help() string {
	c.once.Do(c.init)
	return c.help
}

const synopsis = "Get secrets and store them in a mount path"
const help = `
Usage: keyvault-k8s webhook [options]
  Run the Admission Webhook server for injecting Azure KeyVault Agent containers into pods.
`
