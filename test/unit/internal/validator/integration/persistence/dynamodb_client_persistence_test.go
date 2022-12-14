package persistence

import (
	"errors"
	"github.com/brienze1/crypto-robot-validator/internal/validator/application/config"
	"github.com/brienze1/crypto-robot-validator/internal/validator/application/properties"
	"github.com/brienze1/crypto-robot-validator/internal/validator/domain/adapters"
	"github.com/brienze1/crypto-robot-validator/internal/validator/domain/model"
	"github.com/brienze1/crypto-robot-validator/internal/validator/integration/dto"
	"github.com/brienze1/crypto-robot-validator/internal/validator/integration/persistence"
	"github.com/brienze1/crypto-robot-validator/test/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	clientPersistence adapters.ClientPersistenceAdapter
	logger            = mocks.Logger()
	dynamoDBClient    = mocks.DynamoDBClient()
)

var (
	clientPersisted *dto.Client
	clientUnlocked  *model.Client
	clientLocked    *model.Client
)

func clientPersistenceSetup() {
	config.LoadTestEnv()

	clientPersistence = persistence.DynamoDBClientPersistence(logger, dynamoDBClient)

	logger.Reset()
	dynamoDBClient.Reset()

	clientPersisted = &dto.Client{Id: uuid.NewString(), Locked: false}
	clientUnlocked = &model.Client{Id: uuid.NewString(), Locked: false}
	clientLocked = &model.Client{Id: uuid.NewString(), Locked: true}

	dynamoDBClient.AddItem(clientPersisted.Id, clientPersisted, properties.Properties().Aws.DynamoDB.ClientTableName)
}

func TestGetClientsSuccess(t *testing.T) {
	clientPersistenceSetup()

	client, err := clientPersistence.GetClient(clientPersisted.Id)

	assert.Nilf(t, err, "Should be nil")
	assert.NotNilf(t, client, "Should not be nil")
	assert.Equal(t, clientPersisted.Id, client.Id)
	assert.Equal(t, 1, dynamoDBClient.GetItemCounter)
	assert.Equal(t, 2, logger.InfoCallCounter)
	assert.Equal(t, 0, logger.ErrorCallCounter)
}

func TestGetClientsClientNotFoundFailure(t *testing.T) {
	clientPersistenceSetup()

	client, err := clientPersistence.GetClient(uuid.NewString())

	assert.Equal(t, "Client not found.", err.Error())
	assert.Equal(t, "Client not found.", err.InternalError())
	assert.Equal(t, "Error while using DynamoDB Client table", err.Description())
	assert.Nilf(t, client, "Should be nil")
	assert.Equal(t, 1, dynamoDBClient.GetItemCounter)
	assert.Equal(t, 1, logger.InfoCallCounter)
	assert.Equal(t, 1, logger.ErrorCallCounter)
}

func TestGetClientsDynamoDBClientFailure(t *testing.T) {
	clientPersistenceSetup()

	dynamoDBClient.GetItemError = errors.New("dynamodb client error")

	client, err := clientPersistence.GetClient(uuid.NewString())

	assert.Equal(t, "dynamodb client error", err.Error())
	assert.Equal(t, "Error while trying to get client.", err.InternalError())
	assert.Equal(t, "Error while using DynamoDB Client table", err.Description())
	assert.Nilf(t, client, "Should be nil")
	assert.Equal(t, 1, dynamoDBClient.GetItemCounter)
	assert.Equal(t, 1, logger.InfoCallCounter)
	assert.Equal(t, 1, logger.ErrorCallCounter)
}

func TestGetClientsUnmarshalFailure(t *testing.T) {
	clientPersistenceSetup()

	fakeClient := map[string]interface{}{
		"locked": "test",
	}
	clientId := uuid.NewString()

	dynamoDBClient.AddItem(clientId, fakeClient, properties.Properties().Aws.DynamoDB.ClientTableName)

	client, err := clientPersistence.GetClient(clientId)

	assert.Equal(t, "unmarshal failed, cannot unmarshal string into Go value type bool", err.Error())
	assert.Equal(t, "Error while trying to unmarshal get client response.", err.InternalError())
	assert.Equal(t, "Error while using DynamoDB Client table", err.Description())
	assert.Nilf(t, client, "Should be nil")
	assert.Equal(t, 1, dynamoDBClient.GetItemCounter)
	assert.Equal(t, 1, logger.InfoCallCounter)
	assert.Equal(t, 1, logger.ErrorCallCounter)
}

func TestGetClientsClientLockedFailure(t *testing.T) {
	clientPersistenceSetup()

	dynamoDBClient.AddItem(clientLocked.Id, clientLocked, properties.Properties().Aws.DynamoDB.ClientTableName)

	client, err := clientPersistence.GetClient(clientLocked.Id)

	assert.Equal(t, "Client is locked.", err.Error())
	assert.Equal(t, "Client is locked.", err.InternalError())
	assert.Equal(t, "Error while using DynamoDB Client table", err.Description())
	assert.Nilf(t, client, "Should be nil")
	assert.Equal(t, 1, dynamoDBClient.GetItemCounter)
	assert.Equal(t, 1, logger.InfoCallCounter)
	assert.Equal(t, 1, logger.ErrorCallCounter)
}

