package dto

import (
	"github.com/brienze1/crypto-robot-validator/internal/validator/domain/enum/operation_type"
	"github.com/brienze1/crypto-robot-validator/internal/validator/domain/enum/symbol"
	"github.com/brienze1/crypto-robot-validator/internal/validator/domain/model"
	"time"
)

type OperationRequestDto struct {
	ClientId      string                       `json:"client_id"`
	OperationTypo operation_type.OperationType `json:"operation"`
	Symbol        symbol.Symbol                `json:"symbol"`
	StartTime     time.Time                    `json:"start_time"`
}

func (o *OperationRequestDto) ToModel() *model.OperationRequest {
	return &model.OperationRequest{
		ClientId:  o.ClientId,
		Operation: o.OperationTypo,
		Symbol:    o.Symbol,
		StartTime: o.StartTime,
	}
}
