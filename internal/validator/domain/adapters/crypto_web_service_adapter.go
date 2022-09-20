package adapters

import "github.com/brienze1/crypto-robot-operation-hub/internal/operation-hub/domain/enum/symbol"

type CryptoWebServiceAdapter interface {
	GetCryptoCurrentQuote(symbol symbol.Symbol) (float64, error)
}
