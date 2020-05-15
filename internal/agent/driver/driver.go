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
type Config struct {
	Logger hclog.Logger
}

// Factory is a struct responsible for creating drivers
type Factory interface {
	GetDriver(name string, config Config) (Driver, error)
}

//SimpleFactory is a factory based on driver name
type SimpleFactory struct {
}

// NewSimpleFactory returns SecretURLDriverFactory
func NewSimpleFactory() *SimpleFactory {
	return &SimpleFactory{}
}

// GetDriver returns Driver struct by name
func (df *SimpleFactory) GetDriver(name string, config Config) (Driver, error) {
	var secretDriver Driver

	switch name {
	case keyVaultDriverName:
		return newKeyVaultDriver(config), nil
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
