package agent

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-hclog"
	awsProvider "github.com/jacops/kubers/pkg/agent/provider/aws"
	azureProvider "github.com/jacops/kubers/pkg/agent/provider/azure"
	"github.com/jacops/kubers/pkg/provider"
)

// Provider is an interfce that needs to be implemented by providers
type Provider interface {
	GetSecret(ctx context.Context, secretURL string) (string, error)
}

func newProvider(name string, config *provider.Config, logger hclog.Logger) (Provider, error) {
	switch name {
	case "azure":
		return azureProvider.New(config, logger), nil
	case "aws":
		return awsProvider.New(config, logger), nil
	}

	return nil, fmt.Errorf("Provider not found: %s", name)
}
