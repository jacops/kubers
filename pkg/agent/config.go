package agent

import (
	"encoding/json"

	"github.com/jacops/kubers/pkg/provider"
)

// Config is the top level struct that composes am Agent
// configuration file.
type Config struct {
	Secrets        []*Secret `json:"secret"`
	ProviderName   string            `json:"provider_name"`
	ProviderConfig *provider.Config  `json:"provider_config"`
}

// Secret is a metadata object sued for fetching and storing secret
type Secret struct {
	// Name of the secret used as the filename for the rendered secret file.
	Name string `json:"name"`

	// URL of the secret e.g. keyvault://name/key
	URL string `json:"url"`

	// Mount Path
	MountPath string `json:"mount_path"`

	Value string `json:"value"`
}

// Render is a method used for converting config to json
func (c *Config) Render() ([]byte, error) {
	return json.Marshal(c)
}
