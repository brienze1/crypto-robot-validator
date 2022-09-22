package aws

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/brienze1/crypto-robot-validator/internal/validator/domain/adapters"
	adapters2 "github.com/brienze1/crypto-robot-validator/internal/validator/integration/adapters"
	"github.com/brienze1/crypto-robot-validator/internal/validator/integration/dto"
	"github.com/brienze1/crypto-robot-validator/internal/validator/integration/exceptions"
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
func (s *secretsManagerService) GetSecret(secretName string) (*dto.Secrets, error) {
	s.logger.Info("Get secret starting", secretName)

	result, err := s.secretsManager.GetSecretValue(context.TODO(), &secretsmanager.GetSecretValueInput{SecretId: aws.String(secretName)})
	if err != nil {
		return nil, s.abort(err, "error while getting secret")
	}

	var secrets dto.Secrets
	var secretString, decodedBinarySecret string
	if result.SecretString != nil {
		secretString = *result.SecretString
		err := json.Unmarshal([]byte(secretString), &secrets)
		if err != nil {
			return nil, s.abort(err, "error while unmarshalling secret string")
		}
	} else {
		decodedBinarySecretBytes := make([]byte, base64.StdEncoding.DecodedLen(len(result.SecretBinary)))
		decodedLen, err := base64.StdEncoding.Decode(decodedBinarySecretBytes, result.SecretBinary)
		if err != nil {
			return nil, s.abort(err, "error while decoding secret binary")
		}
		decodedBinarySecret = string(decodedBinarySecretBytes[:decodedLen])
		err = json.Unmarshal([]byte(decodedBinarySecret), &secrets)
		if err != nil {
			return nil, s.abort(err, "error while unmarshalling secret binary")
		}
	}

	s.logger.Info("Get secret finished", secretName)
	return &secrets, nil
}

func (s *secretsManagerService) abort(err error, message string) error {
	secretsManagerError := exceptions.SecretsManagerError(err, message)
	s.logger.Error(secretsManagerError, "Get secret failed: "+message)
	return secretsManagerError
}
