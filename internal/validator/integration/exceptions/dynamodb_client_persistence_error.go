package exceptions

import "github.com/brienze1/crypto-robot-validator/pkg/custom_error"

// DynamoDBClientPersistenceError is the base error class for persistence.DynamoDBClientPersistence.
func DynamoDBClientPersistenceError(err error, internalError string) custom_error.BaseErrorAdapter {
	baseError := custom_error.NewBaseError(err, internalError, "Error while using DynamoDB Client table")
	return baseError
}
