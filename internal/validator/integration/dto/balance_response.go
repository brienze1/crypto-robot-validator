package dto

import (
	"github.com/brienze1/crypto-robot-validator/internal/validator/domain/model"
	"strconv"
)

// BalanceResponse is the model for Biscoint GetBalance response.
type BalanceResponse struct {
	Message string  `json:"message"`
	Balance Balance `json:"data"`
}

func (b *BalanceResponse) ToModel() (*model.Balance, error) {
	brlBalance, err := strconv.ParseFloat(b.Balance.BRL, 64)
	if err != nil {
		return nil, err
	}
	cryptoBalance, err := strconv.ParseFloat(b.Balance.BTC, 64)
	if err != nil {
		return nil, err
	}

	return &model.Balance{
		BrlBalance:    brlBalance,
		CryptoBalance: cryptoBalance,
	}, err
}
