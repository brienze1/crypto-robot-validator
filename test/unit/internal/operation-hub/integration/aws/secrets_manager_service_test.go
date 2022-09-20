package aws_test

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	aws2 "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/brienze1/crypto-robot-operation-hub/internal/operation-hub/application/adapters"
	adapters2 "github.com/brienze1/crypto-robot-operation-hub/internal/operation-hub/domain/adapters"
	adapters3 "github.com/brienze1/crypto-robot-operation-hub/internal/operation-hub/integration/adapters"
	"github.com/brienze1/crypto-robot-operation-hub/internal/operation-hub/integration/aws"
	"github.com/brienze1/crypto-robot-operation-hub/internal/operation-hub/integration/dto"
	"github.com/brienze1/crypto-robot-operation-hub/pkg/custom_error"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

type (
	loggerMock struct {
		adapters2.LoggerAdapter
	}
	secretsManagerMock struct {
	}
)

var (
	loggerInfoCounter                   int
	loggerErrorCounter                  int
	secretsManagerGetSecretValueCounter int
	secretsManagerGetSecretValueReturn  secretsmanager.GetSecretValueOutput
	secretsManagerGetSecretValueError   error
)

func (l loggerMock) Info(string, ...interface{}) {
	loggerInfoCounter++
}

func (l loggerMock) Error(error, string, ...interface{}) {
	loggerErrorCounter++
}

func (s secretsManagerMock) GetSecretValue(_ context.Context, _ *secretsmanager.GetSecretValueInput, _ ...func(*secretsmanager.Options)) (
	*secretsmanager.GetSecretValueOutput,
	error,
) {
	secretsManagerGetSecretValueCounter++
	return &secretsManagerGetSecretValueReturn, secretsManagerGetSecretValueError
}

var (
	secretsManagerService adapters.SecretsManagerServiceAdapter
	logger                adapters2.LoggerAdapter
	secretsManager        adapters3.SecretsManagerAdapter
)

var (
	secrets dto.Secrets
)

func setup() {
	loggerInfoCounter = 0
	loggerErrorCounter = 0
	secretsManagerGetSecretValueCounter = 0
	secretsManagerGetSecretValueError = nil

	secrets = dto.Secrets{
		Host:     uuid.NewString(),
		Port:     1234123,
		User:     uuid.NewString(),
		Password: uuid.NewString(),
		DbName:   uuid.NewString(),
	}

	secretsString, _ := json.Marshal(secrets)

	encodedBinarySecretBytes := make([]byte, base64.StdEncoding.EncodedLen(len(secretsString)))
	base64.StdEncoding.Encode(encodedBinarySecretBytes, secretsString)

	secretsManagerGetSecretValueReturn = secretsmanager.GetSecretValueOutput{
		SecretBinary: encodedBinarySecretBytes,
		SecretString: aws2.String(string(secretsString)),
	}

	logger = loggerMock{}
	secretsManager = secretsManagerMock{}
	secretsManagerService = aws.SecretsManagerService(logger, secretsManager)
}

func TestGetSecretFromStringSuccess(t *testing.T) {
	setup()

	secret, err := secretsManagerService.GetSecret("secretName")

	assert.Nil(t, err)
	assert.Equal(t, &secrets, secret)
	assert.Equal(t, 1, secretsManagerGetSecretValueCounter)
	assert.Equal(t, 2, loggerInfoCounter)
	assert.Equal(t, 0, loggerErrorCounter)
}

func TestGetSecretFromStringFailure(t *testing.T) {
	setup()

	secretsManagerGetSecretValueReturn.SecretString = aws2.String("")

	secret, err := secretsManagerService.GetSecret("secretName")

	assert.Nil(t, secret)
	assert.Equal(t, "unexpected end of JSON input", err.(custom_error.BaseErrorAdapter).Error())
	assert.Equal(t, "error while unmarshalling secret string", err.(custom_error.BaseErrorAdapter).InternalError())
	assert.Equal(t, "Error while trying to get secret from secrets manager", err.(custom_error.BaseErrorAdapter).Description())
	assert.Equal(t, 1, secretsManagerGetSecretValueCounter)
	assert.Equal(t, 1, loggerInfoCounter)
	assert.Equal(t, 1, loggerErrorCounter)
}

func TestGetSecretFromBinarySuccess(t *testing.T) {
	setup()

	secretsManagerGetSecretValueReturn.SecretString = nil

	secret, err := secretsManagerService.GetSecret("secretName")

	assert.Nil(t, err)
	assert.Equal(t, &secrets, secret)
	assert.Equal(t, 1, secretsManagerGetSecretValueCounter)
	assert.Equal(t, 2, loggerInfoCounter)
	assert.Equal(t, 0, loggerErrorCounter)
}

func TestGetSecretFromBinaryUnmarshalFailure(t *testing.T) {
	setup()

	secretsManagerGetSecretValueReturn.SecretString = nil
	secretsManagerGetSecretValueReturn.SecretBinary = []byte{}

	secret, err := secretsManagerService.GetSecret("secretName")

	assert.Nil(t, secret)
	assert.Equal(t, "unexpected end of JSON input", err.(custom_error.BaseErrorAdapter).Error())
	assert.Equal(t, "error while unmarshalling secret binary", err.(custom_error.BaseErrorAdapter).InternalError())
	assert.Equal(t, "Error while trying to get secret from secrets manager", err.(custom_error.BaseErrorAdapter).Description())
	assert.Equal(t, 1, secretsManagerGetSecretValueCounter)
	assert.Equal(t, 1, loggerInfoCounter)
	assert.Equal(t, 1, loggerErrorCounter)
}

func TestGetSecretFromBinaryDecodeFailure(t *testing.T) {
	setup()

	secretsManagerGetSecretValueReturn.SecretString = nil
	secretsManagerGetSecretValueReturn.SecretBinary = []byte{1, 2, 3, 4, 5}

	secret, err := secretsManagerService.GetSecret("secretName")

	assert.Nil(t, secret)
	assert.Equal(t, "illegal base64 data at input byte 0", err.(custom_error.BaseErrorAdapter).Error())
	assert.Equal(t, "error while decoding secret binary", err.(custom_error.BaseErrorAdapter).InternalError())
	assert.Equal(t, "Error while trying to get secret from secrets manager", err.(custom_error.BaseErrorAdapter).Description())
	assert.Equal(t, 1, secretsManagerGetSecretValueCounter)
	assert.Equal(t, 1, loggerInfoCounter)
	assert.Equal(t, 1, loggerErrorCounter)
}

func TestGetSecretSecretsManagerFailure(t *testing.T) {
	setup()

	secretsManagerGetSecretValueError = errors.New("error test")

	secret, err := secretsManagerService.GetSecret("secretName")

	assert.Nil(t, secret)
	assert.Equal(t, "error test", err.(custom_error.BaseErrorAdapter).Error())
	assert.Equal(t, "error while getting secret", err.(custom_error.BaseErrorAdapter).InternalError())
	assert.Equal(t, "Error while trying to get secret from secrets manager", err.(custom_error.BaseErrorAdapter).Description())
	assert.Equal(t, 1, secretsManagerGetSecretValueCounter)
	assert.Equal(t, 1, loggerInfoCounter)
	assert.Equal(t, 1, loggerErrorCounter)
}
