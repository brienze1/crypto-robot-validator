package persistence

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/brienze1/crypto-robot-validator/internal/validator/application/properties"
	"github.com/brienze1/crypto-robot-validator/internal/validator/domain/adapters"
	"github.com/brienze1/crypto-robot-validator/internal/validator/domain/model"
	adapters2 "github.com/brienze1/crypto-robot-validator/internal/validator/integration/adapters"
	"github.com/brienze1/crypto-robot-validator/internal/validator/integration/exceptions"
	"github.com/brienze1/crypto-robot-validator/pkg/custom_error"
)

type dynamoDBOperationPersistence struct {
	logger   adapters.LoggerAdapter
	dynamoDB adapters2.DynamoDBAdapter
}

// DynamoDBOperationPersistence class constructor
func DynamoDBOperationPersistence(logger adapters.LoggerAdapter, dynamoDB adapters2.DynamoDBAdapter) *dynamoDBOperationPersistence {
	return &dynamoDBOperationPersistence{
		logger:   logger,
		dynamoDB: dynamoDB,
	}
}

func (d *dynamoDBOperationPersistence) Save(operation *model.Operation) custom_error.BaseErrorAdapter {
	clientInput, err := attributevalue.MarshalMap(operation)
	if err != nil {
		return d.abort(err, "Error while trying to marshal client.")
	}

	_, err = d.dynamoDB.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: properties.Properties().Aws.DynamoDB.ClientTableName,
		Item:      clientInput,
	})
	if err != nil {
		return d.abort(err, "Error while trying to update client.")
	}

	return nil
}

func (d *dynamoDBOperationPersistence) abort(err error, message string) custom_error.BaseErrorAdapter {
	dynamoDBClientPersistenceError := exceptions.DynamoDBClientPersistenceError(err, message)
	d.logger.Error(dynamoDBClientPersistenceError, "Get clients failed: "+message)
	return dynamoDBClientPersistenceError
}
