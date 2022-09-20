package custom_error

import (
	"errors"
	"github.com/brienze1/crypto-robot-operation-hub/pkg/custom_error"
	"github.com/stretchr/testify/assert"
	"testing"
)

type baseErrorExample struct {
	Message            string
	InternalMessage    string
	DescriptionMessage string
}

func (b baseErrorExample) Error() string {
	return b.Message
}

func (b baseErrorExample) InternalError() string {
	return b.InternalMessage
}

func (b baseErrorExample) Description() string {
	return b.DescriptionMessage
}

var (
	errorTest     error
	baseErrorTest baseErrorExample
)

func setup() {
	errorTest = errors.New("error Message")

	baseErrorTest = baseErrorExample{
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

	baseError := custom_error.NewBaseError(baseErrorTest)

	assert.Equal(t, baseErrorTest.Message, baseError.Message)
	assert.Equal(t, baseErrorTest.InternalMessage, baseError.InternalMessage)
	assert.Equal(t, baseErrorTest.DescriptionMessage, baseError.DescriptionMessage)
}
