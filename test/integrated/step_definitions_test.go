package integrated

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambdacontext"
	"github.com/brienze1/crypto-robot-validator/internal/validator"
	"github.com/brienze1/crypto-robot-validator/internal/validator/application/config"
	"github.com/brienze1/crypto-robot-validator/internal/validator/application/properties"
	"github.com/brienze1/crypto-robot-validator/internal/validator/integration/dto"
	"github.com/brienze1/crypto-robot-validator/test/mocks"
	"github.com/cucumber/godog"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
)

func TestFeatures(test *testing.T) {
	t = test
	suite := godog.TestSuite{
		ScenarioInitializer: func(s *godog.ScenarioContext) {
			InitializeScenario(s)
		},
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"features"},
			TestingT: t,
		},
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run feature tests")
	}
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.Step(`^test env variables were loaded$`, testEnvVariablesWereLoaded)
	ctx.Step(`^dynamoDB is up$`, dynamoDBIsUp)
	ctx.Step(`^redis is up$`, redisIsUp)
	ctx.Step(`^biscoint api is up$`, biscointApiIsUp)
	ctx.Step(`^sns service is up$`, snsServiceIsUp)
	ctx.Step(`^secrets manager service is up$`, secretsManagerServiceIsUp)
	ctx.Step(`^there is a client available on DynamoDB with client id "([^"]*)"$`, thereIsAClientAvailableOnDynamoDBWithClientId)
	ctx.Step(`^client available "([^"]*)" balance is (\d+)\.(\d+)$`, clientAvailableBalanceIs)
	ctx.Step(`^client reserved "([^"]*)" balance is (\d+)\.(\d+)$`, clientReservedBalanceIs)
	ctx.Step(`^client "([^"]*)" balance is (\d+)\.(\d+) on biscoint$`, clientBalanceIsOnBiscoint)
	ctx.Step(`^crypto current "([^"]*)" value is (\d+)\.(\d+) on biscoint$`, cryptoCurrentValueIsOnBiscoint)
	ctx.Step(`^the following credentials available for client id "([^"]*)"$`, theFollowingCredentialsAvailableForClientId)
	ctx.Step(`^the following message is received$`, theFollowingMessageIsReceived)
	ctx.Step(`^there should be (\d+) messages sent via sns$`, thereShouldBeMessagesSentViaSns)
	ctx.Step(`^process should exit with (\d+)$`, processShouldExitWith)
}

var (
	dynamoDB             = mocks.DynamoDBClient()
	biscointApi          = mocks.HttpClient()
	snsClient            = mocks.SNSClient()
	secretsManagerClient = mocks.SecretsManager()
)

var (
	t         *testing.T
	client    *dto.Client
	balance   *dto.BalanceResponse
	coin      *dto.CoinResponse
	handleErr error
)

func testEnvVariablesWereLoaded() error {
	config.LoadTestEnv()
	return nil
}

func dynamoDBIsUp() error {
	config.DependencyInjector().DynamoDBClient = dynamoDB
	return nil
}

func redisIsUp() error {
	config.DependencyInjector().RedisClient = mocks.RedisServer()
	return nil
}

func biscointApiIsUp() error {
	balance = &dto.BalanceResponse{
		Balance: dto.Balance{
			BRL: "0.0",
			BTC: "0.0",
		},
	}

	coin = &dto.CoinResponse{
		Message: "",
		Coin: dto.Coin{
			Symbol:    "BTC",
			Quote:     "BRL",
			BuyValue:  100000.00,
			SellValue: 99000.00,
		},
	}

	biscointApi.SetupServer()
	balanceResponse, _ := json.Marshal(balance)
	coinResponse, _ := json.Marshal(coin)
	biscointApi.GetBalanceResponse = string(balanceResponse)
	biscointApi.GetCryptoResponse = string(coinResponse)
	properties.Properties().BiscointUrl = biscointApi.GetUrl() + "/"
	return nil
}

func snsServiceIsUp() error {
	config.DependencyInjector().SNSClient = snsClient
	return nil
}

