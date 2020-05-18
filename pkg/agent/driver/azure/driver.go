package azure

import (
	"context"
	"fmt"

	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/hashicorp/go-hclog"
	"github.com/jacops/kubers/pkg/driver"
)

type service interface {
	getSecret(ctx context.Context) (string, error)
}

// Driver holds logic for retrieving secrets from Azure services
type Driver struct {
	config *driver.Config
	logger hclog.Logger
}

//New returns Azure driver
func New(config *driver.Config, logger hclog.Logger) *Driver {
	return &Driver{
		config: config,
		logger: logger,
	}
}

// GetSecret is a main function of the driver
func (d *Driver) GetSecret(ctx context.Context, secretURL string) (string, error) {
	var service service
	var secret string

	serviceType := driver.GetServiceTypeFromURL(secretURL)

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
