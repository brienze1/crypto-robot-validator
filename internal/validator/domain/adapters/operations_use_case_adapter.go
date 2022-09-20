package adapters

import (
	"github.com/brienze1/crypto-robot-operation-hub/internal/operation-hub/domain/enum/summary"
)

type OperationUseCaseAdapter interface {
	TriggerOperations(summary summary.Summary) error
}
