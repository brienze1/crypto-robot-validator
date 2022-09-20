package adapters

type EventServiceAdapter interface {
	Send(interface{}) error
}
