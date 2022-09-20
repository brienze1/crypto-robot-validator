package persistence

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	config2 "github.com/brienze1/crypto-robot-operation-hub/internal/operation-hub/application/config"
	"github.com/brienze1/crypto-robot-operation-hub/internal/operation-hub/application/properties"
	"github.com/brienze1/crypto-robot-operation-hub/internal/operation-hub/domain/adapters"
	"github.com/brienze1/crypto-robot-operation-hub/internal/operation-hub/domain/enum/summary"
	"github.com/brienze1/crypto-robot-operation-hub/internal/operation-hub/domain/model"
	adapters2 "github.com/brienze1/crypto-robot-operation-hub/internal/operation-hub/integration/adapters"
	"github.com/brienze1/crypto-robot-operation-hub/internal/operation-hub/integration/dto"
	"github.com/brienze1/crypto-robot-operation-hub/internal/operation-hub/integration/persistence"
	"github.com/brienze1/crypto-robot-operation-hub/pkg/custom_error"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type (
	loggerMock struct {
		adapters.LoggerAdapter
	}
	dynamoDBMock struct {
		adapters2.DynamoDBAdapter
	}
	timeMock struct {
		adapters.TimeAdapter
	}
)

var (
	logger     = loggerMock{}
	dynamoDB   = dynamoDBMock{}
	timeSource = timeMock{}
)

var (
	loggerInfoCounter   int
	loggerErrorCounter  int
	dynamoDBScanCounter int
	dynamoDBScanOutput  dynamodb.ScanOutput
)

var (
	dynamoDBScanError error
)

func (l loggerMock) Info(string, ...interface{}) {
	loggerInfoCounter++
}

func (l loggerMock) Error(error, string, ...interface{}) {
	loggerErrorCounter++
}

func (d dynamoDBMock) Scan(_ context.Context, params *dynamodb.ScanInput, _ ...func(*dynamodb.Options)) (*dynamodb.ScanOutput, error) {
	dynamoDBScanCounter++

	tableName = *params.TableName
	filterExpression = *params.FilterExpression
	projectionExpression = *params.ProjectionExpression
	expressionAttributeNames = params.ExpressionAttributeNames
	expressionAttributeValuesString, _ := json.Marshal(params.ExpressionAttributeValues)
	expressionAttributeValues = string(expressionAttributeValuesString)

	return &dynamoDBScanOutput, dynamoDBScanError
}

func (t timeMock) Now() time.Time {
	return recordedTime
}

var (
	clientPersistence                 adapters.ClientPersistenceAdapter
	config                            model.ClientSearchConfig
	recordedTime                      time.Time
	expectedTableName                 string
	tableName                         string
	expectedClientList                []model.Client
	expectedFilterExpression          string
	filterExpression                  string
	expectedProjectionExpression      string
	projectionExpression              string
	expectedExpressionAttributeNames  map[string]string
	expressionAttributeNames          map[string]string
	expectedExpressionAttributeValues string
	expressionAttributeValues         string
)

func setup() {
	config2.LoadTestEnv()

	loggerInfoCounter = 0
	loggerErrorCounter = 0
	dynamoDBScanCounter = 0
	dynamoDBScanError = nil

	expectedTableName = uuid.NewString()
	properties.Properties().Aws.DynamoDB.ClientTableName = expectedTableName

	clientPersistence = persistence.DynamoDBClientPersistence(logger, dynamoDB, timeSource)

	config = model.ClientSearchConfig{
		Active:        true,
		Locked:        true,
		MinimumCash:   19.0123,
		MinimumCrypto: 0.1123,
		SellWeight:    summary.StrongSell,
		BuyWeight:     summary.StrongBuy,
	}

	recordedTime = time.Now()

	expr, _ := expression.NewBuilder().WithFilter(
		expression.And(
			expression.Name("active").Equal(expression.Value(config.Active)),
			expression.Name("locked").Equal(expression.Value(config.Locked)),
			expression.Name("locked_until").LessThanEqual(expression.Value(recordedTime)),
			expression.Name("cash_amount").GreaterThanEqual(expression.Value(config.MinimumCash)),
			expression.Name("crypto_amount").GreaterThanEqual(expression.Value(config.MinimumCrypto)),
			expression.Name("sell_on").LessThanEqual(expression.Value(config.SellWeight.Value())),
			expression.Name("buy_on").LessThanEqual(expression.Value(config.BuyWeight.Value())),
			expression.Name("symbols").Contains(config.Symbol.Name()),
		),
	).Build()

	expectedFilterExpression = *expr.Filter()
	expectedProjectionExpression = *aws.String("client_id")
	expectedExpressionAttributeNames = expr.Names()
	expectedExpressionAttributeValuesString, _ := json.Marshal(expr.Values())
	expectedExpressionAttributeValues = string(expectedExpressionAttributeValuesString)

	expectedClientList = []model.Client{}
	var items []map[string]types.AttributeValue

	for i := 0; i < 10; i++ {
		expectedClient := model.Client{
			Id: uuid.NewString(),
		}

		expectedClientList = append(
			expectedClientList,
			expectedClient,
		)

		item, _ := attributevalue.MarshalMap(&dto.Client{Id: expectedClient.Id})

		items = append(items, item)
	}

	dynamoDBScanOutput = dynamodb.ScanOutput{
		Items: items,
	}
}

