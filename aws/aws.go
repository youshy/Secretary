package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
)

type Aws struct {
	session *secretsmanager.SecretsManager
}

func NewAwsProvider(session *secretsmanager.SecretsManager) Aws {
	return Aws{
		session,
	}
}

func (a *Aws) GetSecret(secret string) (string, error) {
	input := secretsmanager.GetSecretValueInput{
		SecretId: aws.String(secret),
	}

	output, err := a.session.GetSecretValue(&input)
	if err != nil {
		return "", err
	}

	return aws.StringValue(output.SecretString), nil
}

func (a *Aws) PutSecret(key, secret string) error {
	input := secretsmanager.PutSecretValueInput{
		SecretId:     aws.String(key),
		SecretString: aws.String(secret),
	}

	_, err := a.session.PutSecretValue(&input)
	return err
}
