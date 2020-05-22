package agent

import (
	"encoding/json"
	"testing"
)

func TestConfigCanBeMarshalledToJson(t *testing.T) {
	tests := []struct {
		name                    string
		config                  Config
		expectedNumberOfSecrets int
	}{
		{name: "withnosecrets", config: Config{}, expectedNumberOfSecrets: 0},
		{name: "withnosecrets", config: Config{Secrets: []*Secret{{Name: "dummy-name"}}}, expectedNumberOfSecrets: 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonConfig, err := tt.config.Render()
			var unmarshalledConfig Config
			json.Unmarshal(jsonConfig, &unmarshalledConfig)
			secretsNumber := len(unmarshalledConfig.Secrets)

			if secretsNumber != tt.expectedNumberOfSecrets {
				t.Errorf("Invalid number of secrets after unmarshalling: %d. Expected: %d", secretsNumber, tt.expectedNumberOfSecrets)
			}

			if err != nil {
				t.Errorf("did not expect any errors: %s", err)
			}
		})
	}
}
