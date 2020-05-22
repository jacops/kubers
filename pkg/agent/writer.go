package agent

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/hashicorp/go-hclog"
)

type writeFile func(filename string, data []byte, perm os.FileMode) error

///// MOUNT PATH WRITER \\\\\\

// MountPathWriter is a main writer
type MountPathWriter struct {
	logger    hclog.Logger
	writeFile writeFile
}

// NewMountPathWriter is a factory for MountPathWriter
func NewMountPathWriter(logger hclog.Logger) *MountPathWriter {
	return &MountPathWriter{
		logger:    logger,
		writeFile: ioutil.WriteFile,
	}
}

// WriteSecret is a main method for making a secret available to the other container
func (w *MountPathWriter) WriteSecret(secret *Secret) error {
	fullPath := getFullPath(secret)
	if err := w.writeFile(fullPath, []byte(secret.Value), 0644); err != nil {
		return err
	}

	w.logger.Debug(fmt.Sprintf("Secret %s saved to %s", secret.Name, fullPath))
	return nil
}

////////////////////////////////

func getFullPath(secret *Secret) string {
	return secret.MountPath + "/" + secret.Name
}
