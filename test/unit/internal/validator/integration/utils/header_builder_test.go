package utils

import (
	"errors"
	"github.com/brienze1/crypto-robot-validator/internal/validator/application/config"
	"github.com/brienze1/crypto-robot-validator/internal/validator/application/properties"
	"github.com/brienze1/crypto-robot-validator/internal/validator/integration/adapters"
	"github.com/brienze1/crypto-robot-validator/internal/validator/integration/dto"
	"github.com/brienze1/crypto-robot-validator/internal/validator/integration/utils"
	"github.com/brienze1/crypto-robot-validator/pkg/time_utils"
	"github.com/brienze1/crypto-robot-validator/test/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	headerBuilder            adapters.HeaderBuilderAdapter
	loggerHB                 = mocks.Logger()
	credentialsPersistenceHB = mocks.DynamoDBCredentialPersistence()
	secretsManagerServiceHB  = mocks.SecretsManagerService()
	encryptionServiceHB      = mocks.EncryptionService()
	tokenBuilderHB           = mocks.TokenBuilder()
)

var (
	credentialsPersisted *dto.Credentials
	encryptionSecrets    *dto.EncryptionSecrets
	expectedTokenHB      = uuid.NewString()
	nonceHB              = time_utils.Epoch()
	clientIdHB           = uuid.NewString()
	endpointHB           = "v1/balance"
	payloadHB            = `{}`
)

func setupHB() {
	config.LoadTestEnv()

	loggerHB.Reset()
	credentialsPersistenceHB.Reset()
	secretsManagerServiceHB.Reset()
	encryptionServiceHB.Reset()
	tokenBuilderHB.Reset()

	headerBuilder = utils.HeaderBuilder(loggerHB, credentialsPersistenceHB, secretsManagerServiceHB, encryptionServiceHB, tokenBuilderHB)

	credentialsPersisted = &dto.Credentials{
		ClientId:  clientIdHB,
		ApiKey:    uuid.NewString(),
		ApiSecret: "7571a562734fe71e042c9839ea0beb80f033779aa7379b5b93f0fcc10b0a74db99712329ba92",
	}

	credentialsPersistenceHB.AddCredential(credentialsPersisted)

	encryptionSecrets = &dto.EncryptionSecrets{
		EncryptionKey: "9y$B?E(H+MbQeThWmZq4t7w!z%C*F)J@",
	}

	secretsManagerServiceHB.SetSecret(properties.Properties().Aws.SecretsManager.EncryptionSecretName, encryptionSecrets)

	tokenBuilderHB.ExpectedToken = expectedTokenHB
}

func TestBuildBiscointHeaderSuccess(t *testing.T) {
	setupHB()

	header, err := headerBuilder.BiscointHeader(clientIdHB, endpointHB, payloadHB)

	assert.Nil(t, err)
	assert.NotNil(t, header)
	assert.Equal(t, "application/json", header.Get("Content-Type"))
	assert.Equal(t, nonceHB, header.Get("BSCNT-NONCE"))
	assert.Equal(t, credentialsPersisted.ApiKey, header.Get("BSCNT-APIKEY"))
	assert.Equal(t, expectedTokenHB, header.Get("BSCNT-SIGN"))
	assert.Equal(t, 2, loggerHB.InfoCallCounter)
	assert.Equal(t, 0, loggerHB.ErrorCallCounter)
	assert.Equal(t, 1, credentialsPersistenceHB.GetCredentialsCounter)
	assert.Equal(t, 1, secretsManagerServiceHB.GetSecretCounter)
	assert.Equal(t, 1, encryptionServiceHB.AESDecryptCounter)
	assert.Equal(t, 1, tokenBuilderHB.BuildCounter)
}

