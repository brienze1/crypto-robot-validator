package webservice

import (
	"github.com/brienze1/crypto-robot-validator/internal/validator/application/config"
	"github.com/brienze1/crypto-robot-validator/internal/validator/application/properties"
	"github.com/brienze1/crypto-robot-validator/internal/validator/domain/adapters"
	"github.com/brienze1/crypto-robot-validator/internal/validator/domain/enum/symbol"
	"github.com/brienze1/crypto-robot-validator/internal/validator/integration/webservice"
	"github.com/brienze1/crypto-robot-validator/test/mocks"
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	biscointWebService adapters.CryptoServiceAdapter
	logger             = mocks.Logger()
	client             = mocks.HttpClient()
)

func setup() {
	config.LoadTestEnv()

	logger.Reset()
	client.Reset()
	client.SetupServer()
	client.ServerResponse = `
		{
			"message": "",
			"data": {
				"base": "BTC",
				"quote": "BRL",
				"vol": 11.16484985,
				"low": 94633.34,
				"high": 102068.59,
				"last": 98810.32,
				"ask": 98790.02,
				"askQuoteAmountRef": 1000,
				"askBaseAmountRef": 0.01012248,
				"bid": 97878.96,
				"bidQuoteAmountRef": 1000,
				"bidBaseAmountRef": 0.0102167,
				"timestamp": "2022-09-22T18:07:27.415Z"
			}
		}`
	properties.Properties().BiscointGetCryptoUrl = client.GetUrl()

	biscointWebService = webservice.BiscointWebService(logger, client)
}

func teardown() {
	client.Close()
}

func TestGetCryptoSuccess(t *testing.T) {
	setup()
	defer teardown()

	coin, err := biscointWebService.GetCrypto(symbol.Bitcoin, symbol.Brl)
	assert.Nil(t, err)
	assert.NotNil(t, coin)
	assert.Equal(t, symbol.Bitcoin, coin.Symbol)
	assert.Equal(t, symbol.Brl, coin.Quote)
	assert.Equal(t, 98790.02, coin.BuyValue)
	assert.Equal(t, 97878.96, coin.SellValue)
}

func TestGetCryptoNewRequestFailure(t *testing.T) {
	setup()
	defer teardown()

	properties.Properties().BiscointGetCryptoUrl = string([]byte{0x7f})

	coin, err := biscointWebService.GetCrypto(symbol.Bitcoin, symbol.Brl)
	assert.Equal(t, "parse \"\\x7f\": net/url: invalid control character in URL", err.Error())
	assert.Equal(t, "Error while trying to generate Biscoint get request", err.InternalError())
	assert.Equal(t, "Error while performing Biscoint API request", err.Description())
	assert.Nil(t, coin)
}

func TestGetCryptoDoRequestFailure(t *testing.T) {
	setup()
	defer teardown()

	properties.Properties().BiscointGetCryptoUrl = ""

	coin, err := biscointWebService.GetCrypto(symbol.Bitcoin, symbol.Brl)
	assert.Equal(t, "Get \"?quote=BRL&symbol=BTC\": unsupported protocol scheme \"\"", err.Error())
	assert.Equal(t, "Error while trying to get crypto value from Biscoint", err.InternalError())
	assert.Equal(t, "Error while performing Biscoint API request", err.Description())
	assert.Nil(t, coin)
}

func TestGetCryptoStatusNotOKFailure(t *testing.T) {
	setup()
	defer teardown()

	client.StatusCode = 400

	coin, err := biscointWebService.GetCrypto(symbol.Bitcoin, symbol.Brl)
	assert.Equal(t, "Biscoint API status code not Ok: 400 Bad Request", err.Error())
	assert.Equal(t, "Biscoint API status code not Ok: 400 Bad Request", err.InternalError())
	assert.Equal(t, "Error while performing Biscoint API request", err.Description())
	assert.Nil(t, coin)
}

func TestGetCryptoDecodeFailure(t *testing.T) {
	setup()
	defer teardown()

	client.ServerResponse = `
		{
			"message": "",
			"data": {
				"base": "BTC",
				"quote": "BRL",
				"vol": 11.16484985,
				"low": 94633.34,
				"high": 102068.59,
				"last": 98810.32,
				"ask": "test",
				"askQuoteAmountRef": 1000,
				"askBaseAmountRef": 0.01012248,
				"bid": "test",
				"bidQuoteAmountRef": 1000,
				"bidBaseAmountRef": 0.0102167,
				"timestamp": "2022-09-22T18:07:27.415Z"
			}
		}`

	coin, err := biscointWebService.GetCrypto(symbol.Bitcoin, symbol.Brl)
	assert.Equal(t, "json: cannot unmarshal string into Go struct field Coin.data.ask of type float64", err.Error())
	assert.Equal(t, "Error while trying to decode Biscoint coinResponse API response", err.InternalError())
	assert.Equal(t, "Error while performing Biscoint API request", err.Description())
	assert.Nil(t, coin)
}
