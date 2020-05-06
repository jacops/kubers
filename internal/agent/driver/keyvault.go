package driver

import (
	"context"
	"fmt"
	"net/url"
	"sync"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/keyvault/2016-10-01/keyvault"
	"github.com/Azure/azure-sdk-for-go/services/keyvault/auth"
	"github.com/Azure/go-autorest/autorest/azure"
)

const (
	// KeyVaultScheme is a scheme for this retriever
	KeyVaultScheme = "keyvault"
)

// KeyVaultDriver will get secrets from Azure KeyVault
type KeyVaultDriver struct {
	config Config

	client keyvault.BaseClient
	once   sync.Once
}

func NewKeyVaultDriver(config Config) *KeyVaultDriver {
	return &KeyVaultDriver{
		config: config,
	}
}

func (kvd *KeyVaultDriver) init() {
	authorizer, _ := auth.NewAuthorizerFromEnvironment()

	kvd.client = keyvault.New()
	kvd.client.Authorizer = authorizer
}

func (kvd *KeyVaultDriver) RetrieveSecret(secretURL string) (string, error) {
	kvd.once.Do(kvd.init)
	var secretValue string

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	u, err := url.Parse(secretURL)
	if err != nil {
		return secretValue, err
	}

	if u.Scheme != KeyVaultScheme {
		return secretValue, fmt.Errorf("Invalid KeyVault scheme in the url: %s", u.Scheme)
	}

	keyVaultName := u.Host
	keyName := u.Path

	vaultBaseURL := getVaultURL(keyVaultName)

	secretBundle, err := kvd.client.GetSecret(ctx, vaultBaseURL, keyName, "")
	if err != nil {
		return secretValue, err
	}

	kvd.config.logger.Debug(fmt.Sprintf("Fetched secret %s from %s", keyName, keyVaultName))

	return *secretBundle.Value, nil
}

func getVaultURL(keyVaultName string) string {
	return fmt.Sprintf("https://%s.%s", keyVaultName, azure.PublicCloud.KeyVaultDNSSuffix)
}
