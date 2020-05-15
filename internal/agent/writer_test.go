package agent

import (
	"os"
	"testing"

	"github.com/hashicorp/go-hclog"
)

func TestMountPathWriterIsInitialized(t *testing.T) {
	t.Run("MountPathWriter", func(t *testing.T) {
		logger := hclog.New(&hclog.LoggerOptions{})
		mountPathWriter := NewMountPathWriter(logger)

		if mountPathWriter.logger != logger {
			t.Errorf("Invalid logger returned")
		}
	})
}

func TestMountPathWriterSecretIsOperatingOnRightValues(t *testing.T) {
	mountPathWriter := &MountPathWriter{
		logger: hclog.New(&hclog.LoggerOptions{}),
		writeFile: func(filename string, data []byte, perm os.FileMode) error {
			return nil
		},
	}

	tests := []struct {
		name     string
		writer   SecretsWriter
		metadata *SecretMetadata
		secret   string
	}{
		{writer: mountPathWriter, secret: "changeme", metadata: &SecretMetadata{}},
		{writer: mountPathWriter, secret: "somesecret", metadata: &SecretMetadata{}},
		{writer: mountPathWriter, secret: "#@#F##F#F", metadata: &SecretMetadata{}},
	}

	for _, tt := range tests {
		t.Run(tt.secret, func(t *testing.T) {
			err := tt.writer.WriteSecret(tt.secret, tt.metadata)

			if err != nil {
				t.Errorf("did not expect any errors: %s", err)
			}
		})
	}
}
