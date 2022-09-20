package custom_error

import "github.com/brienze1/crypto-robot-validator/internal/validator/domain/model"

type BaseErrorAdapter interface {
	Error() string
	Description() string
	InternalError() string
	LockedClientId() *string
	LockedClient() *model.Client
}

type BaseError struct {
	Message            string        `json:"error"`
	InternalMessage    string        `json:"internal_error"`
	DescriptionMessage string        `json:"description"`
	ClientId           *string       `json:"-"`
	Client             *model.Client `json:"-"`
}

func NewBaseError(err error, messages ...string) *BaseError {
	internalMessage := ""
	description := ""
	if len(messages) > 0 {
		internalMessage = messages[0]
	}
	if len(messages) > 1 {
		description = messages[1]
	}

	if err == nil {
		return &BaseError{
			Message:            internalMessage,
			InternalMessage:    internalMessage,
			DescriptionMessage: description,
			ClientId:           nil,
			Client:             nil,
		}
	}

	switch e := err.(type) {
	case BaseErrorAdapter:
		return &BaseError{
			Message:            e.Error(),
			InternalMessage:    e.InternalError(),
			DescriptionMessage: e.Description(),
			ClientId:           e.LockedClientId(),
			Client:             e.LockedClient(),
		}
	default:
		return &BaseError{
			Message:            err.Error(),
			InternalMessage:    internalMessage,
			DescriptionMessage: description,
			ClientId:           nil,
			Client:             nil,
		}
	}
}

func (b *BaseError) Error() string {
	return b.Message
}

func (b *BaseError) InternalError() string {
	return b.InternalMessage
}

func (b *BaseError) Description() string {
	return b.DescriptionMessage
}

func (b *BaseError) LockedClientId() *string {
	return b.ClientId
}

func (b *BaseError) LockedClient() *model.Client {
	return b.Client
}
