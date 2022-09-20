package persistence

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/brienze1/crypto-robot-operation-hub/internal/operation-hub/application/properties"
	"github.com/brienze1/crypto-robot-operation-hub/internal/operation-hub/domain/adapters"
	"github.com/brienze1/crypto-robot-operation-hub/internal/operation-hub/domain/model"
	adapters2 "github.com/brienze1/crypto-robot-operation-hub/internal/operation-hub/integration/adapters"
	"github.com/brienze1/crypto-robot-operation-hub/internal/operation-hub/integration/dto"
	"github.com/brienze1/crypto-robot-operation-hub/internal/operation-hub/integration/exceptions"
)

type dynamoDBClientPersistence struct {
	logger     adapters.LoggerAdapter
	dynamoDB   adapters2.DynamoDBAdapter
	timeSource adapters.TimeAdapter
}

func DynamoDBClientPersistence(logger adapters.LoggerAdapter, dynamoDB adapters2.DynamoDBAdapter, timeSource adapters.TimeAdapter) *dynamoDBClientPersistence {
	return &dynamoDBClientPersistence{
		logger:     logger,
		dynamoDB:   dynamoDB,
		timeSource: timeSource,
	}
}

func (d *dynamoDBClientPersistence) GetClients(config *model.ClientSearchConfig) (*[]model.Client, error) {
	d.logger.Info("Get clients start", config)

	expr, err := expression.NewBuilder().WithFilter(
		expression.And(
			expression.Name("active").Equal(expression.Value(config.Active)),
			expression.Name("locked").Equal(expression.Value(config.Locked)),
			expression.Name("locked_until").LessThanEqual(expression.Value(d.timeSource.Now())),
			expression.Name("cash_amount").GreaterThanEqual(expression.Value(config.MinimumCash)),
			expression.Name("crypto_amount").GreaterThanEqual(expression.Value(config.MinimumCrypto)),
			expression.Name("sell_on").LessThanEqual(expression.Value(config.SellWeight.Value())),
			expression.Name("buy_on").LessThanEqual(expression.Value(config.BuyWeight.Value())),
			expression.Name("symbols").Contains(config.Symbol.Name()),
		),
	).Build()
	if err != nil {
		return nil, d.abort(err, "DynamoDb expression builder error")
	}

	result, err := d.dynamoDB.Scan(context.TODO(), &dynamodb.ScanInput{
		TableName:                 aws.String(properties.Properties().Aws.DynamoDB.ClientTableName),
		ProjectionExpression:      aws.String("client_id"),
		FilterExpression:          expr.Filter(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
	})
	if err != nil {
		return nil, d.abort(err, "DynamoDb scan error")
	}

	clientsDto := &[]dto.Client{}
	clients := &[]model.Client{}

	err = attributevalue.UnmarshalListOfMaps(result.Items, clientsDto)
	if err != nil {
		return nil, d.abort(err, "Error while trying to unmarshal dynamoDB result items")
	}

	for _, clientDto := range *clientsDto {
		*clients = append(*clients, *clientDto.ToModel())
	}

	d.logger.Info("Get clients finished", config, clients)
	return clients, nil
}

func (d *dynamoDBClientPersistence) abort(err error, message string) error {
	clientPersistenceError := exceptions.ClientPersistenceError(err, message)
	d.logger.Error(clientPersistenceError, "Get clients failed: "+message)
	return clientPersistenceError
}
