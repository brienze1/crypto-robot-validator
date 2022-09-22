package domain

import (
	"errors"
	"github.com/brienze1/crypto-robot-validator/internal/validator/application/config"
	"github.com/brienze1/crypto-robot-validator/internal/validator/domain/adapters"
	"github.com/brienze1/crypto-robot-validator/internal/validator/domain/enum/operation_type"
	"github.com/brienze1/crypto-robot-validator/internal/validator/domain/enum/summary_type"
	"github.com/brienze1/crypto-robot-validator/internal/validator/domain/enum/symbol"
	"github.com/brienze1/crypto-robot-validator/internal/validator/domain/model"
	"github.com/brienze1/crypto-robot-validator/internal/validator/domain/usecase"
	"github.com/brienze1/crypto-robot-validator/pkg/custom_error"
	"github.com/brienze1/crypto-robot-validator/pkg/time_utils"
	"github.com/brienze1/crypto-robot-validator/test/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

var (
	validationUseCase    adapters.ValidationUseCaseAdapter
	lockPersistence      = mocks.RedisPersistence()
	clientPersistence    = mocks.DynamoDBClientPersistence()
	clientService        = mocks.BiscointWebService()
	cryptoService        = mocks.BiscointWebService()
	operationPersistence = mocks.DynamoDBOperationPersistence()
	eventService         = mocks.SnsEventService()
	logger               = mocks.Logger()
)

var (
	operationRequest *model.OperationRequest
	client           *model.Client
)

func setup() {
	config.LoadTestEnv()
	config.LoadEnv()

	lockPersistence.Reset()
	clientPersistence.Reset()
	clientService.Reset()
	cryptoService.Reset()
	operationPersistence.Reset()
	eventService.Reset()
	logger.Reset()

	validationUseCase = usecase.ValidationUseCase(
		lockPersistence,
		clientPersistence,
		clientService,
		cryptoService,
		operationPersistence,
		eventService,
		logger,
	)

	operationRequest = &model.OperationRequest{
		ClientId:  uuid.NewString(),
		Operation: operation_type.Buy,
		Symbol:    symbol.Bitcoin,
		StartTime: time.Now(),
	}

	client = &model.Client{
		Id:                        operationRequest.ClientId,
		Active:                    true,
		LockedUntil:               time.Now(),
		Locked:                    false,
		CashAvailable:             10000,
		CashAmount:                10000,
		CashReserved:              0,
		CryptoAvailable:           1,
		CryptoAmount:              1,
		CryptoReserved:            0,
		OperationAmountPercentage: 5,
		DayStopLoss:               100,
		MonthStopLoss:             100,
		BuyOn:                     2,
		SellOn:                    2,
		Symbols:                   []string{"BTC"},
		Summary: []model.Summary{
			{
				Type:         summary_type.Day,
				Day:          time.Now().Day(),
				Month:        int(time.Now().Month()),
				Year:         time.Now().Year(),
				AmountSold:   0,
				AmountBought: 0,
				Profit:       50.00,
			},
			{
				Type:         summary_type.Month,
				Day:          time.Now().Day(),
				Month:        int(time.Now().Month()),
				Year:         time.Now().Year(),
				AmountSold:   0,
				AmountBought: 0,
				Profit:       50.00,
			},
		},
	}

	clientPersistence.AddClient(client)

	cryptoService.CoinExpectedBuyValue = 100000.0
	cryptoService.CoinExpectedSellValue = 99000.0

	clientService.ClientCryptoBalance = 1.0
	clientService.ClientBrlBalance = 1000.0
}

func TestValidateBuySuccess(t *testing.T) {
	setup()

	err := validationUseCase.Validate(operationRequest)

	assert.Nil(t, err)
	assert.Equal(t, false, lockPersistence.IsLocked(client.Id))
	assert.Equal(t, false, client.Locked)
	assert.Equal(t, true, client.LockedUntil.Before(time.Now()))
	assert.Equal(t, 1, len(operationPersistence.GetAllOperations()))
	assert.Equal(t, operation_type.Buy, operationPersistence.GetAllOperations()[0].Type)
	assert.Equal(t, symbol.Bitcoin, operationPersistence.GetAllOperations()[0].Quote)
	assert.Equal(t, symbol.Brl, operationPersistence.GetAllOperations()[0].Base)
	assert.Equal(t, client.OperationStopLoss, operationPersistence.GetAllOperations()[0].StopLoss)
	assert.Equal(t, client.CashAvailable*client.OperationAmountPercentage/100, operationPersistence.GetAllOperations()[0].Amount)
	assert.Equal(t, 1, lockPersistence.LockCounter)
	assert.Equal(t, 1, lockPersistence.UnlockCounter)
	assert.Equal(t, 1, clientPersistence.LockCounter)
	assert.Equal(t, 1, clientPersistence.UnlockCounter)
	assert.Equal(t, 1, clientPersistence.GetClientCounter)
	assert.Equal(t, 1, clientService.GetBalanceCounter)
	assert.Equal(t, 1, cryptoService.GetCryptoCounter)
	assert.Equal(t, 1, operationPersistence.SaveCounter)
	assert.Equal(t, 1, eventService.SendCounter)
	assert.Equal(t, 2, logger.InfoCallCounter)
	assert.Equal(t, 0, logger.ErrorCallCounter)
}

func TestValidateBuyLessThanExpectedOperationCashAmountSuccess(t *testing.T) {
	setup()

	clientService.ClientBrlBalance = 499.0

	err := validationUseCase.Validate(operationRequest)

	assert.Nil(t, err)
	assert.Equal(t, false, lockPersistence.IsLocked(client.Id))
	assert.Equal(t, false, client.Locked)
	assert.Equal(t, true, client.LockedUntil.Before(time.Now()))
	assert.Equal(t, 1, len(operationPersistence.GetAllOperations()))
	assert.Equal(t, operation_type.Buy, operationPersistence.GetAllOperations()[0].Type)
	assert.Equal(t, symbol.Bitcoin, operationPersistence.GetAllOperations()[0].Quote)
	assert.Equal(t, symbol.Brl, operationPersistence.GetAllOperations()[0].Base)
	assert.Equal(t, client.OperationStopLoss, operationPersistence.GetAllOperations()[0].StopLoss)
	assert.Equal(t, clientService.ClientBrlBalance, operationPersistence.GetAllOperations()[0].Amount)
	assert.Equal(t, 1, lockPersistence.LockCounter)
	assert.Equal(t, 1, lockPersistence.UnlockCounter)
	assert.Equal(t, 1, clientPersistence.LockCounter)
	assert.Equal(t, 1, clientPersistence.UnlockCounter)
	assert.Equal(t, 1, clientPersistence.GetClientCounter)
	assert.Equal(t, 1, clientService.GetBalanceCounter)
	assert.Equal(t, 1, cryptoService.GetCryptoCounter)
	assert.Equal(t, 1, operationPersistence.SaveCounter)
	assert.Equal(t, 1, eventService.SendCounter)
	assert.Equal(t, 2, logger.InfoCallCounter)
	assert.Equal(t, 0, logger.ErrorCallCounter)
}

func TestValidateSellSuccess(t *testing.T) {
	setup()

	operationRequest.Operation = operation_type.Sell

	err := validationUseCase.Validate(operationRequest)

	assert.Nil(t, err)
	assert.Equal(t, false, lockPersistence.IsLocked(client.Id))
	assert.Equal(t, false, client.Locked)
	assert.Equal(t, true, client.LockedUntil.Before(time.Now()))
	assert.Equal(t, 1, len(operationPersistence.GetAllOperations()))
	assert.Equal(t, operation_type.Sell, operationPersistence.GetAllOperations()[0].Type)
	assert.Equal(t, symbol.Brl, operationPersistence.GetAllOperations()[0].Quote)
	assert.Equal(t, symbol.Bitcoin, operationPersistence.GetAllOperations()[0].Base)
	assert.Equal(t, client.OperationStopLoss, operationPersistence.GetAllOperations()[0].StopLoss)
	assert.Equal(t, client.CryptoAvailable*client.OperationAmountPercentage/100, operationPersistence.GetAllOperations()[0].Amount)
	assert.Equal(t, 1, lockPersistence.LockCounter)
	assert.Equal(t, 1, lockPersistence.UnlockCounter)
	assert.Equal(t, 1, clientPersistence.LockCounter)
	assert.Equal(t, 1, clientPersistence.UnlockCounter)
	assert.Equal(t, 1, clientPersistence.GetClientCounter)
	assert.Equal(t, 1, clientService.GetBalanceCounter)
	assert.Equal(t, 1, cryptoService.GetCryptoCounter)
	assert.Equal(t, 1, operationPersistence.SaveCounter)
	assert.Equal(t, 1, eventService.SendCounter)
	assert.Equal(t, 2, logger.InfoCallCounter)
	assert.Equal(t, 0, logger.ErrorCallCounter)
}

func TestValidateSellLessThanExpectedOperationCashAmountSuccess(t *testing.T) {
	setup()

	clientService.ClientCryptoBalance = 0.00499
	operationRequest.Operation = operation_type.Sell

	err := validationUseCase.Validate(operationRequest)

	assert.Nil(t, err)
	assert.Equal(t, false, lockPersistence.IsLocked(client.Id))
	assert.Equal(t, false, client.Locked)
	assert.Equal(t, true, client.LockedUntil.Before(time.Now()))
	assert.Equal(t, 1, len(operationPersistence.GetAllOperations()))
	assert.Equal(t, operation_type.Sell, operationPersistence.GetAllOperations()[0].Type)
	assert.Equal(t, symbol.Brl, operationPersistence.GetAllOperations()[0].Quote)
	assert.Equal(t, symbol.Bitcoin, operationPersistence.GetAllOperations()[0].Base)
	assert.Equal(t, client.OperationStopLoss, operationPersistence.GetAllOperations()[0].StopLoss)
	assert.Equal(t, clientService.ClientCryptoBalance, operationPersistence.GetAllOperations()[0].Amount)
	assert.Equal(t, 1, lockPersistence.LockCounter)
	assert.Equal(t, 1, lockPersistence.UnlockCounter)
	assert.Equal(t, 1, clientPersistence.LockCounter)
	assert.Equal(t, 1, clientPersistence.UnlockCounter)
	assert.Equal(t, 1, clientPersistence.GetClientCounter)
	assert.Equal(t, 1, clientService.GetBalanceCounter)
	assert.Equal(t, 1, cryptoService.GetCryptoCounter)
	assert.Equal(t, 1, operationPersistence.SaveCounter)
	assert.Equal(t, 1, eventService.SendCounter)
	assert.Equal(t, 2, logger.InfoCallCounter)
	assert.Equal(t, 0, logger.ErrorCallCounter)
}

func TestValidateEventServiceFailure(t *testing.T) {
	setup()

	eventService.SendError = errors.New("send error")

	err := validationUseCase.Validate(operationRequest)

	assert.NotNil(t, err, "Error should not be nil")
	assert.Equal(t, "Send error", err.(custom_error.BaseErrorAdapter).InternalError())
	assert.Equal(t, "Error while publishing SNS event", err.(custom_error.BaseErrorAdapter).Description())
	assert.Equal(t, "send error", err.(custom_error.BaseErrorAdapter).Error())
	assert.Equal(t, false, lockPersistence.IsLocked(client.Id))
	assert.Equal(t, false, client.Locked)
	assert.Equal(t, true, client.LockedUntil.Before(time.Now()))
	assert.Equal(t, 1, len(operationPersistence.GetAllOperations()))
	assert.Equal(t, client.OperationStopLoss, operationPersistence.GetAllOperations()[0].StopLoss)
	assert.Equal(t, client.CashAvailable*client.OperationAmountPercentage/100, operationPersistence.GetAllOperations()[0].Amount)
	assert.Equal(t, 1, lockPersistence.LockCounter)
	assert.Equal(t, 1, lockPersistence.UnlockCounter)
	assert.Equal(t, 1, clientPersistence.LockCounter)
	assert.Equal(t, 1, clientPersistence.UnlockCounter)
	assert.Equal(t, 1, clientPersistence.GetClientCounter)
	assert.Equal(t, 1, clientService.GetBalanceCounter)
	assert.Equal(t, 1, cryptoService.GetCryptoCounter)
	assert.Equal(t, 1, operationPersistence.SaveCounter)
	assert.Equal(t, 1, eventService.SendCounter)
	assert.Equal(t, 1, logger.InfoCallCounter)
	assert.Equal(t, 1, logger.ErrorCallCounter)
}

func TestValidateOperationPersistenceFailure(t *testing.T) {
	setup()

	operationPersistence.SaveError = errors.New("save error")

	err := validationUseCase.Validate(operationRequest)

	assert.NotNil(t, err, "Error should not be nil")
	assert.Equal(t, "save error", err.(custom_error.BaseErrorAdapter).InternalError())
	assert.Equal(t, "Error while using DynamoDB Operation table", err.(custom_error.BaseErrorAdapter).Description())
	assert.Equal(t, "save error", err.(custom_error.BaseErrorAdapter).Error())
	assert.Equal(t, false, lockPersistence.IsLocked(client.Id))
	assert.Equal(t, false, client.Locked)
	assert.Equal(t, true, client.LockedUntil.Before(time.Now()))
	assert.Equal(t, 0, len(operationPersistence.GetAllOperations()))
	assert.Equal(t, 1, lockPersistence.LockCounter)
	assert.Equal(t, 1, lockPersistence.UnlockCounter)
	assert.Equal(t, 1, clientPersistence.LockCounter)
	assert.Equal(t, 1, clientPersistence.UnlockCounter)
	assert.Equal(t, 1, clientPersistence.GetClientCounter)
	assert.Equal(t, 1, clientService.GetBalanceCounter)
	assert.Equal(t, 1, cryptoService.GetCryptoCounter)
	assert.Equal(t, 1, operationPersistence.SaveCounter)
	assert.Equal(t, 0, eventService.SendCounter)
	assert.Equal(t, 1, logger.InfoCallCounter)
	assert.Equal(t, 1, logger.ErrorCallCounter)
}

func TestValidateCreateOperationDayStopLossFailure(t *testing.T) {
	setup()

	client.Summary[0].Profit = -1000.00

	err := validationUseCase.Validate(operationRequest)

	assert.NotNil(t, err, "Error should not be nil")
	assert.Equal(t, "Client day stop loss reached", err.(custom_error.BaseErrorAdapter).InternalError())
	assert.Equal(t, "Error while validating operation", err.(custom_error.BaseErrorAdapter).Description())
	assert.Equal(t, "validation error", err.(custom_error.BaseErrorAdapter).Error())
	assert.Equal(t, false, lockPersistence.IsLocked(client.Id))
	assert.Equal(t, false, client.Locked)
	assert.Equal(t, time_utils.Time().Tomorrow(), client.LockedUntil)
	assert.Equal(t, 0, len(operationPersistence.GetAllOperations()))
	assert.Equal(t, 1, lockPersistence.LockCounter)
	assert.Equal(t, 1, lockPersistence.UnlockCounter)
	assert.Equal(t, 1, clientPersistence.LockCounter)
	assert.Equal(t, 1, clientPersistence.UnlockCounter)
	assert.Equal(t, 1, clientPersistence.GetClientCounter)
	assert.Equal(t, 1, clientService.GetBalanceCounter)
	assert.Equal(t, 1, cryptoService.GetCryptoCounter)
	assert.Equal(t, 0, operationPersistence.SaveCounter)
	assert.Equal(t, 0, eventService.SendCounter)
	assert.Equal(t, 1, logger.InfoCallCounter)
	assert.Equal(t, 1, logger.ErrorCallCounter)
}

func TestValidateCreateOperationMonthStopLossFailure(t *testing.T) {
	setup()

	client.Summary[1].Profit = -1000.00

	err := validationUseCase.Validate(operationRequest)

	assert.NotNil(t, err, "Error should not be nil")
	assert.Equal(t, "Client month stop loss reached", err.(custom_error.BaseErrorAdapter).InternalError())
	assert.Equal(t, "Error while validating operation", err.(custom_error.BaseErrorAdapter).Description())
	assert.Equal(t, "validation error", err.(custom_error.BaseErrorAdapter).Error())
	assert.Equal(t, false, lockPersistence.IsLocked(client.Id))
	assert.Equal(t, false, client.Locked)
	assert.Equal(t, time_utils.Time().NextMonth(), client.LockedUntil)
	assert.Equal(t, 0, len(operationPersistence.GetAllOperations()))
	assert.Equal(t, 1, lockPersistence.LockCounter)
	assert.Equal(t, 1, lockPersistence.UnlockCounter)
	assert.Equal(t, 1, clientPersistence.LockCounter)
	assert.Equal(t, 1, clientPersistence.UnlockCounter)
	assert.Equal(t, 1, clientPersistence.GetClientCounter)
	assert.Equal(t, 1, clientService.GetBalanceCounter)
	assert.Equal(t, 1, cryptoService.GetCryptoCounter)
	assert.Equal(t, 0, operationPersistence.SaveCounter)
	assert.Equal(t, 0, eventService.SendCounter)
	assert.Equal(t, 1, logger.InfoCallCounter)
	assert.Equal(t, 1, logger.ErrorCallCounter)
}

func TestValidateCreateOperationMinCashFailure(t *testing.T) {
	setup()

	clientService.ClientBrlBalance = 10

	err := validationUseCase.Validate(operationRequest)

	assert.NotNil(t, err, "Error should not be nil")
	assert.Equal(t, "Client does not have minimum cash amount", err.(custom_error.BaseErrorAdapter).InternalError())
	assert.Equal(t, "Error while validating operation", err.(custom_error.BaseErrorAdapter).Description())
	assert.Equal(t, "validation error", err.(custom_error.BaseErrorAdapter).Error())
	assert.Equal(t, false, lockPersistence.IsLocked(client.Id))
	assert.Equal(t, false, client.Locked)
	assert.Equal(t, true, client.LockedUntil.Before(time.Now()))
	assert.Equal(t, 0, len(operationPersistence.GetAllOperations()))
	assert.Equal(t, 1, lockPersistence.LockCounter)
	assert.Equal(t, 1, lockPersistence.UnlockCounter)
	assert.Equal(t, 1, clientPersistence.LockCounter)
	assert.Equal(t, 1, clientPersistence.UnlockCounter)
	assert.Equal(t, 1, clientPersistence.GetClientCounter)
	assert.Equal(t, 1, clientService.GetBalanceCounter)
	assert.Equal(t, 1, cryptoService.GetCryptoCounter)
	assert.Equal(t, 0, operationPersistence.SaveCounter)
	assert.Equal(t, 0, eventService.SendCounter)
	assert.Equal(t, 1, logger.InfoCallCounter)
	assert.Equal(t, 1, logger.ErrorCallCounter)
}

func TestValidateCreateOperationMinCryptoFailure(t *testing.T) {
	setup()

	clientService.ClientCryptoBalance = 0.00001
	operationRequest.Operation = operation_type.Sell

	err := validationUseCase.Validate(operationRequest)

	assert.NotNil(t, err, "Error should not be nil")
	assert.Equal(t, "Client does not have minimum crypto amount", err.(custom_error.BaseErrorAdapter).InternalError())
	assert.Equal(t, "Error while validating operation", err.(custom_error.BaseErrorAdapter).Description())
	assert.Equal(t, "validation error", err.(custom_error.BaseErrorAdapter).Error())
	assert.Equal(t, false, lockPersistence.IsLocked(client.Id))
	assert.Equal(t, false, client.Locked)
	assert.Equal(t, true, client.LockedUntil.Before(time.Now()))
	assert.Equal(t, 0, len(operationPersistence.GetAllOperations()))
	assert.Equal(t, 1, lockPersistence.LockCounter)
	assert.Equal(t, 1, lockPersistence.UnlockCounter)
	assert.Equal(t, 1, clientPersistence.LockCounter)
	assert.Equal(t, 1, clientPersistence.UnlockCounter)
	assert.Equal(t, 1, clientPersistence.GetClientCounter)
	assert.Equal(t, 1, clientService.GetBalanceCounter)
	assert.Equal(t, 1, cryptoService.GetCryptoCounter)
	assert.Equal(t, 0, operationPersistence.SaveCounter)
	assert.Equal(t, 0, eventService.SendCounter)
	assert.Equal(t, 1, logger.InfoCallCounter)
	assert.Equal(t, 1, logger.ErrorCallCounter)
}

func TestValidateGetCryptoFailure(t *testing.T) {
	setup()

	cryptoService.GetCryptoError = errors.New("get crypto error")

	err := validationUseCase.Validate(operationRequest)

	assert.NotNil(t, err, "Error should not be nil")
	assert.Equal(t, "GetCrypto error", err.(custom_error.BaseErrorAdapter).InternalError())
	assert.Equal(t, "Error while performing Biscoint API request", err.(custom_error.BaseErrorAdapter).Description())
	assert.Equal(t, "get crypto error", err.(custom_error.BaseErrorAdapter).Error())
	assert.Equal(t, false, lockPersistence.IsLocked(client.Id))
	assert.Equal(t, false, client.Locked)
	assert.Equal(t, true, client.LockedUntil.Before(time.Now()))
	assert.Equal(t, 0, len(operationPersistence.GetAllOperations()))
	assert.Equal(t, 1, lockPersistence.LockCounter)
	assert.Equal(t, 1, lockPersistence.UnlockCounter)
	assert.Equal(t, 1, clientPersistence.LockCounter)
	assert.Equal(t, 1, clientPersistence.UnlockCounter)
	assert.Equal(t, 1, clientPersistence.GetClientCounter)
	assert.Equal(t, 1, clientService.GetBalanceCounter)
	assert.Equal(t, 1, cryptoService.GetCryptoCounter)
	assert.Equal(t, 0, operationPersistence.SaveCounter)
	assert.Equal(t, 0, eventService.SendCounter)
	assert.Equal(t, 1, logger.InfoCallCounter)
	assert.Equal(t, 1, logger.ErrorCallCounter)
}

func TestValidateGetBalanceFailure(t *testing.T) {
	setup()

	clientService.GetBalanceError = errors.New("get balance error")

	err := validationUseCase.Validate(operationRequest)

	assert.NotNil(t, err, "Error should not be nil")
	assert.Equal(t, "GetBalance error", err.(custom_error.BaseErrorAdapter).InternalError())
	assert.Equal(t, "Error while performing Biscoint API request", err.(custom_error.BaseErrorAdapter).Description())
	assert.Equal(t, "get balance error", err.(custom_error.BaseErrorAdapter).Error())
	assert.Equal(t, false, lockPersistence.IsLocked(client.Id))
	assert.Equal(t, false, client.Locked)
	assert.Equal(t, true, client.LockedUntil.Before(time.Now()))
	assert.Equal(t, 0, len(operationPersistence.GetAllOperations()))
	assert.Equal(t, 1, lockPersistence.LockCounter)
	assert.Equal(t, 1, lockPersistence.UnlockCounter)
	assert.Equal(t, 1, clientPersistence.LockCounter)
	assert.Equal(t, 1, clientPersistence.UnlockCounter)
	assert.Equal(t, 1, clientPersistence.GetClientCounter)
	assert.Equal(t, 1, clientService.GetBalanceCounter)
	assert.Equal(t, 0, cryptoService.GetCryptoCounter)
	assert.Equal(t, 0, operationPersistence.SaveCounter)
	assert.Equal(t, 0, eventService.SendCounter)
	assert.Equal(t, 1, logger.InfoCallCounter)
	assert.Equal(t, 1, logger.ErrorCallCounter)
}

func TestValidateLockClientFailure(t *testing.T) {
	setup()

	clientPersistence.LockError = errors.New("lock client error")

	err := validationUseCase.Validate(operationRequest)

	assert.NotNil(t, err, "Error should not be nil")
	assert.Equal(t, "Lock error", err.(custom_error.BaseErrorAdapter).InternalError())
	assert.Equal(t, "Error while using DynamoDB Client table", err.(custom_error.BaseErrorAdapter).Description())
	assert.Equal(t, "lock client error", err.(custom_error.BaseErrorAdapter).Error())
	assert.Equal(t, false, lockPersistence.IsLocked(client.Id))
	assert.Equal(t, false, client.Locked)
	assert.Equal(t, true, client.LockedUntil.Before(time.Now()))
	assert.Equal(t, 0, len(operationPersistence.GetAllOperations()))
	assert.Equal(t, 1, lockPersistence.LockCounter)
	assert.Equal(t, 1, lockPersistence.UnlockCounter)
	assert.Equal(t, 1, clientPersistence.LockCounter)
	assert.Equal(t, 0, clientPersistence.UnlockCounter)
	assert.Equal(t, 1, clientPersistence.GetClientCounter)
	assert.Equal(t, 0, clientService.GetBalanceCounter)
	assert.Equal(t, 0, cryptoService.GetCryptoCounter)
	assert.Equal(t, 0, operationPersistence.SaveCounter)
	assert.Equal(t, 0, eventService.SendCounter)
	assert.Equal(t, 1, logger.InfoCallCounter)
	assert.Equal(t, 1, logger.ErrorCallCounter)
}

func TestValidateGetClientFailure(t *testing.T) {
	setup()

	clientPersistence.GetClientError = errors.New("get client error")

	err := validationUseCase.Validate(operationRequest)

	assert.NotNil(t, err, "Error should not be nil")
	assert.Equal(t, "GetClient error", err.(custom_error.BaseErrorAdapter).InternalError())
	assert.Equal(t, "Error while using DynamoDB Client table", err.(custom_error.BaseErrorAdapter).Description())
	assert.Equal(t, "get client error", err.(custom_error.BaseErrorAdapter).Error())
	assert.Equal(t, false, lockPersistence.IsLocked(client.Id))
	assert.Equal(t, false, client.Locked)
	assert.Equal(t, true, client.LockedUntil.Before(time.Now()))
	assert.Equal(t, 0, len(operationPersistence.GetAllOperations()))
	assert.Equal(t, 1, lockPersistence.LockCounter)
	assert.Equal(t, 1, lockPersistence.UnlockCounter)
	assert.Equal(t, 1, clientPersistence.GetClientCounter)
	assert.Equal(t, 0, clientPersistence.LockCounter)
	assert.Equal(t, 0, clientPersistence.UnlockCounter)
	assert.Equal(t, 0, clientService.GetBalanceCounter)
	assert.Equal(t, 0, cryptoService.GetCryptoCounter)
	assert.Equal(t, 0, operationPersistence.SaveCounter)
	assert.Equal(t, 0, eventService.SendCounter)
	assert.Equal(t, 1, logger.InfoCallCounter)
	assert.Equal(t, 1, logger.ErrorCallCounter)
}

func TestValidateLockFailure(t *testing.T) {
	setup()

	lockPersistence.LockError = errors.New("lock error")

	err := validationUseCase.Validate(operationRequest)

	assert.NotNil(t, err, "Error should not be nil")
	assert.Equal(t, "Lock error", err.(custom_error.BaseErrorAdapter).InternalError())
	assert.Equal(t, "Error while using cache to lock id.", err.(custom_error.BaseErrorAdapter).Description())
	assert.Equal(t, "lock error", err.(custom_error.BaseErrorAdapter).Error())
	assert.Equal(t, false, lockPersistence.IsLocked(client.Id))
	assert.Equal(t, false, client.Locked)
	assert.Equal(t, true, client.LockedUntil.Before(time.Now()))
	assert.Equal(t, 0, len(operationPersistence.GetAllOperations()))
	assert.Equal(t, 1, lockPersistence.LockCounter)
	assert.Equal(t, 0, lockPersistence.UnlockCounter)
	assert.Equal(t, 0, clientPersistence.GetClientCounter)
	assert.Equal(t, 0, clientPersistence.LockCounter)
	assert.Equal(t, 0, clientPersistence.UnlockCounter)
	assert.Equal(t, 0, clientService.GetBalanceCounter)
	assert.Equal(t, 0, cryptoService.GetCryptoCounter)
	assert.Equal(t, 0, operationPersistence.SaveCounter)
	assert.Equal(t, 0, eventService.SendCounter)
	assert.Equal(t, 1, logger.InfoCallCounter)
	assert.Equal(t, 1, logger.ErrorCallCounter)
}

func TestValidateUnlockFailure(t *testing.T) {
	setup()

	lockPersistence.UnlockError = errors.New("unlock error")

	panicFunction := func() {
		_ = validationUseCase.Validate(operationRequest)
	}

	assert.Panicsf(t, panicFunction, "Should panic")
	assert.Equal(t, true, lockPersistence.IsLocked(client.Id))
	assert.Equal(t, false, client.Locked)
	assert.Equal(t, true, client.LockedUntil.Before(time.Now()))
	assert.Equal(t, 1, len(operationPersistence.GetAllOperations()))
	assert.Equal(t, operation_type.Buy, operationPersistence.GetAllOperations()[0].Type)
	assert.Equal(t, symbol.Bitcoin, operationPersistence.GetAllOperations()[0].Quote)
	assert.Equal(t, symbol.Brl, operationPersistence.GetAllOperations()[0].Base)
	assert.Equal(t, client.OperationStopLoss, operationPersistence.GetAllOperations()[0].StopLoss)
	assert.Equal(t, client.CashAvailable*client.OperationAmountPercentage/100, operationPersistence.GetAllOperations()[0].Amount)
	assert.Equal(t, 1, lockPersistence.LockCounter)
	assert.Equal(t, 2, lockPersistence.UnlockCounter)
	assert.Equal(t, 1, clientPersistence.LockCounter)
	assert.Equal(t, 1, clientPersistence.UnlockCounter)
	assert.Equal(t, 1, clientPersistence.GetClientCounter)
	assert.Equal(t, 1, clientService.GetBalanceCounter)
	assert.Equal(t, 1, cryptoService.GetCryptoCounter)
	assert.Equal(t, 1, operationPersistence.SaveCounter)
	assert.Equal(t, 1, eventService.SendCounter)
	assert.Equal(t, 1, logger.InfoCallCounter)
	assert.Equal(t, 1, logger.ErrorCallCounter)
}

func TestValidateClientUnlockFailure(t *testing.T) {
	setup()

	clientPersistence.UnlockError = errors.New("unlock error")

	panicFunction := func() {
		_ = validationUseCase.Validate(operationRequest)
	}

	assert.Panicsf(t, panicFunction, "Should panic")
	assert.Equal(t, true, lockPersistence.IsLocked(client.Id))
	assert.Equal(t, true, client.Locked)
	assert.Equal(t, true, client.LockedUntil.Before(time.Now()))
	assert.Equal(t, 1, len(operationPersistence.GetAllOperations()))
	assert.Equal(t, operation_type.Buy, operationPersistence.GetAllOperations()[0].Type)
	assert.Equal(t, symbol.Bitcoin, operationPersistence.GetAllOperations()[0].Quote)
	assert.Equal(t, symbol.Brl, operationPersistence.GetAllOperations()[0].Base)
	assert.Equal(t, client.OperationStopLoss, operationPersistence.GetAllOperations()[0].StopLoss)
	assert.Equal(t, client.CashAvailable*client.OperationAmountPercentage/100, operationPersistence.GetAllOperations()[0].Amount)
	assert.Equal(t, 1, lockPersistence.LockCounter)
	assert.Equal(t, 0, lockPersistence.UnlockCounter)
	assert.Equal(t, 1, clientPersistence.LockCounter)
	assert.Equal(t, 2, clientPersistence.UnlockCounter)
	assert.Equal(t, 1, clientPersistence.GetClientCounter)
	assert.Equal(t, 1, clientService.GetBalanceCounter)
	assert.Equal(t, 1, cryptoService.GetCryptoCounter)
	assert.Equal(t, 1, operationPersistence.SaveCounter)
	assert.Equal(t, 1, eventService.SendCounter)
	assert.Equal(t, 1, logger.InfoCallCounter)
	assert.Equal(t, 1, logger.ErrorCallCounter)
}
