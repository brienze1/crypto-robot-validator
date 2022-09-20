package domain

import (
	"errors"
	"github.com/brienze1/crypto-robot-operation-hub/internal/operation-hub/application/config"
	"github.com/brienze1/crypto-robot-operation-hub/internal/operation-hub/domain/adapters"
	"github.com/brienze1/crypto-robot-operation-hub/internal/operation-hub/domain/enum/summary"
	"github.com/brienze1/crypto-robot-operation-hub/internal/operation-hub/domain/enum/symbol"
	"github.com/brienze1/crypto-robot-operation-hub/internal/operation-hub/domain/model"
	"github.com/brienze1/crypto-robot-operation-hub/internal/operation-hub/domain/usecase"
	"github.com/brienze1/crypto-robot-operation-hub/pkg/custom_error"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

type (
	loggerMock struct {
		adapters.LoggerAdapter
	}
	cryptoWebServiceMock struct {
		adapters.CryptoWebServiceAdapter
	}
	clientPersistenceMock struct {
		adapters.ClientPersistenceAdapter
	}
	eventServiceMock struct {
		adapters.EventServiceAdapter
	}
)

var (
	loggerInfoCounter                         int
	loggerErrorCounter                        int
	cryptoServiceGetCryptoCurrentQuoteCounter int
	clientPersistenceGetClientsCounter        int
	eventServiceSendCounter                   int
)

var (
	cryptoServiceGetCryptoCurrentQuoteError error
	clientPersistenceGetClientsError        error
	eventServiceSendError                   error
)

func (l loggerMock) Info(string, ...interface{}) {
	loggerInfoCounter++
}

func (l loggerMock) Error(error, string, ...interface{}) {
	loggerErrorCounter++
}
func (c cryptoWebServiceMock) GetCryptoCurrentQuote(symbol.Symbol) (float64, error) {
	cryptoServiceGetCryptoCurrentQuoteCounter++
	return 0, cryptoServiceGetCryptoCurrentQuoteError
}

func (c clientPersistenceMock) GetClients(*model.ClientSearchConfig) (*[]model.Client, error) {
	clientPersistenceGetClientsCounter++
	return &clients, clientPersistenceGetClientsError
}

func (e eventServiceMock) Send(interface{}) error {
	eventServiceSendCounter++
	return eventServiceSendError
}

var (
	operationUseCase  adapters.OperationUseCaseAdapter
	logger            = loggerMock{}
	cryptoWebService  = cryptoWebServiceMock{}
	clientPersistence = clientPersistenceMock{}
	eventService      = eventServiceMock{}
)

var (
	clients []model.Client
)

func setup() {
	config.LoadTestEnv()
	config.LoadEnv()

	loggerInfoCounter = 0
	loggerErrorCounter = 0
	cryptoServiceGetCryptoCurrentQuoteCounter = 0
	clientPersistenceGetClientsCounter = 0
	eventServiceSendCounter = 0
	cryptoServiceGetCryptoCurrentQuoteError = nil
	clientPersistenceGetClientsError = nil
	eventServiceSendError = nil

	operationUseCase = usecase.OperationUseCase(logger, cryptoWebService, clientPersistence, eventService)

	clients = []model.Client{{
		Id: uuid.NewString(),
	}, {
		Id: uuid.NewString(),
	}}
}

func TestTriggerOperationsStrongBuySuccess(t *testing.T) {
	setup()

	err := operationUseCase.TriggerOperations(summary.StrongBuy)

	assert.Nil(t, err)
	assert.Equal(t, 2, loggerInfoCounter)
	assert.Equal(t, 0, loggerErrorCounter)
	assert.Equal(t, 1, cryptoServiceGetCryptoCurrentQuoteCounter)
	assert.Equal(t, 1, clientPersistenceGetClientsCounter)
	assert.Equal(t, 2, eventServiceSendCounter)
}

func TestTriggerOperationsBuySuccess(t *testing.T) {
	setup()

	err := operationUseCase.TriggerOperations(summary.Buy)

	assert.Nil(t, err)
	assert.Equal(t, 2, loggerInfoCounter)
	assert.Equal(t, 0, loggerErrorCounter)
	assert.Equal(t, 1, cryptoServiceGetCryptoCurrentQuoteCounter)
	assert.Equal(t, 1, clientPersistenceGetClientsCounter)
	assert.Equal(t, 2, eventServiceSendCounter)
}

func TestTriggerOperationsSellSuccess(t *testing.T) {
	setup()

	err := operationUseCase.TriggerOperations(summary.Sell)

	assert.Nil(t, err)
	assert.Equal(t, 2, loggerInfoCounter)
	assert.Equal(t, 0, loggerErrorCounter)
	assert.Equal(t, 0, cryptoServiceGetCryptoCurrentQuoteCounter)
	assert.Equal(t, 1, clientPersistenceGetClientsCounter)
	assert.Equal(t, 2, eventServiceSendCounter)
}

func TestTriggerOperationsStrongSellSuccess(t *testing.T) {
	setup()

	err := operationUseCase.TriggerOperations(summary.StrongSell)

	assert.Nil(t, err)
	assert.Equal(t, 2, loggerInfoCounter)
	assert.Equal(t, 0, loggerErrorCounter)
	assert.Equal(t, 0, cryptoServiceGetCryptoCurrentQuoteCounter)
	assert.Equal(t, 1, clientPersistenceGetClientsCounter)
	assert.Equal(t, 2, eventServiceSendCounter)
}

func TestTriggerOperationsNeutralError(t *testing.T) {
	setup()

	err := custom_error.NewBaseError(operationUseCase.TriggerOperations(summary.Neutral))

	assert.Equal(t, "Operation must be of Buy or Sell type", err.Error())
	assert.Equal(t, "Operation must be of Buy or Sell type", err.InternalError())
	assert.Equal(t, "Error while triggering operations", err.Description())
	assert.Equal(t, 1, loggerInfoCounter)
	assert.Equal(t, 1, loggerErrorCounter)
	assert.Equal(t, 0, cryptoServiceGetCryptoCurrentQuoteCounter)
	assert.Equal(t, 0, clientPersistenceGetClientsCounter)
	assert.Equal(t, 0, eventServiceSendCounter)
}

func TestTriggerOperationsBuyCryptoServiceError(t *testing.T) {
	setup()

	cryptoServiceGetCryptoCurrentQuoteError = errors.New(uuid.NewString())

	err := custom_error.NewBaseError(operationUseCase.TriggerOperations(summary.Buy))

	assert.Equal(t, cryptoServiceGetCryptoCurrentQuoteError.Error(), err.Error())
	assert.Equal(t, "Error while trying to get crypto current quote", err.InternalError())
	assert.Equal(t, "Error while triggering operations", err.Description())
	assert.Equal(t, 1, loggerInfoCounter)
	assert.Equal(t, 1, loggerErrorCounter)
	assert.Equal(t, 1, cryptoServiceGetCryptoCurrentQuoteCounter)
	assert.Equal(t, 0, clientPersistenceGetClientsCounter)
	assert.Equal(t, 0, eventServiceSendCounter)
}

func TestTriggerOperationsClientPersistenceError(t *testing.T) {
	setup()

	clientPersistenceGetClientsError = errors.New(uuid.NewString())

	err := custom_error.NewBaseError(operationUseCase.TriggerOperations(summary.Buy))

	assert.Equal(t, clientPersistenceGetClientsError.Error(), err.Error())
	assert.Equal(t, "Error while trying to find clients", err.InternalError())
	assert.Equal(t, "Error while triggering operations", err.Description())
	assert.Equal(t, 1, loggerInfoCounter)
	assert.Equal(t, 1, loggerErrorCounter)
	assert.Equal(t, 1, cryptoServiceGetCryptoCurrentQuoteCounter)
	assert.Equal(t, 1, clientPersistenceGetClientsCounter)
	assert.Equal(t, 0, eventServiceSendCounter)
}

func TestTriggerOperationsEventServiceError(t *testing.T) {
	setup()

	eventServiceSendError = errors.New(uuid.NewString())

	err := custom_error.NewBaseError(operationUseCase.TriggerOperations(summary.Buy))

	assert.Equal(t, eventServiceSendError.Error(), err.Error())
	assert.Equal(t, "Error while trying to send event", err.InternalError())
	assert.Equal(t, "Error while triggering operations", err.Description())
	assert.Equal(t, 1, loggerInfoCounter)
	assert.Equal(t, 1, loggerErrorCounter)
	assert.Equal(t, 1, cryptoServiceGetCryptoCurrentQuoteCounter)
	assert.Equal(t, 1, clientPersistenceGetClientsCounter)
	assert.Equal(t, 1, eventServiceSendCounter)
}
