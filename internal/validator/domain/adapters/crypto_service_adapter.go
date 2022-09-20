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
}
