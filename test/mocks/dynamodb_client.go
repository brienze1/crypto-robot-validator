package mocks

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/brienze1/crypto-robot-validator/internal/validator/integration/dto"
	"github.com/brienze1/crypto-robot-validator/internal/validator/integration/exceptions"
)

type dynamoDBClient struct {
	ScanCounter    int
	ScanError      error
	ScanOutput     *dynamodb.ScanOutput
	GetItemCounter int
	GetItemError   error
	GetItemOutput  *dynamodb.GetItemOutput
	PutItemCounter int
	PutItemError   error
	PutItemOutput  *dynamodb.PutItemOutput
	items          map[string]interface{}
}

func DynamoDBClient() *dynamoDBClient {
	return &dynamoDBClient{
		items: map[string]interface{}{},
	}
}

func (d *dynamoDBClient) Scan(_ context.Context, _ *dynamodb.ScanInput, _ ...func(*dynamodb.Options)) (*dynamodb.ScanOutput, error) {
	return nil, nil
}

func (d *dynamoDBClient) GetItem(_ context.Context, params *dynamodb.GetItemInput, _ ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error) {
	d.GetItemCounter++

	if d.GetItemError != nil {
		return nil, exceptions.DynamoDBClientPersistenceError(d.GetItemError, "GetItem error")
	}

	client := map[string]string{}

	_ = attributevalue.UnmarshalMap(params.Key, &client)

	item := d.items[client["client_id"]]

	var itemOutput map[string]types.AttributeValue

	if item == nil {
		itemOutput = nil
	} else {
		itemOutput, _ = attributevalue.MarshalMap(item)
	}

	return &dynamodb.GetItemOutput{
		Item: itemOutput,
	}, nil
}

func (d *dynamoDBClient) PutItem(_ context.Context, params *dynamodb.PutItemInput, _ ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {
	d.PutItemCounter++

	if d.PutItemError != nil {
		return nil, exceptions.DynamoDBClientPersistenceError(d.PutItemError, "PutItem error")
	}

	client := &dto.Client{}

	_ = attributevalue.UnmarshalMap(params.Item, &client)

	d.AddItem(client.Id, client)

	return nil, nil
}

func (d *dynamoDBClient) AddItem(key string, value interface{}) {
	d.items[key] = value
}

func (d *dynamoDBClient) Reset() {
	d.ScanCounter = 0
	d.ScanError = nil
	d.ScanOutput = &dynamodb.ScanOutput{}
	d.GetItemCounter = 0
	d.GetItemError = nil
	d.GetItemOutput = &dynamodb.GetItemOutput{}
	d.PutItemCounter = 0
	d.PutItemError = nil
	d.PutItemOutput = &dynamodb.PutItemOutput{}
	d.items = map[string]interface{}{}
}
