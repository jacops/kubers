package agent

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/jacops/kubers/pkg/driver"
)

// Agent struct represnts the init continer
type Agent struct {
	config *Config
	logger hclog.Logger
	writer SecretsWriter
	driver driver.Driver
}

// New returns new Agent with writeSecretToMountPath
func New(config *Config, logger hclog.Logger) (*Agent, error) {
	secretDriver, err := newDriver(config.DriverName, config.DriverConfig, logger)
	if err != nil {
		return nil, err
	}

	return &Agent{
		config: config,
		logger: logger,
		writer: NewMountPathWriter(logger),
		driver: secretDriver,
	}, nil
}

// Retrieve will get secrets and save them in specified location
func (a *Agent) Retrieve() (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	errorChannel := make(chan error, len(a.config.Secrets))

	for _, secretMetadata := range a.config.Secrets {
		go func(secretMetadata *SecretMetadata) (err error) {
			defer func() {
				errorChannel <- err
			}()

			secret, err := a.driver.GetSecret(ctx, secretMetadata.URL)
			if err != nil {
				return
			}

			err = a.writer.WriteSecret(secret, secretMetadata)

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
