package fetchsecrets

import (
	"os"
	"testing"
)

func TestCommandEnvs(t *testing.T) {
	var cmd Command
	tests := []struct {
		env    string
		value  string
		cmdPtr *string
	}{
		{env: "AGENT_LOG_LEVEL", value: "info", cmdPtr: &cmd.flagLogLevel},
		{env: "AGENT_LOG_FORMAT", value: "standard", cmdPtr: &cmd.flagLogFormat},
		{env: "AGENT_CONFIG", value: "standard", cmdPtr: &cmd.flagAgentConfig},
	}

	for _, tt := range tests {
		t.Run(tt.env, func(t *testing.T) {
			if err := os.Setenv(tt.env, tt.value); err != nil {
				t.Errorf("got error setting env, shouldn't have: %s", err)
			}
			defer os.Unsetenv(tt.env)

			if err := cmd.parseEnvs(); err != nil {
				t.Errorf("got error parsing envs, shouldn't have: %s", err)
			}

			if *tt.cmdPtr != tt.value {
				t.Errorf("env wasn't parsed, should have been: got %s, expected %s", *tt.cmdPtr, tt.value)
			}
		})
	}
}
