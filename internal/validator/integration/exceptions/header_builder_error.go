package exceptions

import "github.com/brienze1/crypto-robot-validator/pkg/custom_error"

// HeaderBuilderError is the base error class for utils.HeaderBuilder.
func HeaderBuilderError(err error, internalError string) custom_error.BaseErrorAdapter {
	baseError := custom_error.NewBaseError(err, internalError, "Error while building header")
	baseError.SetLocks(true, true)
	return baseError
}
