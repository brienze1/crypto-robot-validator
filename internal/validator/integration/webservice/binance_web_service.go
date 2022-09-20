package webservice

import (
	"encoding/json"
	"github.com/brienze1/crypto-robot-operation-hub/internal/operation-hub/application/properties"
	"github.com/brienze1/crypto-robot-operation-hub/internal/operation-hub/domain/adapters"
	"github.com/brienze1/crypto-robot-operation-hub/internal/operation-hub/domain/enum/symbol"
	adapters2 "github.com/brienze1/crypto-robot-operation-hub/internal/operation-hub/integration/adapters"
	"github.com/brienze1/crypto-robot-operation-hub/internal/operation-hub/integration/dto"
	"github.com/brienze1/crypto-robot-operation-hub/internal/operation-hub/integration/exceptions"
	"net/http"
	"net/url"
)

type binanceWebService struct {
	logger adapters.LoggerAdapter
	client adapters2.HTTPClientAdapter
}

func BinanceWebService(logger adapters.LoggerAdapter, client adapters2.HTTPClientAdapter) *binanceWebService {
	return &binanceWebService{
		logger: logger,
		client: client,
	}
}

const symbolKey = "symbol"

func (b *binanceWebService) GetCryptoCurrentQuote(symbol symbol.Symbol) (float64, error) {
	b.logger.Info("Get crypto current quote start", symbol)

	request, err := http.NewRequest(http.MethodGet, properties.Properties().BinanceCryptoSymbolPriceTickerUrl, nil)
	if err != nil {
		return 0.0, b.abort(err, "Error while trying to generate binance get request")
	}

	query := url.Values{}
	query.Add(symbolKey, string(symbol))
	request.URL.RawQuery = query.Encode()

	response, err := b.client.Do(request)
	if err != nil {
		return 0.0, b.abort(err, "Error while trying to get crypto value from binance")
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return 0.0, b.abort(err, "Binance API status code not Ok: "+response.Status)
	}

	var ticker dto.Ticker
	if err := json.NewDecoder(response.Body).Decode(&ticker); err != nil {
		return 0.0, b.abort(err, "Error while trying to decode Binance ticker API response")
	}

	price, err := ticker.GetPrice()
	if err != nil {
		return 0.0, b.abort(err, "Price value returned from Binance ticker API is not a float")
	}

	b.logger.Info("Get crypto current quote finish", symbol, price)
	return price, nil
}

func (b *binanceWebService) abort(err error, message string, metadata ...interface{}) error {
	binanceWebServiceError := exceptions.BinanceWebServiceError(err, message)
	b.logger.Error(binanceWebServiceError, "Get crypto current quote failed: "+message, metadata)
	return binanceWebServiceError
}
