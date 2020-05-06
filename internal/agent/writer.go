package agent

import (
	"fmt"
	"io/ioutil"

	"github.com/hashicorp/go-hclog"
)

// SecretsWriter is an interface for writing secrets
type SecretsWriter interface {
	WriteSecret(value string, metadata *Secret) error
}

///// MOUNT PATH WRITER \\\\\\

type MountPathWriter struct {
	logger hclog.Logger
}

func (w *MountPathWriter) WriteSecret(value string, metadata *Secret) error {
	fullPath := getFullPath(metadata)
	if err := ioutil.WriteFile(fullPath, []byte(value), 0644); err != nil {
		return err
	}

	w.logger.Debug(fmt.Sprintf("Secret %s saved to %s", metadata.Name, fullPath))
	return nil
}

func NewMountPathWriter(logger hclog.Logger) *MountPathWriter {
	return &MountPathWriter{
		logger: logger,
	}
}

////////////////////////////////

func getFullPath(metadata *Secret) string {
	return metadata.MountPath + "/" + metadata.Name
}
