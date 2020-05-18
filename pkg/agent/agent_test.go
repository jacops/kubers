package agent

import (
	"context"
	"testing"

	"github.com/hashicorp/go-hclog"
	"github.com/jacops/kubers/pkg/provider"
)

func TestAgentCanRetrieveSecrets(t *testing.T) {
	tests := []struct {
		name  string
		agent Agent
	}{
		{name: "NoSecrets", agent: *getBasicAgent(&Config{})},
		{name: "WithSecretsAndDummyWriter", agent: *getBasicAgentWithSecret()},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.agent.Retrieve()
			if err != nil {
				t.Errorf("did not expect any errors: %s", err)
			}
		})
	}
}

func getLogger() hclog.Logger {
	return hclog.New(&hclog.LoggerOptions{
		Name: "handler",
	})
}

func getBasicAgent(config *Config) *Agent {
	return &Agent{
		logger: getLogger(),
		config: config,
	}
}

func getBasicAgentWithSecret() *Agent {

	config := &Config{
		Secrets: []*SecretMetadata{
			{
				Name:      "test",
				URL:       "dummy://test/test",
				MountPath: "/mount/path",
			},
		},
	}

	return &Agent{
		logger: getLogger(),
		config: config,
		provider: &DummyProvider{},
		writer: NewDummyPathWriter(),
	}
}

type DummyProvider struct{}
type DummyPathWriter struct{}

func (dd *DummyProvider) GetSecret(ctx context.Context, secretURL string) (string, error) {
	return "dummy-secret", nil
}

func getOnlyDummySecretsProviderFromMapByURL(secretURL string, providerConfig provider.Config) (provider.Provider, error) {
	return &DummyProvider{}, nil
}

func (w *DummyPathWriter) WriteSecret(value string, metadata *SecretMetadata) error {
	return nil
}

func NewDummyPathWriter() *DummyPathWriter {
	return &DummyPathWriter{}
}
