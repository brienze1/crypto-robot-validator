package exceptions

import "github.com/brienze1/crypto-robot-validator/pkg/custom_error"

// SNSEventServiceError is the base error class for eventservice.SNSEventService.
func SNSEventServiceError(err error, internalError string) custom_error.BaseErrorAdapter {
	baseError := custom_error.NewBaseError(err, internalError, "Error while publishing SNS event")
	baseError.SetLocks(true, true)
	return baseError
}
