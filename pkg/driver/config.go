package driver

import "context"

// Driver is an interfce that needs to be implemented by drivers
type Driver interface {
	GetSecret(ctx context.Context, secretURL string) (string, error)
}

// Config is a main driver configuration
type Config map[string]string
