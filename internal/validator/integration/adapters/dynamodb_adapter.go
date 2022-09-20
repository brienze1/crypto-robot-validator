package adapters

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

// DynamoDBAdapter is an adapter class for AWS DynamoDB v2 SDK.
type DynamoDBAdapter interface {
	// Scan is the same as dynamodb.Client Scan method
	Scan(ctx context.Context, params *dynamodb.ScanInput, optFns ...func(*dynamodb.Options)) (*dynamodb.ScanOutput, error)

	// GetItem is the same as dynamodb.Client GetItem method
	GetItem(ctx context.Context, params *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error)

	// PutItem is the same as dynamodb.Client PutItem method
	PutItem(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error)
}
