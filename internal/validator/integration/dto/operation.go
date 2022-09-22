package dto

import (
	"github.com/brienze1/crypto-robot-validator/internal/validator/domain/enum/operation_type"
	"github.com/brienze1/crypto-robot-validator/internal/validator/domain/enum/status"
	"github.com/brienze1/crypto-robot-validator/internal/validator/domain/enum/symbol"
	"github.com/brienze1/crypto-robot-validator/internal/validator/domain/model"
	"time"
)

type Operation struct {
	Id        string                       `json:"operation_id"`
	Status    status.Status                `json:"status"`
	CreatedAt time.Time                    `json:"created_at"`
	Locked    bool                         `json:"locked"`
	Type      operation_type.OperationType `json:"type"`
	Amount    float64                      `json:"amount"`
	Base      symbol.Symbol                `json:"base"`
	Quote     symbol.Symbol                `json:"quote"`
	StopLoss  float64                      `json:"stop_loss"`
}

func OperationDto(operation *model.Operation) *Operation {
	return &Operation{
		Id:        operation.Id,
		Status:    operation.Status,
		CreatedAt: operation.CreatedAt,
		Locked:    operation.Locked,
		Type:      operation.Type,
		Amount:    operation.Amount,
		Base:      operation.Base,
		Quote:     operation.Quote,
		StopLoss:  operation.StopLoss,
	}
}
