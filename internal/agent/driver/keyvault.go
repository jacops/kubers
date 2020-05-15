package driver

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/keyvault/2016-10-01/keyvault"
	"github.com/Azure/azure-sdk-for-go/services/keyvault/auth"
	"github.com/Azure/go-autorest/autorest/azure"
)

const (
	keyVaultDriverName = "keyvault"
)

// KeyVaultDriver will get secrets from Azure KeyVault
type KeyVaultDriver struct {
	config Config

	client keyvault.BaseClient
}

// GetSecret is a main function of the driver
func (kvd *KeyVaultDriver) GetSecret(ctx context.Context, secretURL string) (string, error) {
	var secretValue string

	keyVaultName, keyName, _ := getServiceKeyNamesPair(secretURL)

	secretBundle, err := kvd.client.GetSecret(ctx, getVaultURL(keyVaultName), keyName, "")
	if err != nil {
		return secretValue, err
	}

	kvd.config.Logger.Debug(fmt.Sprintf("Fetched secret %s from %s", keyName, keyVaultName))

	return *secretBundle.Value, nil
}

func newKeyVaultDriver(config Config) *KeyVaultDriver {
	authorizer, _ := auth.NewAuthorizerFromEnvironment()

	keyvaultDriver := &KeyVaultDriver{
		config: config,
		client: keyvault.New(),
	}
	keyvaultDriver.client.Authorizer = authorizer

	return keyvaultDriver
}

func getVaultURL(keyVaultName string) string {
	return fmt.Sprintf("https://%s.%s", keyVaultName, azure.PublicCloud.KeyVaultDNSSuffix)
}
