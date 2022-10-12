package webservice

import (
	"errors"
	"github.com/brienze1/crypto-robot-validator/internal/validator/application/config"
	"github.com/brienze1/crypto-robot-validator/internal/validator/application/properties"
	"github.com/brienze1/crypto-robot-validator/internal/validator/domain/adapters"
	"github.com/brienze1/crypto-robot-validator/internal/validator/domain/enum/symbol"
	"github.com/brienze1/crypto-robot-validator/internal/validator/integration/webservice"
	"github.com/brienze1/crypto-robot-validator/test/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	biscointWebService adapters.CryptoServiceAdapter
	logger             = mocks.Logger()
	client             = mocks.HttpClient()
	headerBuilder      = mocks.HeaderBuilder()
)

func setup() {
	config.LoadTestEnv()
	properties.Properties().Reload()

	logger.Reset()
	client.Reset()
	client.SetupServer()
	properties.Properties().BiscointUrl = client.GetUrl() + "/"
	properties.Properties().SimulationUrl = client.GetUrl() + "/"
	headerBuilder.Reset()

	biscointWebService = webservice.BiscointWebService(logger, client, headerBuilder)
}

func teardown() {
	client.Close()
}

func TestGetCryptoSuccess(t *testing.T) {
	setup()
	defer teardown()

	client.GetCryptoResponse = `
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

	coin, err := biscointWebService.GetCrypto(symbol.Bitcoin, symbol.Brl)

	assert.Nil(t, err)
	assert.NotNil(t, coin)
	assert.Equal(t, symbol.Bitcoin, coin.Symbol)
	assert.Equal(t, symbol.Brl, coin.Quote)
	assert.Equal(t, 98790.02, coin.BuyValue)
	assert.Equal(t, 97878.96, coin.SellValue)
	assert.Equal(t, 1, client.DoCounter)
	assert.Equal(t, 2, logger.InfoCallCounter)
	assert.Equal(t, 0, logger.ErrorCallCounter)
}

func TestGetCryptoNewRequestFailure(t *testing.T) {
	setup()
	defer teardown()

	properties.Properties().BiscointGetCryptoPath = ""
	properties.Properties().BiscointUrl = string([]byte{0x7f})
	biscointWebService = webservice.BiscointWebService(logger, client, headerBuilder)

	coin, err := biscointWebService.GetCrypto(symbol.Bitcoin, symbol.Brl)

	assert.Equal(t, "parse \"\\x7f\": net/url: invalid control character in URL", err.Error())
	assert.Equal(t, "Error while trying to generate Biscoint get request", err.InternalError())
	assert.Equal(t, "Error while performing Biscoint API request", err.Description())
	assert.Nil(t, coin)
	assert.Equal(t, 0, client.DoCounter)
	assert.Equal(t, 1, logger.InfoCallCounter)
	assert.Equal(t, 1, logger.ErrorCallCounter)
}

func TestGetCryptoDoRequestFailure(t *testing.T) {
	setup()
	defer teardown()

	properties.Properties().BiscointUrl = ""
	properties.Properties().BiscointGetCryptoPath = ""
	biscointWebService = webservice.BiscointWebService(logger, client, headerBuilder)

	coin, err := biscointWebService.GetCrypto(symbol.Bitcoin, symbol.Brl)

	assert.Equal(t, "Get \"?quote=BRL&symbol=BTC\": unsupported protocol scheme \"\"", err.Error())
	assert.Equal(t, "Error while trying to get crypto value from Biscoint", err.InternalError())
	assert.Equal(t, "Error while performing Biscoint API request", err.Description())
	assert.Nil(t, coin)
	assert.Equal(t, 1, client.DoCounter)
	assert.Equal(t, 1, logger.InfoCallCounter)
	assert.Equal(t, 1, logger.ErrorCallCounter)
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
	assert.Equal(t, 1, client.DoCounter)
	assert.Equal(t, 1, logger.InfoCallCounter)
	assert.Equal(t, 1, logger.ErrorCallCounter)
}

func TestGetCryptoDecodeFailure(t *testing.T) {
	setup()
	defer teardown()

	client.GetCryptoResponse = `
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
	assert.Equal(t, 1, client.DoCounter)
	assert.Equal(t, 1, logger.InfoCallCounter)
	assert.Equal(t, 1, logger.ErrorCallCounter)
}

