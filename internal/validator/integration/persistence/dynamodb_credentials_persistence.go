package persistence

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/brienze1/crypto-robot-validator/internal/validator/application/properties"
	"github.com/brienze1/crypto-robot-validator/internal/validator/domain/adapters"
	adapters2 "github.com/brienze1/crypto-robot-validator/internal/validator/integration/adapters"
	"github.com/brienze1/crypto-robot-validator/internal/validator/integration/dto"
	"github.com/brienze1/crypto-robot-validator/internal/validator/integration/exceptions"
	"github.com/brienze1/crypto-robot-validator/pkg/custom_error"
)

type dynamoDBCredentialsPersistence struct {
	logger   adapters.LoggerAdapter
	dynamoDB adapters2.DynamoDBAdapter
}

// DynamoDBCredentialsPersistence class constructor
func DynamoDBCredentialsPersistence(logger adapters.LoggerAdapter, dynamoDB adapters2.DynamoDBAdapter) *dynamoDBCredentialsPersistence {
	return &dynamoDBCredentialsPersistence{
		logger:   logger,
		dynamoDB: dynamoDB,
	}
}

// GetCredentials will find dto.Credentials on credentials DynamoDB repository using clientId as key.
func (d *dynamoDBCredentialsPersistence) GetCredentials(clientId string) (*dto.Credentials, custom_error.BaseErrorAdapter) {
	d.logger.Info("GetCredentials started", clientId)

	response, err := d.dynamoDB.GetItem(context.TODO(), &dynamodb.GetItemInput{
		Key: map[string]types.AttributeValue{
			"client_id": &types.AttributeValueMemberS{Value: clientId},
		},
		TableName: properties.Properties().Aws.DynamoDB.CredentialsTableName,
	})
	if err != nil {
		return nil, d.abort(err, "Error while trying to get credentials.")
	}

	if response.Item == nil {
		return nil, d.abort(err, "Credentials not found.")
	}

	var credentials *dto.Credentials
	err = attributevalue.UnmarshalMap(response.Item, &credentials)
	if err != nil {
		return nil, d.abort(err, "Error while trying to unmarshal get credentials response.")
	}

	d.logger.Info("GetCredentials finished", clientId)
	return credentials, nil
}

func (d *dynamoDBCredentialsPersistence) abort(err error, message string) custom_error.BaseErrorAdapter {
	dynamoDBClientPersistenceError := exceptions.DynamoDBCredentialsPersistenceError(err, message)
	d.logger.Error(dynamoDBClientPersistenceError, "Get credentials failed: "+message)
	return dynamoDBClientPersistenceError
}