func TestGetClientsSuccess(t *testing.T) {
	setup()

	clients, err := clientPersistence.GetClients(&config)

	assert.Nilf(t, err, "Should be nil")
	assert.NotNilf(t, clients, "Should not be nil")
	assert.Greaterf(t, len(*clients), 0, "Client list should not be empty")
	assert.Equal(t, expectedClientList, *clients)
	assert.Equal(t, expectedTableName, tableName)
	assert.Equal(t, expectedFilterExpression, filterExpression)
	assert.Equal(t, expectedProjectionExpression, projectionExpression)
	assert.Equal(t, expectedExpressionAttributeNames, expressionAttributeNames)
	assert.Equal(t, expectedExpressionAttributeValues, expressionAttributeValues)
	assert.Equal(t, 1, dynamoDBScanCounter)
	assert.Equal(t, 2, loggerInfoCounter)
	assert.Equal(t, 0, loggerErrorCounter)
}

func TestGetClientsScanFailure(t *testing.T) {
	setup()

	dynamoDBScanError = errors.New("test error")

	clients, err := clientPersistence.GetClients(&config)

	assert.Equal(t, "test error", err.(custom_error.BaseErrorAdapter).Error())
	assert.Equal(t, "DynamoDb scan error", err.(custom_error.BaseErrorAdapter).InternalError())
	assert.Equal(t, "Error while getting data from Client table", err.(custom_error.BaseErrorAdapter).Description())
	assert.Equal(t, expectedTableName, tableName)
	assert.Equal(t, expectedFilterExpression, filterExpression)
	assert.Equal(t, expectedProjectionExpression, projectionExpression)
	assert.Equal(t, expectedExpressionAttributeNames, expressionAttributeNames)
	assert.Equal(t, expectedExpressionAttributeValues, expressionAttributeValues)
	assert.Nilf(t, clients, "Should be nil")
	assert.Equal(t, 1, dynamoDBScanCounter)
	assert.Equal(t, 1, loggerInfoCounter)
	assert.Equal(t, 1, loggerErrorCounter)
}

func TestGetClientsUnmarshallFailure(t *testing.T) {
	setup()

	dynamoDBScanOutput = dynamodb.ScanOutput{
		Items: []map[string]types.AttributeValue{
			{
				"client_id": &types.AttributeValueMemberBOOL{
					Value: false,
				},
			},
		},
	}

	clients, err := clientPersistence.GetClients(&config)

	assert.NotNilf(t, err, "Should not be nil")
	assert.Equal(t, "unmarshal failed, cannot unmarshal bool into Go value type string", err.(custom_error.BaseErrorAdapter).Error())
	assert.Equal(t, "Error while trying to unmarshal dynamoDB result items", err.(custom_error.BaseErrorAdapter).InternalError())
	assert.Equal(t, "Error while getting data from Client table", err.(custom_error.BaseErrorAdapter).Description())
	assert.Equal(t, expectedTableName, tableName)
	assert.Equal(t, expectedFilterExpression, filterExpression)
	assert.Equal(t, expectedProjectionExpression, projectionExpression)
	assert.Equal(t, expectedExpressionAttributeNames, expressionAttributeNames)
	assert.Equal(t, expectedExpressionAttributeValues, expressionAttributeValues)
	assert.Nilf(t, clients, "Should be nil")
	assert.Equal(t, 1, dynamoDBScanCounter)
	assert.Equal(t, 1, loggerInfoCounter)
	assert.Equal(t, 1, loggerErrorCounter)
}
