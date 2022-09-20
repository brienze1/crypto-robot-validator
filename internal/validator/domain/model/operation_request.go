package model

import (
	"github.com/brienze1/crypto-robot-validator/internal/validator/domain/enum/operation_type"
	"github.com/brienze1/crypto-robot-validator/internal/validator/domain/enum/symbol"
	"time"
)

type OperationRequest struct {
	ClientId  string
	Operation operation_type.OperationType
	Symbol    symbol.Symbol
	StartTime time.Time
}
