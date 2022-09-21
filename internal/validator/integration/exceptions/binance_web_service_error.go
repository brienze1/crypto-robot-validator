package exceptions

import "github.com/brienze1/crypto-robot-validator/pkg/custom_error"

// BiscointWebServiceError is the base error class for webservice.BiscointWebService.
func BiscointWebServiceError(err error, internalError string) custom_error.BaseErrorAdapter {
	baseError := custom_error.NewBaseError(err, internalError, "Error while performing Biscoint API request")
	baseError.SetLocks(true, true)
	return baseError
}
