package exceptions

import "github.com/brienze1/crypto-robot-validator/pkg/custom_error"

// DynamoDBCredentialsPersistenceError is the base error class for persistence.DynamoDBCredentialsPersistence.
func DynamoDBCredentialsPersistenceError(err error, internalError string) custom_error.BaseErrorAdapter {
	baseError := custom_error.NewBaseError(err, internalError, "Error while using DynamoDB Credentials table")
	baseError.SetLocks(true, true)
	return baseError
}