func TestLockSuccess(t *testing.T) {
	clientPersistenceSetup()

	assert.Equal(t, false, clientUnlocked.Locked)

	err := clientPersistence.Lock(clientUnlocked)

	assert.Nilf(t, err, "Should be nil")
	assert.Equal(t, true, clientUnlocked.Locked)
	assert.Equal(t, 1, dynamoDBClient.PutItemCounter)
	assert.Equal(t, 2, logger.InfoCallCounter)
	assert.Equal(t, 0, logger.ErrorCallCounter)

	_, err = clientPersistence.GetClient(clientUnlocked.Id)

	assert.Equal(t, "Client is locked.", err.Error())
	assert.Equal(t, "Client is locked.", err.InternalError())
	assert.Equal(t, "Error while using DynamoDB Client table", err.Description())
}

func TestLockPutItemFailure(t *testing.T) {
	clientPersistenceSetup()

	dynamoDBClient.PutItemError = errors.New("lock error")

	err := clientPersistence.Lock(clientUnlocked)

	assert.Equal(t, "lock error", err.Error())
	assert.Equal(t, "PutItem error", err.InternalError())
	assert.Equal(t, "Error while using DynamoDB Client table", err.Description())
	assert.Equal(t, true, clientUnlocked.Locked)
	assert.Equal(t, 1, dynamoDBClient.PutItemCounter)
	assert.Equal(t, 1, logger.InfoCallCounter)
	assert.Equal(t, 2, logger.ErrorCallCounter)
}

func TestUnlockSuccess(t *testing.T) {
	clientPersistenceSetup()

	assert.Equal(t, true, clientLocked.Locked)

	err := clientPersistence.Unlock(clientLocked)

	assert.Nilf(t, err, "Should be nil")
	assert.Equal(t, false, clientLocked.Locked)
	assert.Equal(t, 1, dynamoDBClient.PutItemCounter)
	assert.Equal(t, 2, logger.InfoCallCounter)
	assert.Equal(t, 0, logger.ErrorCallCounter)

	clientUpdated, err := clientPersistence.GetClient(clientLocked.Id)

	assert.Nilf(t, err, "Should be nil")
	assert.NotNilf(t, clientUpdated, "Should not be nil")
	assert.Equal(t, clientLocked.Id, clientUpdated.Id)
	assert.Equal(t, 1, dynamoDBClient.GetItemCounter)
	assert.Equal(t, 4, logger.InfoCallCounter)
	assert.Equal(t, 0, logger.ErrorCallCounter)
}

func TestUnlockPutItemFailure(t *testing.T) {
	clientPersistenceSetup()

	dynamoDBClient.PutItemError = errors.New("unlock error")

	err := clientPersistence.Unlock(clientLocked)

	assert.Equal(t, "unlock error", err.Error())
	assert.Equal(t, "PutItem error", err.InternalError())
	assert.Equal(t, "Error while using DynamoDB Client table", err.Description())
	assert.Equal(t, false, clientLocked.Locked)
	assert.Equal(t, 1, dynamoDBClient.PutItemCounter)
	assert.Equal(t, 1, logger.InfoCallCounter)
	assert.Equal(t, 2, logger.ErrorCallCounter)
}

func TestLockAndUnlockSuccess(t *testing.T) {
	clientPersistenceSetup()

	assert.Equal(t, false, clientUnlocked.Locked)

	err := clientPersistence.Lock(clientUnlocked)

	assert.Nilf(t, err, "Should be nil")
	assert.Equal(t, true, clientUnlocked.Locked)
	assert.Equal(t, 1, dynamoDBClient.PutItemCounter)
	assert.Equal(t, 2, logger.InfoCallCounter)
	assert.Equal(t, 0, logger.ErrorCallCounter)

	_, err = clientPersistence.GetClient(clientUnlocked.Id)

	assert.Equal(t, "Client is locked.", err.Error())
	assert.Equal(t, "Client is locked.", err.InternalError())
	assert.Equal(t, "Error while using DynamoDB Client table", err.Description())

	err = clientPersistence.Unlock(clientUnlocked)

	assert.Nilf(t, err, "Should be nil")
	assert.Equal(t, false, clientUnlocked.Locked)
	assert.Equal(t, 2, dynamoDBClient.PutItemCounter)
	assert.Equal(t, 5, logger.InfoCallCounter)
	assert.Equal(t, 1, logger.ErrorCallCounter)

	clientUpdated, err := clientPersistence.GetClient(clientUnlocked.Id)
	assert.Nilf(t, err, "Should be nil")
	assert.NotNilf(t, clientUpdated, "Should not be nil")
	assert.Equal(t, clientUnlocked.Id, clientUpdated.Id)
	assert.Equal(t, 2, dynamoDBClient.GetItemCounter)
	assert.Equal(t, 7, logger.InfoCallCounter)
	assert.Equal(t, 1, logger.ErrorCallCounter)
}
