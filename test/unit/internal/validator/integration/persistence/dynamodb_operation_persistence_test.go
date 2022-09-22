package persistence

import (
	"context"
	"errors"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/brienze1/crypto-robot-validator/internal/validator/application/config"
	"github.com/brienze1/crypto-robot-validator/internal/validator/application/properties"
	"github.com/brienze1/crypto-robot-validator/internal/validator/domain/adapters"
	"github.com/brienze1/crypto-robot-validator/internal/validator/domain/model"
	"github.com/brienze1/crypto-robot-validator/internal/validator/integration/dto"
	"github.com/brienze1/crypto-robot-validator/internal/validator/integration/persistence"
	"github.com/brienze1/crypto-robot-validator/test/mocks"
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	operationPersistence adapters.OperationPersistenceAdapter
	loggerMock           = mocks.Logger()
	dynamoDBClientMock   = mocks.DynamoDBClient()
)

var (
	operation *model.Operation
)

func setup() {
	config.LoadTestEnv()

	operationPersistence = persistence.DynamoDBOperationPersistence(loggerMock, dynamoDBClientMock)

	loggerMock.Reset()
	dynamoDBClientMock.Reset()

	operation = model.NewOperation(50.00)
}

func TestSaveSuccess(t *testing.T) {
	setup()

	err := operationPersistence.Save(operation)

	assert.Nilf(t, err, "Should be nil")
	assert.Equal(t, 1, dynamoDBClientMock.PutItemCounter)
	assert.Equal(t, 2, loggerMock.InfoCallCounter)
	assert.Equal(t, 0, loggerMock.ErrorCallCounter)

	response, _ := dynamoDBClientMock.GetItem(context.TODO(), &dynamodb.GetItemInput{
		Key: map[string]types.AttributeValue{
			"operation_id": &types.AttributeValueMemberS{Value: operation.Id},
		},
		TableName: properties.Properties().Aws.DynamoDB.OperationTableName,
	})

	var operationSaved *dto.Operation
	_ = attributevalue.UnmarshalMap(response.Item, &operationSaved)

	assert.Equal(t, operation.Id, operationSaved.Id)
}

func TestSavePutItemFailure(t *testing.T) {
	setup()

	dynamoDBClientMock.PutItemError = errors.New("put item error")

	err := operationPersistence.Save(operation)

	assert.Equal(t, "put item error", err.Error())
	assert.Equal(t, "PutItem error", err.InternalError())
	assert.Equal(t, "Error while using DynamoDB Operation table", err.Description())
	assert.Equal(t, 1, dynamoDBClientMock.PutItemCounter)
	assert.Equal(t, 1, loggerMock.InfoCallCounter)
	assert.Equal(t, 1, loggerMock.ErrorCallCounter)
}
