package usecase

import (
	"github.com/brienze1/crypto-robot-validator/internal/validator/domain/adapters"
	"github.com/brienze1/crypto-robot-validator/internal/validator/domain/enum/symbol"
	"github.com/brienze1/crypto-robot-validator/internal/validator/domain/exceptions"
	"github.com/brienze1/crypto-robot-validator/internal/validator/domain/model"
	"github.com/brienze1/crypto-robot-validator/pkg/custom_error"
)

type validationUseCase struct {
	lockDB        adapters.LockPersistenceAdapter
	clientDB      adapters.ClientPersistenceAdapter
	clientService adapters.ClientServiceAdapter
	cryptoService adapters.CryptoServiceAdapter
	operationDB   adapters.OperationPersistenceAdapter
	eventService  adapters.EventServiceAdapter
	logger        adapters.LoggerAdapter
}

// ValidationUseCase constructor for class.
func ValidationUseCase(
	lockDB adapters.LockPersistenceAdapter,
	clientDB adapters.ClientPersistenceAdapter,
	clientService adapters.ClientServiceAdapter,
	cryptoService adapters.CryptoServiceAdapter,
	operationDB adapters.OperationPersistenceAdapter,
	eventService adapters.EventServiceAdapter,
	logger adapters.LoggerAdapter,
) *validationUseCase {
	return &validationUseCase{
		lockDB:        lockDB,
		clientDB:      clientDB,
		clientService: clientService,
		cryptoService: cryptoService,
		operationDB:   operationDB,
		eventService:  eventService,
		logger:        logger,
	}
}

// Validate if operation can be executed. client_id key will be locked in cache and locked flag will be set to true on
// client DB during execution of method. After the operation request is validated with client config, an operation is
// created and sent to execution via SNS topic.
func (v *validationUseCase) Validate(operationRequest *model.OperationRequest) error {
	v.logger.Info("Validate start", operationRequest)

	err := v.lockDB.Lock(operationRequest.ClientId)
	if err != nil {
		return v.abort(err, "Error while trying to lock client_id")
	}

	client, err := v.clientDB.GetClient(operationRequest.ClientId)
	if err != nil {
		return v.abort(err, "Error while trying get client from DB")
	}

	err = v.clientDB.Lock(client)
	if err != nil {
		return v.abort(err, "Error while trying to lock client DB")
	}

	balance, err := v.clientService.GetBalance(client.Id)
	if err != nil {
		return v.abort(err, "Error while trying to lock client DB")
	}

	client.SetBalance(balance)

	coin, err := v.cryptoService.GetCrypto(operationRequest.Symbol, symbol.Brl)
	if err != nil {
		return v.abort(err, "Error while trying to get coin from crypto service")
	}

	operation, err := client.CreateOperation(operationRequest, coin)
	if err != nil {
		return v.abort(err, "Error while trying to create operation")
	}

	err = v.operationDB.Save(operation)
	if err != nil {
		return v.abort(err, "Error while trying to save operation")
	}

	err = v.eventService.Send(operation)
	if err != nil {
		return v.abort(err, "Error while trying to send operation event")
	}

	err = v.lockDB.Unlock(client.Id)
	if err != nil {
		return v.abort(err, "Error while trying to unlock client_id")
	}

	err = v.clientDB.Unlock(client)
	if err != nil {
		return v.abort(err, "Error while trying to unlock client DB")
	}

	v.logger.Info("Validate finish", operationRequest, client, operation)
	return nil
}

func (v *validationUseCase) abort(err custom_error.BaseErrorAdapter, message string) error {
	validationError := exceptions.ValidationError(err, message)
	v.logger.Error(validationError, "Validate failed: "+message)

	if err.LockedClientId() != nil {
		ex := v.lockDB.Unlock(*err.LockedClientId())
		if ex != nil {
			panic(ex)
		}
	}

	if err.LockedClient() != nil {
		ex := v.clientDB.Unlock(err.LockedClient())
		if ex != nil {
			panic(ex)
		}
	}

	return validationError
}
