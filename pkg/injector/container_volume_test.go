package injector

import (
	"testing"

	"github.com/jacops/kubers/pkg/agent"
)

func TestVolumes(t *testing.T) {
	tests := []struct {
		name            string
		injector        AgentInjector
		expectedVolumes int
	}{
		{name: "no secrets", expectedVolumes: 1, injector: AgentInjector{}},
		{name: "no secrets", expectedVolumes: 1, injector: AgentInjector{
			Secrets: []*agent.Secret{{Name: "secret"}},
		}},
		{name: "no secrets", expectedVolumes: 2, injector: AgentInjector{
			Secrets: []*agent.Secret{{Name: "secret", MountPath: "/path"}},
		}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			volumes := tt.injector.ContainerVolumes()
			mounts := tt.injector.ContainerVolumeMounts()

			if len(volumes) != tt.expectedVolumes {
				t.Errorf("Invalid number of volumes. Got: %d, expected: %d", len(volumes), tt.expectedVolumes)
			}

			if len(mounts) != tt.expectedVolumes {
				t.Errorf("Invalid number of mounts. Got: %d, expected: %d", len(mounts), tt.expectedVolumes)
			}
		})
	}
}
