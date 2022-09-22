package adapters

import (
	"github.com/brienze1/crypto-robot-validator/internal/validator/domain/model"
	"github.com/brienze1/crypto-robot-validator/pkg/custom_error"
)

type OperationPersistenceAdapter interface {
	// Save model.Operation in operation repository.
	Save(operation *model.Operation) custom_error.BaseErrorAdapter
}
