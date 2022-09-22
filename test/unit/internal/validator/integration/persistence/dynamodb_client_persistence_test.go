package persistence

import (
	"github.com/brienze1/crypto-robot-validator/internal/validator/application/config"
	"github.com/brienze1/crypto-robot-validator/internal/validator/domain/adapters"
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
)

func setup() {
	config.LoadTestEnv()

	clientPersistence = persistence.DynamoDBClientPersistence(logger, dynamoDBClient)

	dynamoDBClient.Reset()

	clientPersisted = &dto.Client{Id: uuid.NewString()}

	dynamoDBClient.AddItem(clientPersisted.Id, clientPersisted)
}

func TestGetClientsSuccess(t *testing.T) {
	setup()

	client, err := clientPersistence.GetClient(clientPersisted.Id)

	assert.Nilf(t, err, "Should be nil")
	assert.NotNilf(t, client, "Should not be nil")
	assert.Equal(t, clientPersisted.Id, client.Id)
}
