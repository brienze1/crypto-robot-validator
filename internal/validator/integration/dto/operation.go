package dto

import (
	"github.com/brienze1/crypto-robot-validator/internal/validator/domain/enum/operation_type"
	"github.com/brienze1/crypto-robot-validator/internal/validator/domain/enum/status"
	"github.com/brienze1/crypto-robot-validator/internal/validator/domain/enum/symbol"
	"github.com/brienze1/crypto-robot-validator/internal/validator/domain/model"
	"time"
)

type Operation struct {
	Id        string                       `dynamodbav:"operation_id"`
	Status    status.Status                `dynamodbav:"status"`
	CreatedAt time.Time                    `dynamodbav:"created_at"`
	Locked    bool                         `dynamodbav:"locked"`
	Type      operation_type.OperationType `dynamodbav:"type"`
	Amount    float64                      `dynamodbav:"amount"`
	Base      symbol.Symbol                `dynamodbav:"base"`
	Quote     symbol.Symbol                `dynamodbav:"quote"`
	StopLoss  float64                      `dynamodbav:"stop_loss"`
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
