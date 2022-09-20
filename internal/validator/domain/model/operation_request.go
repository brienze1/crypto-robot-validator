package model

import (
	"github.com/brienze1/crypto-robot-operation-hub/internal/operation-hub/domain/enum/summary"
	"time"
)

type OperationRequest struct {
	ClientId  string `json:"client_id"`
	Operation string `json:"operation"`
	Symbol    string `json:"symbol"`
	StartTime string `json:"start_time"`
}

func NewOperationRequest(client Client, summary summary.Summary) *OperationRequest {
	return &OperationRequest{
		ClientId:  client.Id,
		Operation: summary.OperationTypeString(),
		Symbol:    summary.Name(),
		StartTime: time.Now().String(),
	}
}
