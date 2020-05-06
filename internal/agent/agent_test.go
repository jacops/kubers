package agent

import (
	"testing"

	"github.com/hashicorp/go-hclog"
	"github.com/jacops/azure-keyvault-k8s/internal/agent/driver"
)

func TestAgentCanRetrieveSecrets(t *testing.T) {
	tests := []struct {
		name  string
		agent Agent
	}{
		{name: "NoSecrets", agent: *getBasicAgent(Config{})},
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

func getBasicAgent(config Config) *Agent {
	return &Agent{
		logger: getLogger(),
		config: config,
	}
}

func getBasicAgentWithSecret() *Agent {

	config := Config{
		Secrets: []*Secret{
			{
				Name:      "test",
				URL:       "dummy://test/test",
				MountPath: "/mount/path",
			},
		},
	}

	return &Agent{
		logger:         getLogger(),
		config:         config,
		getDriverByURL: getOnlyDummySecretsDriverFromMapByURL,
		writer:         NewDummyPathWriter(),
	}
}

type DummyDriver struct{}
type DummyPathWriter struct{}

func (dd *DummyDriver) RetrieveSecret(secretURL string) (string, error) {
	return "dummy-secret", nil
}

func getOnlyDummySecretsDriverFromMapByURL(secretURL string, driverConfig driver.Config) (SecretsDriver, error) {
	return &DummyDriver{}, nil
}

func (w *DummyPathWriter) WriteSecret(value string, metadata *Secret) error {
	return nil
}

func NewDummyPathWriter() *DummyPathWriter {
	return &DummyPathWriter{}
}

// func (r *DummyRetriever) retrieveSecret(secretURL string) (string, error) {
// 	return "", nil
// }
