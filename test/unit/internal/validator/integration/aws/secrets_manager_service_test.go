package aws_test

import (
	"errors"
	"github.com/brienze1/crypto-robot-validator/internal/validator/integration/adapters"
	"github.com/brienze1/crypto-robot-validator/internal/validator/integration/aws"
	"github.com/brienze1/crypto-robot-validator/internal/validator/integration/dto"
	"github.com/brienze1/crypto-robot-validator/pkg/custom_error"
	"github.com/brienze1/crypto-robot-validator/test/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	secretsManagerService adapters.SecretsManagerServiceAdapter
	secretsManager        = mocks.SecretsManager()
	logger                = mocks.Logger()
)

var (
	secrets dto.RedisSecrets
)

func setup() {
	logger.Reset()
	secretsManager.Reset()

	secrets = dto.RedisSecrets{
		Address:    uuid.NewString(),
		Password:   uuid.NewString(),
		User:       uuid.NewString(),
		DatabaseId: 12345,
	}

	secretsManager.SetSecret("secretName", secrets)

	secretsManagerService = aws.SecretsManagerService(logger, secretsManager)
}

func TestGetSecretFromStringSuccess(t *testing.T) {
	setup()

	secret := &dto.RedisSecrets{}
	err := secretsManagerService.GetSecret("secretName", secret)

	assert.Nil(t, err)
	assert.Equal(t, &secrets, secret)
	assert.Equal(t, 1, secretsManager.GetSecretValueCounter)
	assert.Equal(t, 2, logger.InfoCallCounter)
	assert.Equal(t, 0, logger.ErrorCallCounter)
}

func TestGetSecretFromStringFailure(t *testing.T) {
	setup()

	secretsManager.SetSecret("secretName", "")

	secret := &dto.RedisSecrets{}
	err := secretsManagerService.GetSecret("secretName", secret)

	assert.Equal(t, "json: cannot unmarshal string into Go value of type dto.RedisSecrets", err.(custom_error.BaseErrorAdapter).Error())
	assert.Equal(t, "error while unmarshalling secret string", err.(custom_error.BaseErrorAdapter).InternalError())
	assert.Equal(t, "Error while trying to get secret from secrets manager", err.(custom_error.BaseErrorAdapter).Description())
	assert.Equal(t, 1, secretsManager.GetSecretValueCounter)
	assert.Equal(t, 1, logger.InfoCallCounter)
	assert.Equal(t, 1, logger.ErrorCallCounter)
}

func TestGetSecretFromBinarySuccess(t *testing.T) {
	setup()

	secretsManager.ReturnEmptyString = true

	secret := &dto.RedisSecrets{}
	err := secretsManagerService.GetSecret("secretName", secret)

	assert.Nil(t, err)
	assert.Equal(t, &secrets, secret)
	assert.Equal(t, 1, secretsManager.GetSecretValueCounter)
	assert.Equal(t, 2, logger.InfoCallCounter)
	assert.Equal(t, 0, logger.ErrorCallCounter)
}

func TestGetSecretFromBinaryUnmarshalFailure(t *testing.T) {
	setup()

	secretsManager.ReturnEmptyString = true
	secretsManager.ReturnEmptyBinary = true

	secret := &dto.RedisSecrets{}
	err := secretsManagerService.GetSecret("secretName", secret)

	assert.Equal(t, "unexpected end of JSON input", err.(custom_error.BaseErrorAdapter).Error())
	assert.Equal(t, "error while unmarshalling secret binary", err.(custom_error.BaseErrorAdapter).InternalError())
	assert.Equal(t, "Error while trying to get secret from secrets manager", err.(custom_error.BaseErrorAdapter).Description())
	assert.Equal(t, 1, secretsManager.GetSecretValueCounter)
	assert.Equal(t, 1, logger.InfoCallCounter)
	assert.Equal(t, 1, logger.ErrorCallCounter)
}

func TestGetSecretSecretsManagerFailure(t *testing.T) {
	setup()

	secretsManager.GetSecretValueError = errors.New("error test")

	secret := &dto.RedisSecrets{}
	err := secretsManagerService.GetSecret("secretName", secret)

	assert.Equal(t, "error test", err.(custom_error.BaseErrorAdapter).Error())
	assert.Equal(t, "error while getting secret", err.(custom_error.BaseErrorAdapter).InternalError())
	assert.Equal(t, "Error while trying to get secret from secrets manager", err.(custom_error.BaseErrorAdapter).Description())
	assert.Equal(t, 1, secretsManager.GetSecretValueCounter)
	assert.Equal(t, 1, logger.InfoCallCounter)
	assert.Equal(t, 1, logger.ErrorCallCounter)
}
