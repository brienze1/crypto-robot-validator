package adapters

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
)

// HandlerAdapter is an adapter class. Used for handler.Handler implementation.
type HandlerAdapter interface {
	Handle(context context.Context, event events.SQSEvent) error
}
