package model

import "github.com/brienze1/crypto-robot-validator/internal/validator/domain/enum/symbol"

type Coin struct {
	Symbol    symbol.Symbol
	Quote     symbol.Symbol
	BuyValue  float64
	SellValue float64
}
