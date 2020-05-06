package agent

import (
	"fmt"
	"net/url"

	"github.com/jacops/azure-keyvault-k8s/internal/agent/driver"
)

type getSecretsDriverByURL = func(secretURL string, driverConfig driver.Config) (SecretsDriver, error)

// SecretsDriver is an interfce that needs to be implemented by drivers
type SecretsDriver interface {
	RetrieveSecret(secretURL string) (string, error)
}

func getSecretsDriverFromMapByURL(secretURL string, driverConfig driver.Config) (SecretsDriver, error) {
	var secretsDriver SecretsDriver

	u, err := url.Parse(secretURL)
	if err != nil {
		return secretsDriver, err
	}

	switch u.Scheme {
	case driver.KeyVaultScheme:
		secretsDriver = driver.NewKeyVaultDriver(driverConfig)
	default:
		return secretsDriver, fmt.Errorf("Unable to find a driver: %s", u.Scheme)
	}

	return secretsDriver, nil
}