func TestGetBalanceSuccess(t *testing.T) {
	setup()
	defer teardown()

	client.GetBalanceResponse = `
		{
			"message": "",
			"data": {
				"BRL": "9949.75",
    			"BTC": "0.00138164"
			}
		}`

	balance, err := biscointWebService.GetBalance(uuid.NewString(), false)

	assert.Nil(t, err)
	assert.NotNil(t, balance)
	assert.Equal(t, 9949.75, balance.BrlBalance)
	assert.Equal(t, 0.00138164, balance.CryptoBalance)
	assert.Equal(t, 1, client.DoCounter)
	assert.Equal(t, 1, headerBuilder.BiscointHeaderCounter)
	assert.Equal(t, 2, logger.InfoCallCounter)
	assert.Equal(t, 0, logger.ErrorCallCounter)
}

func TestGetBalanceSimulationSuccess(t *testing.T) {
	setup()
	defer teardown()

	client.GetBalanceResponse = `
		{
			"message": "",
			"data": {
				"BRL": "9949.75",
    			"BTC": "0.00138164"
			}
		}`

	balance, err := biscointWebService.GetBalance(uuid.NewString(), true)

	assert.Nil(t, err)
	assert.NotNil(t, balance)
	assert.Equal(t, 9949.75, balance.BrlBalance)
	assert.Equal(t, 0.00138164, balance.CryptoBalance)
	assert.Equal(t, 1, client.DoCounter)
	assert.Equal(t, 1, headerBuilder.BiscointHeaderCounter)
	assert.Equal(t, 2, logger.InfoCallCounter)
	assert.Equal(t, 0, logger.ErrorCallCounter)
}

func TestGetBalanceCreateRequestFailed(t *testing.T) {
	setup()
	defer teardown()

	properties.Properties().BiscointGetBalancePath = ""
	properties.Properties().BiscointUrl = string([]byte{0x7f})
	biscointWebService = webservice.BiscointWebService(logger, client, headerBuilder)

	balance, err := biscointWebService.GetBalance(uuid.NewString(), false)

	assert.Equal(t, "parse \"\\x7f\": net/url: invalid control character in URL", err.Error())
	assert.Equal(t, "Error while trying to generate Biscoint get request", err.InternalError())
	assert.Equal(t, "Error while performing Biscoint API request", err.Description())
	assert.Nil(t, balance)
	assert.Equal(t, 0, client.DoCounter)
	assert.Equal(t, 0, headerBuilder.BiscointHeaderCounter)
	assert.Equal(t, 1, logger.InfoCallCounter)
	assert.Equal(t, 1, logger.ErrorCallCounter)
}

func TestGetBalanceBuildHeaderFailed(t *testing.T) {
	setup()
	defer teardown()

	headerBuilder.BiscointHeaderError = errors.New("error building header")

	balance, err := biscointWebService.GetBalance(uuid.NewString(), false)

	assert.Equal(t, "error building header", err.Error())
	assert.Equal(t, "header builder error", err.InternalError())
	assert.Equal(t, "Error while building header", err.Description())
	assert.Nil(t, balance)
	assert.Equal(t, 0, client.DoCounter)
	assert.Equal(t, 1, headerBuilder.BiscointHeaderCounter)
	assert.Equal(t, 1, logger.InfoCallCounter)
	assert.Equal(t, 1, logger.ErrorCallCounter)
}

