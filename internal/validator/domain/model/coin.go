package model

import (
	"github.com/brienze1/crypto-robot-validator/internal/validator/application/properties"
	"github.com/brienze1/crypto-robot-validator/internal/validator/domain/enum/operation_type"
	"github.com/brienze1/crypto-robot-validator/internal/validator/domain/enum/symbol"
	"math"
)

type Coin struct {
	Symbol    symbol.Symbol
	Quote     symbol.Symbol
	BuyValue  float64
	SellValue float64
}

func (c Coin) GetMinOperationValue(operationType operation_type.OperationType) float64 {
	switch operationType {
	case operation_type.Buy:
		return c.BuyValue * properties.Properties().MinimumCryptoBuyOperation
	case operation_type.Sell:
		return properties.Properties().MinimumCryptoSellOperation
	}
	return math.MaxFloat64
}
