package agent

import (
	"encoding/json"

	"github.com/jacops/kubers/pkg/provider"
)

// Config is the top level struct that composes am Agent
// configuration file.
type Config struct {
	Secrets      []*SecretMetadata `json:"secret"`
	ProviderName   string            `json:"provider_name"`
	ProviderConfig *provider.Config    `json:"provider_config"`
}

// SecretMetadata is a metadata object sued for fetching and storing secret
type SecretMetadata struct {
	// Name of the secret used as the filename for the rendered secret file.
	Name string `json:"name"`

	// URL of the secret e.g. keyvault://name/key
	URL string `json:"url"`

	// Mount Path
	MountPath string `json:"mount_path"`
}

// Render is a method used for converting config to json
func (c *Config) Render() ([]byte, error) {
	return json.Marshal(c)
}
