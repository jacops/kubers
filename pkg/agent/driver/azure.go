package driver

import (
	"context"
	"fmt"
	"net/url"

	"github.com/Azure/azure-sdk-for-go/services/keyvault/2016-10-01/keyvault"
	"github.com/Azure/azure-sdk-for-go/services/keyvault/auth"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/hashicorp/go-hclog"
)

// AzureDriver holds logic for retrieving secrets from Azure services
type AzureDriver struct {
	config *Config
	logger hclog.Logger
}

func newAzureDriver(config *Config, logger hclog.Logger) *AzureDriver {
	return &AzureDriver{
		config: config,
		logger: logger,
	}
}

// GetSecret is a main function of the driver
func (azureDriver *AzureDriver) GetSecret(ctx context.Context, secretURL string) (string, error) {
	serviceType := getServiceTypeFromURL(secretURL)

	switch serviceType {
	case "keyvault":
		return azureDriver.getSecretFromKeyVault(ctx, secretURL)
	}

	return "", fmt.Errorf("Unrecognized service: %s", serviceType)
}

func (azureDriver *AzureDriver) getSecretFromKeyVault(ctx context.Context, secretURL string) (string, error) {
	var secretValue string

	authorizer, _ := auth.NewAuthorizerFromEnvironment()
	client := keyvault.New()
	client.Authorizer = authorizer
	keyVaultName, keyName, _ := getServiceKeyNamesPair(secretURL)

	secretBundle, err := client.GetSecret(ctx, getVaultURL(keyVaultName), keyName, "")
	if err != nil {
		return secretValue, err
	}

	azureDriver.logger.Debug(fmt.Sprintf("Fetched secret %s from %s", keyName, keyVaultName))

	return *secretBundle.Value, nil
}

func getServiceTypeFromURL(secretURL string) string {
	u, err := url.Parse(secretURL)
	if err != nil {
		return ""
	}

	return u.Scheme
}

func getVaultURL(keyVaultName string) string {
	return fmt.Sprintf("https://%s.%s", keyVaultName, azure.PublicCloud.KeyVaultDNSSuffix)
}
