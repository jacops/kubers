package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/hashicorp/go-hclog"
	"github.com/jacops/kubers/pkg/provider"
)

type service interface {
	getSecret(ctx context.Context) (string, error)
}

// Provider holds logic for retrieving secrets from AWS services
type Provider struct {
	region string

	logger hclog.Logger
}

//New returns AWS provider
func New(config *provider.Config, logger hclog.Logger) *Provider {
	return &Provider{
		region: (*config)["region"],
		logger: logger,
	}
}

// GetSecret is a main function of the provider
func (d *Provider) GetSecret(ctx context.Context, secretURL string) (string, error) {
	var service service
	var secret string

	serviceType := provider.GetServiceTypeFromURL(secretURL)
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(d.region),
	}))

	switch serviceType {
	case "secretsmanager":
		service = newSecretsManager(secretURL, sess)
	default:
		return secret, fmt.Errorf("Unrecognized service: %s", serviceType)
	}

	return service.getSecret(ctx)
}
