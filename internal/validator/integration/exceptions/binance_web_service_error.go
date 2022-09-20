package exceptions

import "github.com/brienze1/crypto-robot-operation-hub/pkg/custom_error"

func BinanceWebServiceError(err error, internalError string) *custom_error.BaseError {
	return custom_error.NewBaseError(err, internalError, "Error while performing Binance API request")
}