func TestGetBalanceDoRequestFailed(t *testing.T) {
	setup()
	defer teardown()

	client.DoError = errors.New("do error")

	balance, err := biscointWebService.GetBalance(uuid.NewString(), false)

	assert.Equal(t, "do error", err.Error())
	assert.Equal(t, "Error while trying to get balance from Biscoint", err.InternalError())
	assert.Equal(t, "Error while performing Biscoint API request", err.Description())
	assert.Nil(t, balance)
	assert.Equal(t, 1, client.DoCounter)
	assert.Equal(t, 1, headerBuilder.BiscointHeaderCounter)
	assert.Equal(t, 1, logger.InfoCallCounter)
	assert.Equal(t, 1, logger.ErrorCallCounter)
}

func TestGetBalanceStatusFailed(t *testing.T) {
	setup()
	defer teardown()

	client.StatusCode = 400

	balance, err := biscointWebService.GetBalance(uuid.NewString(), false)

	assert.Equal(t, "Biscoint API status code not Ok: 400 Bad Request", err.Error())
	assert.Equal(t, "Biscoint API status code not Ok: 400 Bad Request", err.InternalError())
	assert.Equal(t, "Error while performing Biscoint API request", err.Description())
	assert.Nil(t, balance)
	assert.Equal(t, 1, client.DoCounter)
	assert.Equal(t, 1, headerBuilder.BiscointHeaderCounter)
	assert.Equal(t, 1, logger.InfoCallCounter)
	assert.Equal(t, 1, logger.ErrorCallCounter)
}

func TestGetBalanceBRLDecodeFailed(t *testing.T) {
	setup()
	defer teardown()

	client.GetBalanceResponse = `
		{
			"message": "",
			"data": {
				"BRL": error,
    			"BTC": "0.00138164"
			}
		}`

	balance, err := biscointWebService.GetBalance(uuid.NewString(), false)

	assert.Equal(t, "invalid character 'e' looking for beginning of value", err.Error())
	assert.Equal(t, "Error while trying to decode Biscoint balanceResponse API response", err.InternalError())
	assert.Equal(t, "Error while performing Biscoint API request", err.Description())
	assert.Nil(t, balance)
	assert.Equal(t, 1, client.DoCounter)
	assert.Equal(t, 1, headerBuilder.BiscointHeaderCounter)
	assert.Equal(t, 1, logger.InfoCallCounter)
	assert.Equal(t, 1, logger.ErrorCallCounter)
}

func TestGetBalanceToModelBRLFailed(t *testing.T) {
	setup()
	defer teardown()

	client.GetBalanceResponse = `
		{
			"message": "",
			"data": {
				"BRL": "error",
    			"BTC": "0.00138164"
			}
		}`

	balance, err := biscointWebService.GetBalance(uuid.NewString(), false)

	assert.Equal(t, "strconv.ParseFloat: parsing \"error\": invalid syntax", err.Error())
	assert.Equal(t, "Could not convert Biscoint Get Balance response to model", err.InternalError())
	assert.Equal(t, "Error while performing Biscoint API request", err.Description())
	assert.Nil(t, balance)
	assert.Equal(t, 1, client.DoCounter)
	assert.Equal(t, 1, headerBuilder.BiscointHeaderCounter)
	assert.Equal(t, 1, logger.InfoCallCounter)
	assert.Equal(t, 1, logger.ErrorCallCounter)
}

func TestGetBalanceToModelBTCFailed(t *testing.T) {
	setup()
	defer teardown()

	client.GetBalanceResponse = `
		{
			"message": "",
			"data": {
				"BRL": "1000.00",
    			"BTC": "error"
			}
		}`

	balance, err := biscointWebService.GetBalance(uuid.NewString(), false)

	assert.Equal(t, "strconv.ParseFloat: parsing \"error\": invalid syntax", err.Error())
	assert.Equal(t, "Could not convert Biscoint Get Balance response to model", err.InternalError())
	assert.Equal(t, "Error while performing Biscoint API request", err.Description())
	assert.Nil(t, balance)
	assert.Equal(t, 1, client.DoCounter)
	assert.Equal(t, 1, headerBuilder.BiscointHeaderCounter)
	assert.Equal(t, 1, logger.InfoCallCounter)
	assert.Equal(t, 1, logger.ErrorCallCounter)
}
