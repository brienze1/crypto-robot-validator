package model

import (
	"github.com/brienze1/crypto-robot-validator/internal/validator/domain/enum/operation_type"
	"github.com/brienze1/crypto-robot-validator/internal/validator/domain/enum/status"
	"github.com/brienze1/crypto-robot-validator/internal/validator/domain/enum/symbol"
	"github.com/google/uuid"
	"time"
)

type Operation struct {
	Id        string
	Status    status.Status
	CreatedAt time.Time
	Locked    bool
	Type      operation_type.OperationType
	Amount    float64
	Base      symbol.Symbol
	Quote     symbol.Symbol
	StopLoss  float64
}

func NewOperation(stopLoss float64) *Operation {
	return &Operation{
		Id:        uuid.NewString(),
		Status:    status.Created,
		CreatedAt: time.Now(),
		Locked:    false,
		StopLoss:  stopLoss,
	}
}
