package agent

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/jacops/kubers/internal/agent/driver"
)

// Agent struct represnts the init continer
type Agent struct {
	config        Config
	logger        hclog.Logger
	writer        SecretsWriter
	driverFactory driver.Factory
}

// New returns new Agent with writeSecretToMountPath
func New(config Config, logger hclog.Logger) *Agent {
	return &Agent{
		config:        config,
		logger:        logger,
		writer:        NewMountPathWriter(logger),
		driverFactory: driver.NewSimpleFactory(),
	}
}

func (a *Agent) getDriver(name string) (driver.Driver, error) {
	return a.driverFactory.GetDriver(name, getDefaultDriverConfig(a.logger))
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

			secretDriver, err := a.getDriver(getDriverNameByURL(secretMetadata.URL))
			if err != nil {
				return
			}

			secret, err := secretDriver.GetSecret(ctx, secretMetadata.URL)
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

func getDefaultDriverConfig(logger hclog.Logger) driver.Config {
	return driver.Config{
		Logger: logger,
	}
}

func getDriverNameByURL(secretURL string) string {
	u, _ := url.Parse(secretURL)
	return u.Scheme
}
