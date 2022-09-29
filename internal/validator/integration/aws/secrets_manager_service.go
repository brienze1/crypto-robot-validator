package aws

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/brienze1/crypto-robot-validator/internal/validator/domain/adapters"
	adapters2 "github.com/brienze1/crypto-robot-validator/internal/validator/integration/adapters"
	"github.com/brienze1/crypto-robot-validator/internal/validator/integration/exceptions"
	"github.com/brienze1/crypto-robot-validator/pkg/custom_error"
)

type secretsManagerService struct {
	logger         adapters.LoggerAdapter
	secretsManager adapters2.SecretsManagerAdapter
}

// SecretsManagerService constructor method, used to inject dependencies.
func SecretsManagerService(logger adapters.LoggerAdapter, secretsManager adapters2.SecretsManagerAdapter) *secretsManagerService {
	return &secretsManagerService{
		logger:         logger,
		secretsManager: secretsManager,
	}
}

// GetSecret is used to retrieve secrets from secrets manager, returns *dto.Secrets.
func (s *secretsManagerService) GetSecret(secretName string, secretObject any) custom_error.BaseErrorAdapter {
	s.logger.Info("Get secret starting", secretName)

	result, err := s.secretsManager.GetSecretValue(context.TODO(), &secretsmanager.GetSecretValueInput{SecretId: aws.String(secretName)})
	if err != nil {
		return s.abort(err, "error while getting secret")
	}

	var secretString, decodedBinarySecret string
	if result.SecretString != nil {
		secretString = *result.SecretString
		err := json.Unmarshal([]byte(secretString), secretObject)
		if err != nil {
			return s.abort(err, "error while unmarshalling secret string")
		}
	} else {
		decodedBinarySecretBytes := make([]byte, base64.StdEncoding.DecodedLen(len(result.SecretBinary)))
		decodedLen, err := base64.StdEncoding.Decode(decodedBinarySecretBytes, result.SecretBinary)
		if err != nil {
			return s.abort(err, "error while decoding secret binary")
		}
		decodedBinarySecret = string(decodedBinarySecretBytes[:decodedLen])
		err = json.Unmarshal([]byte(decodedBinarySecret), secretObject)
		if err != nil {
			return s.abort(err, "error while unmarshalling secret binary")
		}
	}

	s.logger.Info("Get secret finished", secretName)
	return nil
}

func (s *secretsManagerService) abort(err error, message string) custom_error.BaseErrorAdapter {
	secretsManagerError := exceptions.SecretsManagerError(err, message)
	s.logger.Error(secretsManagerError, "Get secret failed: "+message)
	return secretsManagerError
}
