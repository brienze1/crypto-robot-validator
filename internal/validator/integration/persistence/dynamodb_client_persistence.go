package persistence

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/brienze1/crypto-robot-validator/internal/validator/application/properties"
	"github.com/brienze1/crypto-robot-validator/internal/validator/domain/adapters"
	"github.com/brienze1/crypto-robot-validator/internal/validator/domain/model"
	adapters2 "github.com/brienze1/crypto-robot-validator/internal/validator/integration/adapters"
	"github.com/brienze1/crypto-robot-validator/internal/validator/integration/dto"
	"github.com/brienze1/crypto-robot-validator/internal/validator/integration/exceptions"
	"github.com/brienze1/crypto-robot-validator/pkg/custom_error"
)

type dynamoDBClientPersistence struct {
	logger   adapters.LoggerAdapter
	dynamoDB adapters2.DynamoDBAdapter
}

// DynamoDBClientPersistence class constructor
func DynamoDBClientPersistence(logger adapters.LoggerAdapter, dynamoDB adapters2.DynamoDBAdapter) *dynamoDBClientPersistence {
	return &dynamoDBClientPersistence{
		logger:   logger,
		dynamoDB: dynamoDB,
	}
}

// GetClient will find model.Client on client DynamoDB repository using clientId as key.
func (d *dynamoDBClientPersistence) GetClient(clientId string) (*model.Client, custom_error.BaseErrorAdapter) {
	d.logger.Info("GetClient started", clientId)

	response, err := d.dynamoDB.GetItem(context.TODO(), &dynamodb.GetItemInput{
		Key: map[string]types.AttributeValue{
			"client_id": &types.AttributeValueMemberS{Value: clientId},
		},
		TableName: properties.Properties().Aws.DynamoDB.ClientTableName,
	})
	if err != nil {
		return nil, d.abort(err, "Error while trying to get client.")
	}

	if response.Item == nil {
		return nil, d.abort(err, "Client not found.")
	}

	var client *dto.Client
	err = attributevalue.UnmarshalMap(response.Item, &client)
	if err != nil {
		return nil, d.abort(err, "Error while trying to unmarshal get client response.")
	}

	if client.Locked {
		return nil, d.abort(err, "Client is locked.")
	}

	d.logger.Info("GetClient finished", clientId, client)
	return client.ToModel(), nil
}

// Lock will update model.Client setting flag locked as true on client DynamoDB repository. Returns error if client is
// already locked.
func (d *dynamoDBClientPersistence) Lock(client *model.Client) custom_error.BaseErrorAdapter {
	d.logger.Info("Lock started", client)

	client.Lock()

	clientDto := dto.ClientDto(client)

	err := d.update(clientDto)
	if err != nil {
		return d.abort(err, "Error while trying to lock client.")
	}

	d.logger.Info("Lock finished", client)
	return nil
}

// Unlock will update model.Client setting flag locked as false on client DynamoDB repository.
func (d *dynamoDBClientPersistence) Unlock(client *model.Client) custom_error.BaseErrorAdapter {
	d.logger.Info("Unlock started", client)

	client.Unlock()

	clientDto := dto.ClientDto(client)

	err := d.update(clientDto)
	if err != nil {
		return d.abort(err, "Error while trying to unlock client.")
	}

	d.logger.Info("Unlock finished", client)
	return nil
}

func (d *dynamoDBClientPersistence) update(client *dto.Client) custom_error.BaseErrorAdapter {
	clientInput, err := attributevalue.MarshalMap(client)
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

func (d *dynamoDBClientPersistence) abort(err error, message string) custom_error.BaseErrorAdapter {
	dynamoDBClientPersistenceError := exceptions.DynamoDBClientPersistenceError(err, message)
	d.logger.Error(dynamoDBClientPersistenceError, "Get clients failed: "+message)
	return dynamoDBClientPersistenceError
}
