package custom_error

import (
	"errors"
	"github.com/brienze1/crypto-robot-validator/pkg/custom_error"
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	errorTest     error
	baseErrorTest custom_error.BaseError
)

func setup() {
	errorTest = errors.New("error Message")

	baseErrorTest = custom_error.BaseError{
		Message:            "error Message",
		InternalMessage:    "error InternalMessage",
		DescriptionMessage: "error DescriptionMessage",
	}
}

func TestBaseExceptionFromErrorSuccess(t *testing.T) {
	setup()

	baseError := custom_error.NewBaseError(errorTest)

	assert.Equal(t, errorTest.Error(), baseError.Message)
	assert.Equal(t, "", baseError.InternalMessage)
	assert.Equal(t, "", baseError.DescriptionMessage)
}

func TestBaseExceptionFromBaseErrorSuccess(t *testing.T) {
	setup()

	baseError := custom_error.NewBaseError(&baseErrorTest)

	assert.Equal(t, baseErrorTest.Message, baseError.Message)
	assert.Equal(t, baseErrorTest.InternalMessage, baseError.InternalMessage)
	assert.Equal(t, baseErrorTest.DescriptionMessage, baseError.DescriptionMessage)
}

func TestBaseExceptionFromNilErrorSuccess(t *testing.T) {
	setup()

	baseError := custom_error.NewBaseError(nil, baseErrorTest.Message, baseErrorTest.DescriptionMessage)

	assert.Equal(t, baseErrorTest.Message, baseError.Message)
	assert.Equal(t, baseErrorTest.Message, baseError.InternalMessage)
	assert.Equal(t, baseErrorTest.DescriptionMessage, baseError.DescriptionMessage)
}
