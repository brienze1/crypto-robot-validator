package exceptions

import "github.com/brienze1/crypto-robot-operation-hub/pkg/custom_error"

func ValidationError(err error, internalError string) custom_error.BaseErrorAdapter {
	return custom_error.NewBaseError(err, internalError, "Error while validating operation")
}
