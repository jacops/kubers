package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/services/keyvault/2016-10-01/keyvault"
	"github.com/Azure/azure-sdk-for-go/services/keyvault/auth"
	"github.com/jacops/kubers/pkg/provider"
)

type keyVault struct {
	name   string
	key    string
	client keyvault.BaseClient
}

func newKeyVault(secretURL string) *keyVault {
	authorizer, _ := auth.NewAuthorizerFromEnvironment()
	client := keyvault.New()
	client.Authorizer = authorizer

	keyVaultName, keyName, _ := provider.GetServiceKeyNamesPair(secretURL)

	return &keyVault{
		name:   keyVaultName,
		key:    keyName,
		client: client,
	}
}

func (kv *keyVault) getSecret(ctx context.Context) (string, error) {
	secretBundle, err := kv.client.GetSecret(ctx, getVaultURL(kv.name), kv.key, "")
	if err != nil {
		return "", err
	}

	return *secretBundle.Value, nil
}
