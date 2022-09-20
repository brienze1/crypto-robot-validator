package adapters

type LoggerAdapter interface {
	SetCorrelationID(correlationId string)
	Info(message string, metadata ...interface{})
	Error(err error, message string, metadata ...interface{})
}
