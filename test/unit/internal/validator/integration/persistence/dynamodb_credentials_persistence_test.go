package persistence

import (
	"errors"
	"github.com/brienze1/crypto-robot-validator/internal/validator/integration/adapters"
	"github.com/brienze1/crypto-robot-validator/internal/validator/integration/dto"
	"github.com/brienze1/crypto-robot-validator/internal/validator/integration/persistence"
	"github.com/brienze1/crypto-robot-validator/test/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	credentialsPersistence adapters.CredentialsPersistenceAdapter
	loggerCredentials      = mocks.Logger()
	dynamoDBCredentials    = mocks.DynamoDBClient()
)

var (
	credentials *dto.Credentials
)

func setupCredentialsPersistence() {
	loggerCredentials.Reset()
	dynamoDBCredentials.Reset()

	credentials = &dto.Credentials{
		ClientId:  uuid.NewString(),
		ApiKey:    uuid.NewString(),
		ApiSecret: uuid.NewString(),
	}

	dynamoDBCredentials.AddItem(credentials.ClientId, credentials)

	credentialsPersistence = persistence.DynamoDBCredentialsPersistence(loggerCredentials, dynamoDBCredentials)
}

func TestGetCredentialsSuccess(t *testing.T) {
	setupCredentialsPersistence()

	credentialsPersisted, err := credentialsPersistence.GetCredentials(credentials.ClientId)

	assert.Nil(t, err)
	assert.NotNil(t, credentialsPersisted)
	assert.Equal(t, credentials.ClientId, credentialsPersisted.ClientId)
	assert.Equal(t, credentials.ApiKey, credentialsPersisted.ApiKey)
	assert.Equal(t, credentials.ApiSecret, credentialsPersisted.ApiSecret)
	assert.Equal(t, 1, dynamoDBCredentials.GetItemCounter)
	assert.Equal(t, 2, loggerCredentials.InfoCallCounter)
	assert.Equal(t, 0, loggerCredentials.ErrorCallCounter)
}

func TestGetCredentialsDynamoDBErrorFailure(t *testing.T) {
	setupCredentialsPersistence()

	dynamoDBCredentials.GetItemError = errors.New("get item error")

	credentialsPersisted, err := credentialsPersistence.GetCredentials(credentials.ClientId)

	assert.Nil(t, credentialsPersisted)
	assert.NotNil(t, err)
	assert.Equal(t, "get item error", err.Error())
	assert.Equal(t, "Error while trying to get credentials.", err.InternalError())
	assert.Equal(t, "Error while using DynamoDB Credentials table", err.Description())
	assert.Equal(t, 1, dynamoDBCredentials.GetItemCounter)
	assert.Equal(t, 1, loggerCredentials.InfoCallCounter)
	assert.Equal(t, 1, loggerCredentials.ErrorCallCounter)
}

func TestGetCredentialsNotFoundFailure(t *testing.T) {
	setupCredentialsPersistence()

	dynamoDBCredentials.Reset()

	credentialsPersisted, err := credentialsPersistence.GetCredentials(credentials.ClientId)

	assert.Nil(t, credentialsPersisted)
	assert.NotNil(t, err)
	assert.Equal(t, "Credentials not found.", err.Error())
	assert.Equal(t, "Credentials not found.", err.InternalError())
	assert.Equal(t, "Error while using DynamoDB Credentials table", err.Description())
	assert.Equal(t, 1, dynamoDBCredentials.GetItemCounter)
	assert.Equal(t, 1, loggerCredentials.InfoCallCounter)
	assert.Equal(t, 1, loggerCredentials.ErrorCallCounter)
}

func TestGetCredentialsUnmarshalFailure(t *testing.T) {
	setupCredentialsPersistence()

	fakeClient := map[string]interface{}{
		"client_id": false,
	}

	dynamoDBCredentials.AddItem(credentials.ClientId, fakeClient)

	credentialsPersisted, err := credentialsPersistence.GetCredentials(credentials.ClientId)

	assert.Nil(t, credentialsPersisted)
	assert.NotNil(t, err)
	assert.Equal(t, "unmarshal failed, cannot unmarshal bool into Go value type string", err.Error())
	assert.Equal(t, "Error while trying to unmarshal get credentials response.", err.InternalError())
	assert.Equal(t, "Error while using DynamoDB Credentials table", err.Description())
	assert.Equal(t, 1, dynamoDBCredentials.GetItemCounter)
	assert.Equal(t, 1, loggerCredentials.InfoCallCounter)
	assert.Equal(t, 1, loggerCredentials.ErrorCallCounter)
}
