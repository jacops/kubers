package agent

import (
	"context"

	"github.com/hashicorp/go-hclog"
)

// SecretsWriter is an interface for writing secrets
type SecretsWriter interface {
	WriteSecret(secret *Secret) error
}

// Agent struct represnts the init continer
type Agent struct {
	logger   hclog.Logger
	pipeline *Pipeline
	secrets  []*Secret
}

// New returns new Agent with writeSecretToMountPath
func New(config *Config, logger hclog.Logger) (*Agent, error) {
	writer := NewMountPathWriter(logger)
	provider, err := newProvider(config.ProviderName, config.ProviderConfig, logger)
	if err != nil {
		return nil, err
	}

	pipeline := NewPipeline(writer, provider, logger, 3)

	return &Agent{
		logger:   logger,
		pipeline: pipeline,
		secrets:  config.Secrets,
	}, nil
}

// Process will get secrets and save them in specified location
func (a *Agent) Process(ctx context.Context) error {
	return a.pipeline.Do(ctx, a.secrets)
}
