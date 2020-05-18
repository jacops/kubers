package injector

import (
	"testing"

	"github.com/hashicorp/vault/sdk/helper/strutil"
)

func TestContainerEnvs(t *testing.T) {

	tests := []struct {
		agent        AgentInjector
		expectedEnvs []string
	}{
		{AgentInjector{Annotations: map[string]string{AnnotationAgentProvider: "azure"}}, []string{"AGENT_CONFIG"}},
		{AgentInjector{Annotations: map[string]string{AnnotationAgentProviderAzureCredentialsSecret: "secret", AnnotationAgentProvider: "azure"}}, []string{"AGENT_CONFIG", "AZURE_TENANT_ID", "AZURE_SUBSCRIPTION_ID", "AZURE_CLIENT_ID", "AZURE_CLIENT_SECRET"}},
	}
	for _, tt := range tests {
		envs, err := tt.agent.ContainerEnvVars(true)

		if err != nil {
			t.Errorf("got error, shouldn't have: %s", err)
		}

		if len(envs) != len(tt.expectedEnvs) {
			t.Errorf("number of envs mismatch, wanted %d, got %d", len(tt.expectedEnvs), len(envs))
		}

		for _, env := range envs {
			if !strutil.StrListContains(tt.expectedEnvs, env.Name) {
				t.Errorf("unexpected env found %s", env.Name)
			}
		}
	}
}
