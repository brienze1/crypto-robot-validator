package exceptions

import (
	"errors"
	"github.com/brienze1/crypto-robot-validator/pkg/custom_error"
)

func ValidationError(err error, internalError string) custom_error.BaseErrorAdapter {
	return custom_error.NewBaseError(err, internalError, "Error while validating operation")
}

func NewValidationError(err string) custom_error.BaseErrorAdapter {
	baseError := custom_error.NewBaseError(errors.New("validation error"), err, "Error while validating operation")
	baseError.SetLocks(true, true)
	return baseError
}
