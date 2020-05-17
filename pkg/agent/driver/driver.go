package driver

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/hashicorp/go-hclog"
)

// Driver is an interfce that needs to be implemented by drivers
type Driver interface {
	GetSecret(ctx context.Context, secretURL string) (string, error)
}

// Config is a main driver configuration
type Config map[string]string

// New returns Driver struct by name
func New(name string, config *Config, logger hclog.Logger) (Driver, error) {
	var secretDriver Driver

	switch name {
	case "azure":
		return newAzureDriver(config, logger), nil
	}

	return secretDriver, fmt.Errorf("Driver not found: %s", name)
}

func getServiceKeyNamesPair(secretURL string) (serviceName string, keyName string, err error) {
	u, err := url.Parse(secretURL)
	if err != nil {
		return
	}

	serviceName = strings.Trim(u.Host, "/")
	keyName = strings.Trim(u.Path, "/")

	return
}
