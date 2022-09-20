package webservice

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/brienze1/crypto-robot-operation-hub/internal/operation-hub/application/config"
	"github.com/brienze1/crypto-robot-operation-hub/internal/operation-hub/application/properties"
	"github.com/brienze1/crypto-robot-operation-hub/internal/operation-hub/domain/adapters"
	"github.com/brienze1/crypto-robot-operation-hub/internal/operation-hub/domain/enum/symbol"
	"github.com/brienze1/crypto-robot-operation-hub/internal/operation-hub/integration/dto"
	"github.com/brienze1/crypto-robot-operation-hub/internal/operation-hub/integration/webservice"
	"github.com/brienze1/crypto-robot-operation-hub/pkg/custom_error"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

type (
	loggerMock struct {
		adapters.LoggerAdapter
	}
	clientMock struct {
	}
)

var (
	loggerInfoCounter  int
	loggerErrorCounter int
	clientDoCounter    int
	binanceAPICounter  int
)

var (
	binanceAPIError  error
	binanceAPIStatus int
)

func (l loggerMock) Info(string, ...interface{}) {
	loggerInfoCounter++
}

func (l loggerMock) Error(error, string, ...interface{}) {
	loggerErrorCounter++
}

func (c clientMock) Do(req *http.Request) (*http.Response, error) {
	clientDoCounter++

	if binanceAPIError != nil {
		return nil, binanceAPIError
	}

	realClient := http.Client{}
	return realClient.Do(req)
}

var (
	binanceWebService adapters.CryptoWebServiceAdapter
	logger            = loggerMock{}
	client            = clientMock{}
)

var (
	ticker        dto.Ticker
	expectedPrice float64
)

func setup() {
	config.LoadTestEnv()
	properties.Properties().BinanceCryptoSymbolPriceTickerUrl = ""

	loggerInfoCounter = 0
	loggerErrorCounter = 0
	binanceAPICounter = 0
	clientDoCounter = 0
	binanceAPIError = nil
	binanceAPIStatus = 200

	binanceWebService = webservice.BinanceWebService(logger, client)

	expectedPrice = 21537.81000000
	ticker = dto.Ticker{
		Symbol: string(symbol.Bitcoin),
		Price:  fmt.Sprintf("%f", expectedPrice),
	}

}

func setupServer(responseBody interface{}) *httptest.Server {
	response, _ := json.Marshal(responseBody)
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		binanceAPICounter++

		if binanceAPIStatus != 200 {
			http.Error(w, "error test", binanceAPIStatus)
		}

		_, _ = w.Write(response)
	}))
}

func teardown(server *httptest.Server) {
	server.Close()
}

func TestBinanceWebServiceSuccess(t *testing.T) {
	setup()
	server := setupServer(ticker)
	properties.Properties().BinanceCryptoSymbolPriceTickerUrl = server.URL

	price, err := binanceWebService.GetCryptoCurrentQuote(symbol.Bitcoin)
	assert.Nil(t, err)
	assert.Equal(t, expectedPrice, price)
	assert.Equal(t, 1, binanceAPICounter)
	assert.Equal(t, 1, clientDoCounter)
	assert.Equal(t, 2, loggerInfoCounter)
	assert.Equal(t, 0, loggerErrorCounter)

	defer teardown(server)
}

func TestBinanceWebServiceNewRequestError(t *testing.T) {
	setup()
	properties.Properties().BinanceCryptoSymbolPriceTickerUrl = string([]byte{0x7f})

	price, err := binanceWebService.GetCryptoCurrentQuote(symbol.Bitcoin)
	assert.Equal(t, "Error while trying to generate binance get request", err.(custom_error.BaseErrorAdapter).InternalError())
	assert.Equal(t, "parse \"\\x7f\": net/url: invalid control character in URL", err.(custom_error.BaseErrorAdapter).Error())
	assert.Equal(t, "Error while performing Binance API request", err.(custom_error.BaseErrorAdapter).Description())
	assert.Equal(t, 0.0, price)
	assert.Equal(t, 0, binanceAPICounter)
	assert.Equal(t, 0, clientDoCounter)
	assert.Equal(t, 1, loggerInfoCounter)
	assert.Equal(t, 1, loggerErrorCounter)
}

