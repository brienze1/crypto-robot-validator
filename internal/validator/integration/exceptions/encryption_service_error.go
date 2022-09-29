package exceptions

import "github.com/brienze1/crypto-robot-validator/pkg/custom_error"

// EncryptionServiceError is the base error class for utils.EncryptionService.
func EncryptionServiceError(err error, internalError string) custom_error.BaseErrorAdapter {
	baseError := custom_error.NewBaseError(err, internalError, "Error while performing encryption")
	baseError.SetLocks(true, true)
	return baseError
}
