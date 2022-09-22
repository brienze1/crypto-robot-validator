package exceptions

import "github.com/brienze1/crypto-robot-validator/pkg/custom_error"

// DynamoDBOperationPersistenceError is the base error class for persistence.DynamoDBOperationPersistence.
func DynamoDBOperationPersistenceError(err error, internalError string) custom_error.BaseErrorAdapter {
	baseError := custom_error.NewBaseError(err, internalError, "Error while using DynamoDB Operation table")
	baseError.SetLocks(true, true)
	return baseError
}
