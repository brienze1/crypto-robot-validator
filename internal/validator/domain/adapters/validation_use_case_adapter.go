package adapters

import (
	"github.com/brienze1/crypto-robot-validator/internal/validator/domain/model"
)

// ValidationUseCaseAdapter adapter for usecase.ValidationUseCase.
type ValidationUseCaseAdapter interface {
	// Validate if operation can be executed. client_id key will be locked in cache and locked flag will be set to true on
	// client DB during execution of method. After the operation request is validated with client config, an operation is
	// created and sent to execution via SNS topic.
	Validate(operationRequest *model.OperationRequest) error
}