func secretsManagerServiceIsUp() error {
	encryptionSecret := &dto.EncryptionSecrets{
		EncryptionKey: "9y$B?E(H+MbQeThWmZq4t7w!z%C*F)J@",
	}
	secretsManagerClient.SetSecret(properties.Properties().Aws.SecretsManager.EncryptionSecretName, encryptionSecret)
	config.DependencyInjector().SecretsManager = secretsManagerClient
	return nil
}

func thereIsAClientAvailableOnDynamoDBWithClientId(clientId string) error {
	client = &dto.Client{
		Id:      clientId,
		Active:  true,
		Symbols: []string{"BTC"},
	}

	dynamoDB.AddItem(clientId, client, properties.Properties().Aws.DynamoDB.ClientTableName)
	return nil
}

func clientAvailableBalanceIs(balanceType string, value float64) error {
	if balanceType == "brl" {
		client.CashAvailable = value
		dynamoDB.AddItem(client.Id, client, properties.Properties().Aws.DynamoDB.ClientTableName)
	} else if balanceType == "btc" {
		client.CryptoAvailable = value
		dynamoDB.AddItem(client.Id, client, properties.Properties().Aws.DynamoDB.ClientTableName)
	}
	return nil
}

func clientReservedBalanceIs(balanceType string, value float64) error {
	if balanceType == "brl" {
		client.CashReserved = value
		dynamoDB.AddItem(client.Id, client, properties.Properties().Aws.DynamoDB.ClientTableName)
	} else if balanceType == "btc" {
		client.CryptoReserved = value
		dynamoDB.AddItem(client.Id, client, properties.Properties().Aws.DynamoDB.ClientTableName)
	}
	return nil
}

func clientBalanceIsOnBiscoint(balanceType string, value float64) error {
	if balanceType == "brl" {
		balance.Balance.BRL = strconv.FormatFloat(value, 'f', 2, 64)
		balanceResponse, _ := json.Marshal(balance)
		biscointApi.GetBalanceResponse = string(balanceResponse)
	} else if balanceType == "btc" {
		balance.Balance.BTC = strconv.FormatFloat(value, 'f', 8, 64)
		balanceResponse, _ := json.Marshal(balance)
		biscointApi.GetBalanceResponse = string(balanceResponse)
	}
	return nil
}

func cryptoCurrentValueIsOnBiscoint(operationType string, value float64) error {
	if operationType == "buy" {
		coin.Coin.BuyValue = value
		coinResponse, _ := json.Marshal(coin)
		biscointApi.GetCryptoResponse = string(coinResponse)
	} else if operationType == "sell" {
		coin.Coin.SellValue = value
		coinResponse, _ := json.Marshal(coin)
		biscointApi.GetCryptoResponse = string(coinResponse)
	}
	return nil
}

func theFollowingCredentialsAvailableForClientId(clientId string, credentialsJson *godog.DocString) error {
	var credentials dto.Credentials
	_ = json.Unmarshal([]byte(credentialsJson.Content), &credentials)
	dynamoDB.AddItem(clientId, credentials, properties.Properties().Aws.DynamoDB.CredentialsTableName)
	return nil
}

func theFollowingMessageIsReceived(messageReceived *godog.DocString) error {
	event := createSQSEvent(messageReceived.Content)
	ctx := createContext()

	return validator.Main().Handle(ctx, event)
}

func thereShouldBeMessagesSentViaSns(numberOfMessages int) error {
	assert.Equal(t, numberOfMessages, snsClient.NumberOfMessagesSent)
	return nil
}

func processShouldExitWith(status int) error {
	if status == 0 {
		assert.Nil(t, handleErr)
	} else if status == 1 {
		assert.NotNil(t, handleErr)
	}
	return nil
}

func createSQSEvent(message string) events.SQSEvent {
	snsEventMessage, _ := json.Marshal(createSNSEvent(message))

	return events.SQSEvent{
		Records: []events.SQSMessage{
			{
				Body: string(snsEventMessage),
			},
		},
	}
}

func createSNSEvent(message string) events.SNSEntity {
	return events.SNSEntity{
		Message: message,
	}
}

type ctx struct {
	context.Context
	awsRequestId string
}

func (ctx ctx) Value(any) any {
	return &lambdacontext.LambdaContext{
		AwsRequestID: ctx.awsRequestId,
	}
}

func createContext() *ctx {
	return &ctx{
		awsRequestId: uuid.NewString(),
	}
}