func TestBuildBiscointHeaderCredentialsFailure(t *testing.T) {
	setupHB()

	credentialsPersistenceHB.GetCredentialsError = errors.New("get credentials error")

	header, err := headerBuilder.BiscointHeader(clientIdHB, endpointHB, payloadHB)

	assert.NotNil(t, err)
	assert.Nil(t, header)
	assert.Equal(t, "get credentials error", err.Error())
	assert.Equal(t, "GetCredentials error", err.InternalError())
	assert.Equal(t, "Error while using DynamoDB Credentials table", err.Description())
	assert.Equal(t, 1, loggerHB.InfoCallCounter)
	assert.Equal(t, 1, loggerHB.ErrorCallCounter)
	assert.Equal(t, 1, credentialsPersistenceHB.GetCredentialsCounter)
	assert.Equal(t, 0, secretsManagerServiceHB.GetSecretCounter)
	assert.Equal(t, 0, encryptionServiceHB.AESDecryptCounter)
	assert.Equal(t, 0, tokenBuilderHB.BuildCounter)
}

func TestBuildBiscointHeaderGetSecretFailure(t *testing.T) {
	setupHB()

	secretsManagerServiceHB.GetSecretError = errors.New("get secret error")

	header, err := headerBuilder.BiscointHeader(clientIdHB, endpointHB, payloadHB)

	assert.NotNil(t, err)
	assert.Nil(t, header)
	assert.Equal(t, "get secret error", err.Error())
	assert.Equal(t, "secrets manager error", err.InternalError())
	assert.Equal(t, "Error while trying to get secret from secrets manager", err.Description())
	assert.Equal(t, 1, loggerHB.InfoCallCounter)
	assert.Equal(t, 1, loggerHB.ErrorCallCounter)
	assert.Equal(t, 1, credentialsPersistenceHB.GetCredentialsCounter)
	assert.Equal(t, 1, secretsManagerServiceHB.GetSecretCounter)
	assert.Equal(t, 0, encryptionServiceHB.AESDecryptCounter)
	assert.Equal(t, 0, tokenBuilderHB.BuildCounter)
}

func TestBuildBiscointHeaderAESDecryptFailure(t *testing.T) {
	setupHB()

	encryptionServiceHB.AESDecryptError = errors.New("AES decrypt error")

	header, err := headerBuilder.BiscointHeader(clientIdHB, endpointHB, payloadHB)

	assert.NotNil(t, err)
	assert.Nil(t, header)
	assert.Equal(t, "AES decrypt error", err.Error())
	assert.Equal(t, "AES decrypt error", err.InternalError())
	assert.Equal(t, "Error while performing encryption", err.Description())
	assert.Equal(t, 1, loggerHB.InfoCallCounter)
	assert.Equal(t, 1, loggerHB.ErrorCallCounter)
	assert.Equal(t, 1, credentialsPersistenceHB.GetCredentialsCounter)
	assert.Equal(t, 1, secretsManagerServiceHB.GetSecretCounter)
	assert.Equal(t, 1, encryptionServiceHB.AESDecryptCounter)
	assert.Equal(t, 0, tokenBuilderHB.BuildCounter)
}

func TestBuildBiscointHeaderTokenBuildFailure(t *testing.T) {
	setupHB()

	tokenBuilderHB.BuildError = errors.New("token build error")

	header, err := headerBuilder.BiscointHeader(clientIdHB, endpointHB, payloadHB)

	assert.NotNil(t, err)
	assert.Nil(t, header)
	assert.Equal(t, "token build error", err.Error())
	assert.Equal(t, "token build error", err.InternalError())
	assert.Equal(t, "Error while building token", err.Description())
	assert.Equal(t, 1, loggerHB.InfoCallCounter)
	assert.Equal(t, 1, loggerHB.ErrorCallCounter)
	assert.Equal(t, 1, credentialsPersistenceHB.GetCredentialsCounter)
	assert.Equal(t, 1, secretsManagerServiceHB.GetSecretCounter)
	assert.Equal(t, 1, encryptionServiceHB.AESDecryptCounter)
	assert.Equal(t, 1, tokenBuilderHB.BuildCounter)
}
