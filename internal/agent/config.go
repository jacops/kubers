package agent

import "encoding/json"

// Secret is a metadata struct for fetching and storing secret
type Secret struct {
	// Name of the secret used as the filename for the rendered secret file.
	Name string `json:"name"`

	// URL of the secret e.g. keyvault://name/key
	URL string `json:"url"`

	// Mount Path
	MountPath string `json:"mount_path"`
}

// Config is the top level struct that composes am Agent
// configuration file.
type Config struct {
	Secrets []*Secret `json:"secret"`
}

// Render is a method used for converting config to json
func (c *Config) Render() ([]byte, error) {
	return json.Marshal(c)
}
