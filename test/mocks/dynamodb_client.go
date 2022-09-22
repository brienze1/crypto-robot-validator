package mocks

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
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

func (d *dynamoDBClient) Scan(ctx context.Context, params *dynamodb.ScanInput, optFns ...func(*dynamodb.Options)) (*dynamodb.ScanOutput, error) {
	return nil, nil
}

func (d *dynamoDBClient) GetItem(_ context.Context, params *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error) {
	client := map[string]string{}

	_ = attributevalue.UnmarshalMap(params.Key, &client)

	item := d.items[client["client_id"]]

	itemOutput, _ := attributevalue.MarshalMap(item)

	return &dynamodb.GetItemOutput{
		Item: itemOutput,
	}, nil
}

func (d *dynamoDBClient) PutItem(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {
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
