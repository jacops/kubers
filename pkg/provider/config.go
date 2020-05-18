package provider

import "context"

// Provider is an interfce that needs to be implemented by providers
type Provider interface {
	GetSecret(ctx context.Context, secretURL string) (string, error)
}

// Config is a main provider configuration
type Config map[string]string
