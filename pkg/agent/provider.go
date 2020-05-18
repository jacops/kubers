package agent

import (
	"fmt"

	"github.com/hashicorp/go-hclog"
	awsProvider "github.com/jacops/kubers/pkg/agent/provider/aws"
	azureProvider "github.com/jacops/kubers/pkg/agent/provider/azure"
	"github.com/jacops/kubers/pkg/provider"
)

func newProvider(name string, config *provider.Config, logger hclog.Logger) (provider.Provider, error) {
	switch name {
	case "azure":
		return azureProvider.New(config, logger), nil
	case "aws":
		return awsProvider.New(config, logger), nil
	}

	return nil, fmt.Errorf("Provider not found: %s", name)
}
