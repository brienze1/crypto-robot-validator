package exceptions

import "github.com/brienze1/crypto-robot-validator/pkg/custom_error"

func HandlerError(err error, internalError string) custom_error.BaseErrorAdapter {
	return custom_error.NewBaseError(err, internalError, "Error occurred while handling the event")
}
