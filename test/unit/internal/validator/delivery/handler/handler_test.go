package handler

import (
	"context"
	"errors"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambdacontext"
	"github.com/brienze1/crypto-robot-validator/internal/validator/delivery/adapters"
	"github.com/brienze1/crypto-robot-validator/internal/validator/delivery/handler"
	"github.com/brienze1/crypto-robot-validator/pkg/custom_error"
	"github.com/brienze1/crypto-robot-validator/test/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

type (
	ctx struct {
		context.Context
	}
)

var (
	awsRequestIdExpected string
)

var (
	validationUseCase = mocks.ValidationUseCase()
	logger            = mocks.Logger()
	handlerImpl       adapters.HandlerAdapter
)

func (c ctx) Value(any) any {
	return &lambdacontext.LambdaContext{
		AwsRequestID: awsRequestIdExpected,
	}
}

func setup() {
	handlerImpl = handler.Handler(validationUseCase, logger)

	logger.Reset()
	validationUseCase.Reset()
	awsRequestIdExpected = uuid.NewString()
}

func TestHandlerSuccess(t *testing.T) {
	setup()

	ctx := ctx{}
	event := *createSQSEvent()

	err := handlerImpl.Handle(ctx, event)

	assert.Nil(t, err, "Error should be nil")
	assert.Equal(t, 1, validationUseCase.ValidateCallCounter, "validate should be called once")
	assert.Equal(t, 2, logger.InfoCallCounter, "logger info should be called twice")
	assert.Equal(t, 0, logger.ErrorCallCounter, "logger exceptions should not be called")
	assert.Equal(t, awsRequestIdExpected, logger.CorrelationId, "Logger correlationId is same as context awsRequestId")
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
	assert.Equal(t, 0, validationUseCase.ValidateCallCounter, "validationUseCase should not be called")
	assert.Equal(t, 1, logger.InfoCallCounter, "logger info should be called once")
	assert.Equal(t, 1, logger.ErrorCallCounter, "logger exceptions should be called once")
	assert.Equal(t, awsRequestIdExpected, logger.CorrelationId, "Logger correlationId is same as context awsRequestId")
}

func TestHandlerOperationUseCaseError(t *testing.T) {
	setup()

	ctx := ctx{}
	event := *createSQSEvent()
	expectedErrorMsg := uuid.NewString()
	validationUseCase.ValidateError = errors.New(expectedErrorMsg)

	var err = handlerImpl.Handle(ctx, event)

	assert.NotNil(t, err, "Error should not be nil")
	assert.Equal(t, "Error while trying to run ValidationUseCase", err.(custom_error.BaseErrorAdapter).InternalError())
	assert.Equal(t, "Error occurred while handling the event", err.(custom_error.BaseErrorAdapter).Description())
	assert.Equal(t, expectedErrorMsg, err.Error())
	assert.Equal(t, 1, validationUseCase.ValidateCallCounter, "validationUseCase should not be called")
	assert.Equal(t, 1, logger.InfoCallCounter, "logger info should be called once")
	assert.Equal(t, 1, logger.ErrorCallCounter, "logger exceptions should be called once")
	assert.Equal(t, awsRequestIdExpected, logger.CorrelationId, "Logger correlationId is same as context awsRequestId")
}

func createSQSEvent() *events.SQSEvent {
	operationRequest := `{
	  "client_id": "aa324edf-99fa-4a95-b9c4-a588d1ccb441e",
	  "operation": "BUY",
	  "symbol": "BTC",
	  "start_time": "2022-09-17T12:05:07.45066-03:00"
	}`

	return &events.SQSEvent{
		Records: []events.SQSMessage{
			{
				Body: operationRequest,
			},
		},
	}
}
