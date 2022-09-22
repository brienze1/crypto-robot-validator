package mocks

import (
	"github.com/brienze1/crypto-robot-validator/internal/validator/domain/enum/symbol"
	"github.com/brienze1/crypto-robot-validator/internal/validator/domain/model"
	"github.com/brienze1/crypto-robot-validator/internal/validator/integration/exceptions"
	"github.com/brienze1/crypto-robot-validator/pkg/custom_error"
)

type biscointWebService struct {
	GetCryptoCounter      int
	GetCryptoError        error
	CoinExpectedBuyValue  float64
	CoinExpectedSellValue float64
	GetBalanceCounter     int
	GetBalanceError       error
	ClientBrlBalance      float64
	ClientCryptoBalance   float64
}

func BiscointWebService() *biscointWebService {
	return &biscointWebService{}
}

func (b *biscointWebService) GetCrypto(symbol symbol.Symbol, quote symbol.Symbol) (*model.Coin, custom_error.BaseErrorAdapter) {
	b.GetCryptoCounter++

	if b.GetCryptoError != nil {
		return nil, exceptions.BiscointWebServiceError(b.GetCryptoError, "GetCrypto error")
	}

	return &model.Coin{
		Symbol:    symbol,
		Quote:     quote,
		BuyValue:  b.CoinExpectedBuyValue,
		SellValue: b.CoinExpectedSellValue,
	}, nil
}

func (b *biscointWebService) GetBalance(clientId string) (*model.Balance, custom_error.BaseErrorAdapter) {
	b.GetBalanceCounter++

	if b.GetBalanceError != nil {
		return nil, exceptions.BiscointWebServiceError(b.GetBalanceError, "GetBalance error")
	}

	return &model.Balance{
		BrlBalance:    b.ClientBrlBalance,
		CryptoBalance: b.ClientCryptoBalance,
	}, nil
}

func (b *biscointWebService) Reset() {
	b.GetCryptoCounter = 0
	b.GetCryptoError = nil
	b.CoinExpectedBuyValue = 0
	b.CoinExpectedSellValue = 0
	b.GetBalanceCounter = 0
	b.GetBalanceError = nil
	b.ClientBrlBalance = 0
	b.ClientCryptoBalance = 0
}
