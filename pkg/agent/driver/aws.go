package driver

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/hashicorp/go-hclog"
)

// AWSDriver holds logic for retrieving secrets from AWS services
type AWSDriver struct {
	region string

	logger hclog.Logger
}

func newAWSDriver(config *Config, logger hclog.Logger) *AWSDriver {
	if _, ok := (*config)["region"]; !ok {
		logger.Info("nie ma")
	}

	return &AWSDriver{
		region: (*config)["region"],
		logger: logger,
	}
}

// GetSecret is a main function of the driver
func (awsDriver *AWSDriver) GetSecret(ctx context.Context, secretURL string) (string, error) {
	serviceType := getServiceTypeFromURL(secretURL)

	switch serviceType {
	case "secretmanager":
		return awsDriver.getSecretFromSecretManager(ctx, secretURL)
	}

	return "", fmt.Errorf("Unrecognized service: %s", serviceType)
}

func (awsDriver *AWSDriver) getSecretFromSecretManager(ctx context.Context, secretURL string) (string, error) {
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(awsDriver.region),
	}))

	keyName, _ := getKeyName(secretURL)

	svc := secretsmanager.New(sess)
	input := &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(keyName),
	}

	result, err := svc.GetSecretValue(input)
	if err != nil {
		return "", err
	}

	awsDriver.logger.Debug(fmt.Sprintf("Fetched secret %s from Secret Manager", keyName))

	return *result.SecretString, nil
}
