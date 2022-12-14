package webservice

import (
	"encoding/json"
	"github.com/brienze1/crypto-robot-validator/internal/validator/application/properties"
	"github.com/brienze1/crypto-robot-validator/internal/validator/domain/adapters"
	"github.com/brienze1/crypto-robot-validator/internal/validator/domain/enum/symbol"
	"github.com/brienze1/crypto-robot-validator/internal/validator/domain/model"
	adapters2 "github.com/brienze1/crypto-robot-validator/internal/validator/integration/adapters"
	"github.com/brienze1/crypto-robot-validator/internal/validator/integration/dto"
	"github.com/brienze1/crypto-robot-validator/internal/validator/integration/exceptions"
	"github.com/brienze1/crypto-robot-validator/pkg/custom_error"
	"github.com/brienze1/crypto-robot-validator/pkg/log"
	"io"
	"net/http"
	"net/url"
)

type biscointWebService struct {
	logger                 adapters.LoggerAdapter
	client                 adapters2.HTTPClientAdapter
	headerBuilder          adapters2.HeaderBuilderAdapter
	biscointUrl            string
	simulationUrl          string
	biscointGetCryptoPath  string
	biscointGetBalancePath string
}

// BiscointWebService class constructor.
func BiscointWebService(logger adapters.LoggerAdapter, client adapters2.HTTPClientAdapter, headerBuilder adapters2.HeaderBuilderAdapter) *biscointWebService {
	return &biscointWebService{
		logger:                 logger,
		client:                 client,
		headerBuilder:          headerBuilder,
		biscointUrl:            properties.Properties().BiscointUrl,
		simulationUrl:          properties.Properties().SimulationUrl,
		biscointGetCryptoPath:  properties.Properties().BiscointGetCryptoPath,
		biscointGetBalancePath: properties.Properties().BiscointGetBalancePath,
	}
}

const symbolKey = "symbol"
const quoteKey = "quote"

// TODO missing header generation

// GetCrypto finds and return a model.Coin object containing values to buy and sell a crypto coin based on symbol
// and quote (symbol.Symbol).
func (b *biscointWebService) GetCrypto(symbol symbol.Symbol, quote symbol.Symbol) (*model.Coin, custom_error.BaseErrorAdapter) {
	b.logger.Info("Get crypto start", symbol, quoteKey)

	request, err := http.NewRequest(http.MethodGet, b.biscointUrl+b.biscointGetCryptoPath, nil)
	if err != nil {
		return nil, b.abort(err, "Error while trying to generate Biscoint get request")
	}

	log.Logger().Info("request", request)

	query := url.Values{}
	query.Add(symbolKey, symbol.Name())
	query.Add(quoteKey, quote.Name())
	request.URL.RawQuery = query.Encode()

	response, err := b.client.Do(request)
	if err != nil {
		return nil, b.abort(err, "Error while trying to get crypto value from Biscoint")
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(response.Body)

	if response.StatusCode != http.StatusOK {
		return nil, b.abort(err, "Biscoint API status code not Ok: "+response.Status)
	}

	var coinResponse dto.CoinResponse
	if err := json.NewDecoder(response.Body).Decode(&coinResponse); err != nil {
		return nil, b.abort(err, "Error while trying to decode Biscoint coinResponse API response")
	}

	coin := coinResponse.Coin.ToModel()

	b.logger.Info("Get crypto finish", symbol, quote, coin)
	return coin, nil
}

// GetBalance will search for client balance on external service. ClientId is used to get the apiKey in credentials DB.
func (b *biscointWebService) GetBalance(clientId string, useSimulation bool) (*model.Balance, custom_error.BaseErrorAdapter) {
	b.logger.Info("Get balance start", clientId, quoteKey)

	biscointUrl := b.biscointUrl
	if useSimulation {
		biscointUrl = b.simulationUrl
	}
	request, err := http.NewRequest(http.MethodPost, biscointUrl+b.biscointGetBalancePath, nil)
	if err != nil {
		return nil, b.abort(err, "Error while trying to generate Biscoint get request")
	}

	request.Header, err = b.headerBuilder.BiscointHeader(clientId, b.biscointGetBalancePath, `{}`)
	if err != nil {
		return nil, b.abort(err, "Error while trying to generate Biscoint header")
	}

	response, err := b.client.Do(request)
	if err != nil {
		return nil, b.abort(err, "Error while trying to get balance from Biscoint")
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(response.Body)

	if response.StatusCode != http.StatusOK {
		return nil, b.abort(err, "Biscoint API status code not Ok: "+response.Status, response.Body)
	}

	var balanceResponse dto.BalanceResponse
	if err := json.NewDecoder(response.Body).Decode(&balanceResponse); err != nil {
		return nil, b.abort(err, "Error while trying to decode Biscoint balanceResponse API response")
	}

	balance, err := balanceResponse.ToModel()
	if err != nil {
		return nil, b.abort(err, "Could not convert Biscoint Get Balance response to model")
	}

	b.logger.Info("Get balance finish", clientId, quoteKey, balance)
	return balance, nil
}

func (b *biscointWebService) abort(err error, message string, metadata ...interface{}) custom_error.BaseErrorAdapter {
	biscointWebServiceError := exceptions.BiscointWebServiceError(err, message)
	b.logger.Error(biscointWebServiceError, "Biscoint API failed: "+message, metadata)
	return biscointWebServiceError
}
