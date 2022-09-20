package log

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/brienze1/crypto-robot-operation-hub/internal/operation-hub/domain/adapters"
	"github.com/brienze1/crypto-robot-operation-hub/pkg/custom_error"
	logg "github.com/brienze1/crypto-robot-operation-hub/pkg/log"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"testing"
)

var (
	logger            adapters.LoggerAdapter
	buf               bytes.Buffer
	testCorrelationId string
	testMessage       string
	testMetadata1     interface{}
	testMetadata2     interface{}
	testError         custom_error.BaseErrorAdapter
)

func setup() {
	logger = logg.Logger()

	log.SetOutput(&buf)

	testCorrelationId = uuid.NewString()
	testMessage = uuid.NewString()
	testMetadata1 = map[string]interface{}{"Key1": 123, "Key2": "value2", "Test": uuid.NewString()}
	testMetadata2 = map[string]interface{}{"Key3": 1234, "Key4": "value3", "Test2": uuid.NewString()}
	testError = custom_error.NewBaseError(errors.New("test error message"))
}

func teardown() {
	log.SetOutput(os.Stderr)
	logger.SetCorrelationID("")
}

func TestLoggerSetCorrelationIdInfoSuccess(t *testing.T) {
	setup()
	defer func() {
		teardown()
	}()

	logger.SetCorrelationID(testCorrelationId)
	logger.Info(testMessage)

	assert.Contains(t, buf.String(), testCorrelationId)
	assert.Contains(t, buf.String(), testMessage)
}

func TestLoggerSetCorrelationIdErrorSuccess(t *testing.T) {
	setup()
	defer func() {
		teardown()
	}()

	logger.SetCorrelationID(testCorrelationId)
	logger.Error(testError, testMessage)

	assert.Contains(t, buf.String(), testCorrelationId)
	assert.Contains(t, buf.String(), testError.Error())
	assert.Contains(t, buf.String(), testMessage)
}

func TestLoggerInfoMessageOnlySuccess(t *testing.T) {
	setup()
	defer func() {
		teardown()
	}()

	logger.Info(testMessage)

	assert.Contains(t, buf.String(), testMessage)
}

func TestLoggerInfoOneMetadataSuccess(t *testing.T) {
	setup()
	defer func() {
		teardown()
	}()

	logger.Info(testMessage, testMetadata1)

	testMetadata1String, _ := json.Marshal(testMetadata1)

	assert.Contains(t, buf.String(), testMessage)
	assert.Contains(t, buf.String(), string(testMetadata1String))
}

func TestLoggerInfoTwoMetadataSuccess(t *testing.T) {
	setup()
	defer func() {
		teardown()
	}()

	logger.Info(testMessage, testMetadata1, testMetadata2)

	testMetadata1String, _ := json.Marshal(testMetadata1)
	testMetadata2String, _ := json.Marshal(testMetadata2)

	assert.Contains(t, buf.String(), testMessage)
	assert.Contains(t, buf.String(), string(testMetadata1String))
	assert.Contains(t, buf.String(), string(testMetadata2String))
}

func TestLoggerErrorMessageOnlySuccess(t *testing.T) {
	setup()
	defer func() {
		teardown()
	}()

	logger.Error(testError, testMessage)

	assert.Contains(t, buf.String(), testError.Error())
	assert.Contains(t, buf.String(), testMessage)
}

func TestLoggerErrorOneMetadataSuccess(t *testing.T) {
	setup()
	defer func() {
		teardown()
	}()

	logger.Error(testError, testMessage, testMetadata1)

	testMetadata1String, _ := json.Marshal(testMetadata1)

	assert.Contains(t, buf.String(), testError.Error())
	assert.Contains(t, buf.String(), testMessage)
	assert.Contains(t, buf.String(), string(testMetadata1String))
}

func TestLoggerErrorTwoMetadataSuccess(t *testing.T) {
	setup()
	defer func() {
		teardown()
	}()

	logger.Error(testError, testMessage, testMetadata1, testMetadata2)

	testMetadata1String, _ := json.Marshal(testMetadata1)
	testMetadata2String, _ := json.Marshal(testMetadata2)

	assert.Contains(t, buf.String(), testError.Error())
	assert.Contains(t, buf.String(), testMessage)
	assert.Contains(t, buf.String(), string(testMetadata1String))
	assert.Contains(t, buf.String(), string(testMetadata2String))
}
