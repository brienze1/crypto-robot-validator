package custom_error

type BaseErrorAdapter interface {
	Error() string
	Description() string
	InternalError() string
}

type BaseError struct {
	Message            string `json:"error"`
	InternalMessage    string `json:"internal_error"`
	DescriptionMessage string `json:"description"`
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
		}
	}

	switch e := err.(type) {
	case BaseErrorAdapter:
		return &BaseError{
			Message:            e.Error(),
			InternalMessage:    e.InternalError(),
			DescriptionMessage: e.Description(),
		}
	default:
		return &BaseError{
			Message:            err.Error(),
			InternalMessage:    internalMessage,
			DescriptionMessage: description,
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
