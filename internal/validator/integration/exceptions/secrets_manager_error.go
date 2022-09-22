package exceptions

import "github.com/brienze1/crypto-robot-validator/pkg/custom_error"

func SecretsManagerError(err error, internalError string) *custom_error.BaseError {
	return custom_error.NewBaseError(err, internalError, "Error while trying to get secret from secrets manager")
}
