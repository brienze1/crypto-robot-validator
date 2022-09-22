package exceptions

import "github.com/brienze1/crypto-robot-validator/pkg/custom_error"

// RedisPersistenceLockError is the base error class for persistence.RedisPersistence Lock method.
func RedisPersistenceLockError(err error, internalError string) custom_error.BaseErrorAdapter {
	baseError := custom_error.NewBaseError(err, internalError, "Error while using cache to lock id.")
	baseError.SetLocks(false, false)
	return baseError
}

// RedisPersistenceUnlockError is the base error class for persistence.RedisPersistence Unlock method.
func RedisPersistenceUnlockError(err error, internalError string) custom_error.BaseErrorAdapter {
	baseError := custom_error.NewBaseError(err, internalError, "Error while using cache to unlock id.")
	baseError.SetLocks(true, false)
	return baseError
}
