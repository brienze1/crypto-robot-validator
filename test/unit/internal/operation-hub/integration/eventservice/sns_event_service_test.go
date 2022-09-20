package eventservice

import (
	"context"
	"encoding/json"
	"errors"
	sns2 "github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/brienze1/crypto-robot-operation-hub/internal/operation-hub/application/config"
	"github.com/brienze1/crypto-robot-operation-hub/internal/operation-hub/application/properties"
	"github.com/brienze1/crypto-robot-operation-hub/internal/operation-hub/domain/adapters"
	adapters2 "github.com/brienze1/crypto-robot-operation-hub/internal/operation-hub/integration/adapters"
	"github.com/brienze1/crypto-robot-operation-hub/internal/operation-hub/integration/eventservice"
	"github.com/brienze1/crypto-robot-operation-hub/pkg/custom_error"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

type (
	loggerMock struct {
		adapters.LoggerAdapter
	}
	snsMock struct {
		adapters2.SNSAdapter
	}
)

var (
	logger = loggerMock{}
	sns    = snsMock{}
)

var (
	loggerInfoCounter  int
	loggerErrorCounter int
	snsPublishCounter  int
	snsPublishInput    sns2.PublishInput
)

var (
	snsPublishError error
)

func (l loggerMock) Info(string, ...interface{}) {
	loggerInfoCounter++
}

func (l loggerMock) Error(error, string, ...interface{}) {
	loggerErrorCounter++
}

func (s snsMock) Publish(_ context.Context, input *sns2.PublishInput, _ ...func(*sns2.Options)) (*sns2.PublishOutput, error) {
	snsPublishCounter++
	snsPublishInput = *input

	return nil, snsPublishError
}

var (
	snsEventService adapters.EventServiceAdapter
	payload         interface{}
	payloadString   string
)

func setup() {
	config.LoadTestEnv()

	loggerInfoCounter = 0
	loggerErrorCounter = 0
	snsPublishCounter = 0
	snsPublishInput = sns2.PublishInput{}
	snsPublishError = nil

	snsEventService = eventservice.SNSEventService(logger, sns)

	payload = map[string]string{
		"message": uuid.NewString(),
		"test":    uuid.NewString(),
		"test2":   uuid.NewString(),
		"test3":   uuid.NewString(),
	}
	payloadByte, _ := json.Marshal(payload)
	payloadString = string(payloadByte)
}

func TestSendSuccess(t *testing.T) {
	setup()

	err := snsEventService.Send(payload)

	assert.Nil(t, err)
	assert.Equal(t, 1, snsPublishCounter)
	assert.Equal(t, payloadString, *snsPublishInput.Message)
	assert.Equal(t, properties.Properties().CryptoOperationTriggerTopicArn, *snsPublishInput.TopicArn)
	assert.Equal(t, 2, loggerInfoCounter)
	assert.Equal(t, 0, loggerErrorCounter)
}

func TestSendPublishFailure(t *testing.T) {
	setup()

	snsPublishError = errors.New("test error")

	err := snsEventService.Send(payload)

	assert.Equal(t, "test error", err.(custom_error.BaseErrorAdapter).Error())
	assert.Equal(t, "Error while trying to publish", err.(custom_error.BaseErrorAdapter).InternalError())
	assert.Equal(t, "Error while publishing SNS event", err.(custom_error.BaseErrorAdapter).Description())
	assert.Equal(t, 1, snsPublishCounter)
	assert.Equal(t, payloadString, *snsPublishInput.Message)
	assert.Equal(t, properties.Properties().CryptoOperationTriggerTopicArn, *snsPublishInput.TopicArn)
	assert.Equal(t, 1, loggerInfoCounter)
	assert.Equal(t, 1, loggerErrorCounter)
}

func TestSendMarshalFailure(t *testing.T) {
	setup()

	payload = make(chan int)

	err := snsEventService.Send(payload)

	assert.Equal(t, "json: unsupported type: chan int", err.(custom_error.BaseErrorAdapter).Error())
	assert.Equal(t, "Error while trying create string message", err.(custom_error.BaseErrorAdapter).InternalError())
	assert.Equal(t, "Error while publishing SNS event", err.(custom_error.BaseErrorAdapter).Description())
	assert.Equal(t, 0, snsPublishCounter)
	assert.Equal(t, 1, loggerInfoCounter)
	assert.Equal(t, 1, loggerErrorCounter)
}
