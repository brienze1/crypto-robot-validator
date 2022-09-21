package custom_error

type BaseErrorAdapter interface {
	Error() string
	Description() string
	InternalError() string
	LockedClientId() bool
	LockedClient() bool
	SetLocks(clientIdLock, clientLock bool)
}

type BaseError struct {
	Message            string `json:"error"`
	InternalMessage    string `json:"internal_error"`
	DescriptionMessage string `json:"description"`
	lockedClientId     bool
	lockedClient       bool
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
			lockedClientId:     false,
			lockedClient:       false,
		}
	}

	switch e := err.(type) {
	case BaseErrorAdapter:
		return &BaseError{
			Message:            e.Error(),
			InternalMessage:    e.InternalError(),
			DescriptionMessage: e.Description(),
			lockedClientId:     e.LockedClientId(),
			lockedClient:       e.LockedClient(),
		}
	default:
		return &BaseError{
			Message:            err.Error(),
			InternalMessage:    internalMessage,
			DescriptionMessage: description,
			lockedClientId:     false,
			lockedClient:       false,
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

func (b *BaseError) LockedClientId() bool {
	return b.lockedClientId
}

func (b *BaseError) LockedClient() bool {
	return b.lockedClient
}

func (b *BaseError) SetLocks(clientIdLock, clientLock bool) {
	b.lockedClientId = clientIdLock
	b.lockedClient = clientLock
}
