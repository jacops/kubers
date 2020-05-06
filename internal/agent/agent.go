package agent

import (
	"fmt"

	"github.com/hashicorp/go-hclog"
	"github.com/jacops/azure-keyvault-k8s/internal/agent/driver"
)

// Agent struct represnts the init continer
type Agent struct {
	config Config
	logger hclog.Logger
	writer SecretsWriter

	getDriverByURL getSecretsDriverByURL
}

// NewAgent returns new Agent with writeSecretToMountPath
func NewAgent(config Config, logger hclog.Logger) *Agent {
	return &Agent{
		config:         config,
		logger:         logger,
		writer:         NewMountPathWriter(logger),
		getDriverByURL: getSecretsDriverFromMapByURL,
	}
}

// Retrieve will get secrets and save them in specified location
func (a *Agent) Retrieve() (err error) {
	errorChannel := make(chan error, len(a.config.Secrets))

	driverConfig := driver.NewConfig(a.logger)

	for _, secretMetadata := range a.config.Secrets {
		go func(secretMetadata *Secret) (err error) {
			defer func() {
				errorChannel <- err
			}()

			driver, err := a.getDriverByURL(secretMetadata.URL, driverConfig)
			if err != nil {
				return
			}

			secret, err := driver.RetrieveSecret(secretMetadata.URL)
			if err != nil {
				return
			}

			if err = a.writer.WriteSecret(secret, secretMetadata); err != nil {
				return
			}

			return
		}(secretMetadata)
	}

	a.logger.Debug(fmt.Sprintf("Waiting for secrets..."))

	for c := 0; c < cap(errorChannel); c++ {
		err = <-errorChannel
		if err != nil {
			return err
		}
	}

	a.logger.Info(fmt.Sprintf("Secrets were successfully processed..."))
	return nil
}
