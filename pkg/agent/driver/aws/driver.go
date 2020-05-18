package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/hashicorp/go-hclog"
	"github.com/jacops/kubers/pkg/driver"
)

type service interface {
	getSecret(ctx context.Context) (string, error)
}

// Driver holds logic for retrieving secrets from AWS services
type Driver struct {
	region string

	logger hclog.Logger
}

//New returns AWS driver
func New(config *driver.Config, logger hclog.Logger) *Driver {
	return &Driver{
		region: (*config)["region"],
		logger: logger,
	}
}

// GetSecret is a main function of the driver
func (d *Driver) GetSecret(ctx context.Context, secretURL string) (string, error) {
	var service service
	var secret string

	serviceType := driver.GetServiceTypeFromURL(secretURL)
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
