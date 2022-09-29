package adapters

import (
	"github.com/brienze1/crypto-robot-validator/internal/validator/domain/enum/symbol"
	"github.com/brienze1/crypto-robot-validator/internal/validator/domain/model"
	"github.com/brienze1/crypto-robot-validator/pkg/custom_error"
)

type CryptoServiceAdapter interface {
	// GetCrypto finds and return a model.Coin object containing values to buy and sell a crypto coin based on symbol
	// and quote (symbol.Symbol).
	GetCrypto(symbol symbol.Symbol, quote symbol.Symbol) (*model.Coin, custom_error.BaseErrorAdapter)

	// GetBalance will search for client balance on external service. ClientId is used to get the apiKey in credentials
	// DB.
	GetBalance(clientId string, useSimulation bool) (*model.Balance, custom_error.BaseErrorAdapter)
}
