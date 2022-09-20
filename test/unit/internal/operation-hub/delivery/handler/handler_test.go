package handler

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambdacontext"
	"github.com/brienze1/crypto-robot-operation-hub/internal/operation-hub/delivery/adapters"
	"github.com/brienze1/crypto-robot-operation-hub/internal/operation-hub/delivery/dto"
	"github.com/brienze1/crypto-robot-operation-hub/internal/operation-hub/delivery/handler"
	adapters2 "github.com/brienze1/crypto-robot-operation-hub/internal/operation-hub/domain/adapters"
	"github.com/brienze1/crypto-robot-operation-hub/internal/operation-hub/domain/enum/summary"
	"github.com/brienze1/crypto-robot-operation-hub/pkg/custom_error"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type (
	operationUseCaseMock struct {
		adapters2.OperationUseCaseAdapter
	}
	loggerMock struct {
		adapters2.LoggerAdapter
	}
	ctx struct {
		context.Context
	}
)

var (
	clientActionsUseCaseCallCounter int
	operationUseCaseError           error
	loggerInfoCallCounter           int
	loggerErrorCallCounter          int
	awsRequestIdExpected            string
	awsRequestIdReceived            string
)

var (
	clientActionsUseCase = operationUseCaseMock{}
	logger               = loggerMock{}
	handlerImpl          adapters.HandlerAdapter
)

func (operationUseCaseMock operationUseCaseMock) TriggerOperations(summary.Summary) error {
	clientActionsUseCaseCallCounter += 1
	return operationUseCaseError
}

func (loggerMock loggerMock) SetCorrelationID(id string) {
	awsRequestIdReceived = id
}

func (loggerMock loggerMock) Info(string, ...interface{}) {
	loggerInfoCallCounter++
}

func (loggerMock loggerMock) Error(error, string, ...interface{}) {
	loggerErrorCallCounter++
}

func (ctx ctx) Value(any) any {
	return &lambdacontext.LambdaContext{
		AwsRequestID: awsRequestIdExpected,
	}
}

func setup() {
	handlerImpl = handler.Handler(clientActionsUseCase, logger)

	clientActionsUseCaseCallCounter = 0
	operationUseCaseError = nil
	loggerInfoCallCounter = 0
	loggerErrorCallCounter = 0
	awsRequestIdReceived = ""
	awsRequestIdExpected = uuid.NewString()
}

func TestHandlerSuccess(t *testing.T) {
	setup()

	ctx := ctx{}
	event := *createSQSEvent()

	err := handlerImpl.Handle(ctx, event)

	assert.Nil(t, err, "Error should be nil")
	assert.Equal(t, 1, clientActionsUseCaseCallCounter, "clientActionsUseCase should be called once")
	assert.Equal(t, 2, loggerInfoCallCounter, "logger info should be called twice")
	assert.Equal(t, 0, loggerErrorCallCounter, "logger exceptions should not be called")
	assert.Equal(t, awsRequestIdExpected, awsRequestIdReceived, "Logger correlationId is same as context awsRequestId")
}

func TestHandlerJsonSQSError(t *testing.T) {
	setup()

	ctx := ctx{}
	event := *createSQSEvent()
	event.Records[0].Body = ""

	var err = handlerImpl.Handle(ctx, event)

	assert.NotNil(t, err, "Error should not be nil")
	assert.Equal(t, "Error while trying to parse the SNS message", err.(custom_error.BaseErrorAdapter).InternalError())
	assert.Equal(t, "Error occurred while handling the event", err.(custom_error.BaseErrorAdapter).Description())
	assert.Equal(t, "unexpected end of JSON input", err.(custom_error.BaseErrorAdapter).Error())
	assert.Equal(t, 0, clientActionsUseCaseCallCounter, "clientActionsUseCase should not be called")
	assert.Equal(t, 1, loggerInfoCallCounter, "logger info should be called once")
	assert.Equal(t, 1, loggerErrorCallCounter, "logger exceptions should be called once")
	assert.Equal(t, awsRequestIdExpected, awsRequestIdReceived, "Logger correlationId is same as context awsRequestId")
}

func TestHandlerJsonSNSError(t *testing.T) {
	setup()

	ctx := ctx{}
	event := *createSQSEvent()
	snsMessage := createSNSEvent("")
	snsMessageString, _ := json.Marshal(snsMessage)
	event.Records[0].Body = string(snsMessageString)

	var err = handlerImpl.Handle(ctx, event)

	assert.NotNil(t, err, "Error should not be nil")
	assert.Equal(t, "Error while trying to parse the analysis object", err.(custom_error.BaseErrorAdapter).InternalError())
	assert.Equal(t, "Error occurred while handling the event", err.(custom_error.BaseErrorAdapter).Description())
	assert.Equal(t, "unexpected end of JSON input", err.(custom_error.BaseErrorAdapter).Error())
	assert.Equal(t, 0, clientActionsUseCaseCallCounter, "clientActionsUseCase should not be called")
	assert.Equal(t, 1, loggerInfoCallCounter, "logger info should be called once")
	assert.Equal(t, 1, loggerErrorCallCounter, "logger exceptions should be called once")
	assert.Equal(t, awsRequestIdExpected, awsRequestIdReceived, "Logger correlationId is same as context awsRequestId")
}

func TestHandlerOperationUseCaseError(t *testing.T) {
	setup()

	ctx := ctx{}
	event := *createSQSEvent()
	operationUseCaseError = errors.New(uuid.NewString())

	var err = handlerImpl.Handle(ctx, event)

	assert.NotNil(t, err, "Error should not be nil")
	assert.Equal(t, "Error while trying to run OperationUseCase", err.(custom_error.BaseErrorAdapter).InternalError())
	assert.Equal(t, "Error occurred while handling the event", err.(custom_error.BaseErrorAdapter).Description())
	assert.Equal(t, operationUseCaseError.Error(), err.(custom_error.BaseErrorAdapter).Error())
	assert.Equal(t, 1, clientActionsUseCaseCallCounter, "clientActionsUseCase should not be called")
	assert.Equal(t, 1, loggerInfoCallCounter, "logger info should be called once")
	assert.Equal(t, 1, loggerErrorCallCounter, "logger exceptions should be called once")
	assert.Equal(t, awsRequestIdExpected, awsRequestIdReceived, "Logger correlationId is same as context awsRequestId")
}

func createSQSEvent() *events.SQSEvent {
	analysisDto := dto.AnalysisDto{
		Summary:   summary.StrongBuy,
		Timestamp: time.Now().Format("2022-01-01 13:01:01"),
	}

	analysisMessage, _ := json.Marshal(analysisDto)

	snsEventMessage, _ := json.Marshal(createSNSEvent(string(analysisMessage)))

	return &events.SQSEvent{
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