func TestBinanceWebServiceClientDoError(t *testing.T) {
	setup()
	binanceAPIError = errors.New(uuid.NewString())

	price, err := binanceWebService.GetCryptoCurrentQuote(symbol.Bitcoin)
	assert.Equal(t, "Error while trying to get crypto value from binance", err.(custom_error.BaseErrorAdapter).InternalError())
	assert.Equal(t, binanceAPIError.Error(), err.(custom_error.BaseErrorAdapter).Error())
	assert.Equal(t, "Error while performing Binance API request", err.(custom_error.BaseErrorAdapter).Description())
	assert.Equal(t, 0.0, price)
	assert.Equal(t, 0, binanceAPICounter)
	assert.Equal(t, 1, clientDoCounter)
	assert.Equal(t, 1, loggerInfoCounter)
	assert.Equal(t, 1, loggerErrorCounter)
}

func TestBinanceWebServiceBinanceAPIStatusError(t *testing.T) {
	setup()
	binanceAPIStatus = http.StatusBadRequest
	server := setupServer(ticker)
	properties.Properties().BinanceCryptoSymbolPriceTickerUrl = server.URL

	price, err := binanceWebService.GetCryptoCurrentQuote(symbol.Bitcoin)
	assert.Equal(t, "Binance API status code not Ok: 400 Bad Request", err.(custom_error.BaseErrorAdapter).InternalError())
	assert.Equal(t, "Binance API status code not Ok: 400 Bad Request", err.(custom_error.BaseErrorAdapter).Error())
	assert.Equal(t, "Error while performing Binance API request", err.(custom_error.BaseErrorAdapter).Description())
	assert.Equal(t, 0.0, price)
	assert.Equal(t, 1, binanceAPICounter)
	assert.Equal(t, 1, clientDoCounter)
	assert.Equal(t, 1, loggerInfoCounter)
	assert.Equal(t, 1, loggerErrorCounter)

	teardown(server)
}

func TestBinanceWebServiceBinanceAPIResponseDecodeError(t *testing.T) {
	setup()
	wrongResponse := map[string]int{
		"Symbol": 10,
	}
	server := setupServer(wrongResponse)
	properties.Properties().BinanceCryptoSymbolPriceTickerUrl = server.URL

	price, err := binanceWebService.GetCryptoCurrentQuote(symbol.Bitcoin)
	assert.Equal(t, "Error while trying to decode Binance ticker API response", err.(custom_error.BaseErrorAdapter).InternalError())
	assert.Equal(t, "json: cannot unmarshal number into Go struct field Ticker.symbol of type string", err.(custom_error.BaseErrorAdapter).Error())
	assert.Equal(t, "Error while performing Binance API request", err.(custom_error.BaseErrorAdapter).Description())
	assert.Equal(t, 0.0, price)
	assert.Equal(t, 1, binanceAPICounter)
	assert.Equal(t, 1, clientDoCounter)
	assert.Equal(t, 1, loggerInfoCounter)
	assert.Equal(t, 1, loggerErrorCounter)

	teardown(server)
}

func TestBinanceWebServiceBinanceAPIResponseTickerNotFloatError(t *testing.T) {
	setup()
	wrongResponse := map[string]string{
		"Symbol": "BTCBUSD",
		"Price":  "string",
	}
	server := setupServer(wrongResponse)
	properties.Properties().BinanceCryptoSymbolPriceTickerUrl = server.URL

	price, err := binanceWebService.GetCryptoCurrentQuote(symbol.Bitcoin)
	assert.Equal(t, "Price value returned from Binance ticker API is not a float", err.(custom_error.BaseErrorAdapter).InternalError())
	assert.Equal(t, "strconv.ParseFloat: parsing \"string\": invalid syntax", err.(custom_error.BaseErrorAdapter).Error())
	assert.Equal(t, "Error while performing Binance API request", err.(custom_error.BaseErrorAdapter).Description())
	assert.Equal(t, 0.0, price)
	assert.Equal(t, 1, binanceAPICounter)
	assert.Equal(t, 1, clientDoCounter)
	assert.Equal(t, 1, loggerInfoCounter)
	assert.Equal(t, 1, loggerErrorCounter)

	teardown(server)
}
