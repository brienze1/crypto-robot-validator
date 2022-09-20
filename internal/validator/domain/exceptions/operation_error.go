package exceptions

import "github.com/brienze1/crypto-robot-operation-hub/pkg/custom_error"

func OperationError(err error, internalError string) *custom_error.BaseError {
	return custom_error.NewBaseError(err, internalError, "Error while triggering operations")
}
