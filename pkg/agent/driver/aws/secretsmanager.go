package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/jacops/kubers/pkg/driver"
)

type secretsManager struct {
	key     string
	session *session.Session
}

func newSecretsManager(secretURL string, sess *session.Session) *secretsManager {
	keyName, _ := driver.GetKeyName(secretURL)

	return &secretsManager{
		key:     keyName,
		session: sess,
	}
}

func (sm *secretsManager) getSecret(ctx context.Context) (string, error) {
	input := &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(sm.key),
	}

	svc := secretsmanager.New(sm.session)
	result, err := svc.GetSecretValue(input)
	if err != nil {
		return "", err
	}

	return *result.SecretString, nil
}
