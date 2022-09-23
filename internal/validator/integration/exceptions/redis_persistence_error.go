package exceptions

import "github.com/brienze1/crypto-robot-validator/pkg/custom_error"

// RedisPersistenceLockError is the base error class for persistence.RedisPersistence Lock method.
func RedisPersistenceLockError(err error, internalError string, lock bool) custom_error.BaseErrorAdapter {
	baseError := custom_error.NewBaseError(err, internalError, "Error while using Redis cache")
	baseError.SetLocks(lock, false)
	return baseError
}
