package integrated

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambdacontext"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/brienze1/crypto-robot-validator/internal/validator"
	"github.com/brienze1/crypto-robot-validator/internal/validator/application/config"
	"github.com/brienze1/crypto-robot-validator/internal/validator/application/properties"
	"github.com/brienze1/crypto-robot-validator/internal/validator/domain/model"
	"github.com/brienze1/crypto-robot-validator/test/mocks"
	"github.com/cucumber/godog"
	"github.com/google/uuid"
	"net/http"
	"net/http/httptest"
	"time"
)

//func TestFeatures(t *testing.T) {
//	suite := godog.TestSuite{
//		ScenarioInitializer: func(s *godog.ScenarioContext) {
//			InitializeScenario(s)
//		},
//		Options: &godog.Options{
//			Format:   "pretty",
//			Paths:    []string{"features"},
//			TestingT: t,
//		},
//	}
//
//	if suite.Run() != 0 {
//		t.Fatal("non-zero status returned, failed to run feature tests")
//	}
//}

func InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.Step(`^test env variables were loaded$`, testEnvVariablesWereLoaded)
	ctx.Step(`^dynamoDB is "([^"]*)"$`, dynamoDBIs)
	ctx.Step(`^binance api is "([^"]*)"$`, binanceApiIs)
	ctx.Step(`^sns service is "([^"]*)"$`, snsServiceIs)
	ctx.Step(`^I receive message with summary equals "([^"]*)"$`, iReceiveMessageWithSummaryEquals)
	ctx.Step(`^there are (\d+) clients available in DB`, thereAreClientsAvailableInDB)
	ctx.Step(`^handler is triggered$`, handlerIsTriggered)
	ctx.Step(`^there should be (\d+) messages sent via sns$`, thereShouldBeMessagesSentViaSns)
	ctx.Step(`^sns messages payload should have all client_id\'s got from clients table$`, snsMessagesPayloadShouldHaveAllClientIdsGotFromClientsTable)
	ctx.Step(`^sns messages payload symbol should be equal "([^"]*)"$`, snsMessagesPayloadSymbolShouldBeEqual)
	ctx.Step(`^sns messages payload operation should be equal "([^"]*)"$`, snsMessagesPayloadOperationShouldBeEqual)
	ctx.Step(`^process should exit with (\d+)$`, processShouldExitWith)
}

type (
	snsClientMock struct {
	}
	contextMock struct {
		context.Context
	}
)

var (
	persistedClients        []*model.Client
	snsClientError          error
	snsClientPublishCounter = 0
	snsClientPublishInputs  []*model.OperationRequest
	handlerError            error
)

func (s *snsClientMock) Publish(_ context.Context, input *sns.PublishInput, _ ...func(*sns.Options)) (*sns.PublishOutput, error) {
	snsClientPublishCounter++
	request := model.OperationRequest{}
	_ = json.Unmarshal([]byte(*input.Message), &request)
	snsClientPublishInputs = append(snsClientPublishInputs, &request)
	return nil, snsClientError
}

func (ctx contextMock) Value(any) any {
	return &lambdacontext.LambdaContext{
		AwsRequestID: uuid.NewString(),
	}
}

var (
	ctx   contextMock
	event *events.SQSEvent
)

func testEnvVariablesWereLoaded() {
	config.LoadTestEnv()
}

func dynamoDBIs(_ string) error {
	config.DependencyInjector().DynamoDBClient = mocks.DynamoDBClient()

	return nil
}

func binanceApiIs(status string) error {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if status != "up" {
			http.Error(w, "error test", 500)
		}

		//_, _ = w.Write(response)
	}))

	properties.Properties().BiscointGetCryptoUrl = server.URL

	return nil
}

func snsServiceIs(status string) error {
	snsClientPublishCounter = 0
	snsClientPublishInputs = []*model.OperationRequest{}
	config.DependencyInjector().SNSClient = &snsClientMock{}

	if status != "up" {
		snsClientError = errors.New("sns client not up")
	}

	return nil
}

func iReceiveMessageWithSummaryEquals(_ string) error {
	//event = createSQSEvent(summaryValue)

	ctx = contextMock{}

	return nil
}

func thereAreClientsAvailableInDB(numberOfClients int) error {
	persistedClients = []*model.Client{}
	for i := 1; i <= numberOfClients; i++ {
		client := model.Client{
			Id:           uuid.NewString(),
			Active:       true,
			LockedUntil:  time.Now().Add(-time.Second * 15),
			Locked:       false,
			CashAmount:   10000.0,
			CryptoAmount: 1.0,
			BuyOn:        1,
			SellOn:       1,
			Symbols:      []string{"BTC", "SOL"},
		}
		persistedClients = append(persistedClients, &client)
	}

	return nil
}

func handlerIsTriggered() error {
	config.LoadTestEnv()
	config.DependencyInjector().Logger = mocks.Logger()

	handlerError = validator.Main().Handle(ctx, *event)

	return nil
}

func thereShouldBeMessagesSentViaSns(numberOfMessages int) error {
	err := assertEqual(snsClientPublishCounter, numberOfMessages)
	if err != nil {
		return err
	}

	err = assertEqual(len(snsClientPublishInputs), numberOfMessages)
	if err != nil {
		return err
	}

	return err
}

func snsMessagesPayloadShouldHaveAllClientIdsGotFromClientsTable() error {
	for _, client := range persistedClients {
		found := false

		for _, request := range snsClientPublishInputs {
			err := assertEqual(request.ClientId, client.Id)
			if err == nil {
				found = true
			}
		}

		if !found {
			err := errors.New("client id should have been sent to sns")
			return err
		}
	}

	return nil
}

func snsMessagesPayloadSymbolShouldBeEqual(value string) error {
	for _, request := range snsClientPublishInputs {
		err := assertEqual(request.Symbol, value)
		if err != nil {
			return err
		}
	}

	return nil
}

func snsMessagesPayloadOperationShouldBeEqual(value string) error {
	for _, request := range snsClientPublishInputs {
		err := assertEqual(request.Operation, value)
		if err != nil {
			return err
		}
	}

	return nil
}

func processShouldExitWith(status int) error {
	if status == 0 && handlerError != nil {
		return errors.New("should have exited with status 0 but instead finished with:" + handlerError.Error())
	} else if status == 1 && handlerError == nil {
		return errors.New("should have exited with status 1 but instead finished with 0")
	}
	return nil
}

func assertEqual(val1, val2 interface{}) error {
	if val1 == val2 {
		return nil
	}
	val1String, _ := json.Marshal(val1)
	val2String, _ := json.Marshal(val2)
	return errors.New(string(val1String) + " should be equal to " + string(val2String))
}

//func createSQSEvent(summary summary.Summary) *events.SQSEvent {
//	analysisDto := dto2.AnalysisDto{
//		Summary:   summary,
//		Timestamp: time.Now().Format("2022-01-01 13:01:01"),
//	}
//
//	analysisMessage, _ := json.Marshal(analysisDto)
//
//	snsEventMessage, _ := json.Marshal(createSNSEvent(string(analysisMessage)))
//
//	return &events.SQSEvent{
//		Records: []events.SQSMessage{
//			{
//				Body: string(snsEventMessage),
//			},
//		},
//	}
//}

//func createSNSEvent(message string) events.SNSEntity {
//	return events.SNSEntity{
//		Message: message,
//	}
//}
