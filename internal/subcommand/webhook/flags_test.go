package webhook

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
		{env: "KUBERSD_LISTEN", value: ":8080", cmdPtr: &cmd.flagListen},
		{env: "KUBERSD_TLS_KEY_FILE", value: "server.key", cmdPtr: &cmd.flagKeyFile},
		{env: "KUBERSD_TLS_CERT_FILE", value: "server.crt", cmdPtr: &cmd.flagCertFile},
		{env: "KUBERSD_TLS_AUTO_HOSTS", value: "foobar.com", cmdPtr: &cmd.flagAutoHosts},
		{env: "KUBERSD_TLS_AUTO", value: "mutationWebhook", cmdPtr: &cmd.flagAutoName},
		{env: "KUBERSD_LOG_LEVEL", value: "info", cmdPtr: &cmd.flagLogLevel},
		{env: "KUBERSD_LOG_FORMAT", value: "standard", cmdPtr: &cmd.flagLogFormat},
	}

	for _, tt := range tests {
		t.Run(tt.env, func(t *testing.T) {
			if err := os.Setenv(tt.env, tt.value); err != nil {
				t.Errorf("got error setting env, shouldn't have: %s", err)
			}
			defer os.Unsetenv(tt.env)

			if err := cmd.parseKubersDEnvs(); err != nil {
				t.Errorf("got error parsing envs, shouldn't have: %s", err)
			}

			if *tt.cmdPtr != tt.value {
				t.Errorf("env wasn't parsed, should have been: got %s, expected %s", *tt.cmdPtr, tt.value)
			}
		})
	}
}

func TestCommandAgentEnvs(t *testing.T) {
	var cmd Command
	tests := []struct {
		env    string
		value  string
		cmdPtr *string
	}{
		{env: "KUBERS_AGENT_LOG_LEVEL", value: "debug", cmdPtr: &cmd.flagAgentLogLevel},
		{env: "KUBERS_AGENT_LOG_FORMAT", value: "standard", cmdPtr: &cmd.flagAgentLogFormat},
		{env: "KUBERS_AGENT_PROVIDER", value: "aws", cmdPtr: &cmd.flagAgentProvider},
		{env: "KUBERS_AGENT_PROVIDER_AWS_REGION", value: "eu-west-2", cmdPtr: &cmd.flagAgentProviderAWSRegion},
	}

	for _, tt := range tests {
		t.Run(tt.env, func(t *testing.T) {
			if err := os.Setenv(tt.env, tt.value); err != nil {
				t.Errorf("got error setting env, shouldn't have: %s", err)
			}
			defer os.Unsetenv(tt.env)

			if err := cmd.parseKubersAgentEnvs(); err != nil {
				t.Errorf("got error parsing envs, shouldn't have: %s", err)
			}

			if *tt.cmdPtr != tt.value {
				t.Errorf("env wasn't parsed, should have been: got %s, expected %s", *tt.cmdPtr, tt.value)
			}
		})
	}
}
