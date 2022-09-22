package adapters

import (
	"github.com/brienze1/crypto-robot-validator/pkg/custom_error"
)

type EventServiceAdapter interface {
	// Send event containing object to topic.
	Send(object interface{}) custom_error.BaseErrorAdapter
}
