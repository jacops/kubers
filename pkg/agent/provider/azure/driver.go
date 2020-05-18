package azure

import (
	"context"
	"fmt"

	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/hashicorp/go-hclog"
	"github.com/jacops/kubers/pkg/provider"
)

type service interface {
	getSecret(ctx context.Context) (string, error)
}

// Provider holds logic for retrieving secrets from Azure services
type Provider struct {
	config *provider.Config
	logger hclog.Logger
}

//New returns Azure provider
func New(config *provider.Config, logger hclog.Logger) *Provider {
	return &Provider{
		config: config,
		logger: logger,
	}
}

// GetSecret is a main function of the provider
func (d *Provider) GetSecret(ctx context.Context, secretURL string) (string, error) {
	var service service
	var secret string

	serviceType := provider.GetServiceTypeFromURL(secretURL)

	switch serviceType {
	case "keyvault":
		service = newKeyVault(secretURL)
	default:
		return secret, fmt.Errorf("Unrecognized service: %s", serviceType)
	}

	return service.getSecret(ctx)
}

func getVaultURL(keyVaultName string) string {
	return fmt.Sprintf("https://%s.%s", keyVaultName, azure.PublicCloud.KeyVaultDNSSuffix)
}
