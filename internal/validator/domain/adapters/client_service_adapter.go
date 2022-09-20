package adapters

import (
	"github.com/brienze1/crypto-robot-validator/internal/validator/domain/model"
	"github.com/brienze1/crypto-robot-validator/pkg/custom_error"
)

type ClientServiceAdapter interface {
	// GetBalance will search for client balance on external service. ClientId is used to get the apiKey in credentials
	// DB.
	GetBalance(clientId string) (*model.Balance, custom_error.BaseErrorAdapter)
}
