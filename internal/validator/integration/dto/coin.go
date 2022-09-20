package dto

import (
	"github.com/brienze1/crypto-robot-validator/internal/validator/domain/enum/symbol"
	"github.com/brienze1/crypto-robot-validator/internal/validator/domain/model"
)

// Coin is the model for Biscoint GetCoin response.
type Coin struct {
	Symbol    symbol.Symbol `json:"base"`
	Quote     symbol.Symbol `json:"quote"`
	BuyValue  float64       `json:"ask"`
	SellValue float64       `json:"bid"`
}

// ToModel returns model.Coin from dto.Coin.
func (c *Coin) ToModel() *model.Coin {
	return &model.Coin{
		Symbol:    c.Symbol,
		Quote:     c.Quote,
		BuyValue:  c.BuyValue,
		SellValue: c.SellValue,
	}
}
