package events

import "time"

// Carries the event data
type EventInterface interface {
	GetName() string
	GetDateTime() time.Time
	GetPayload() interface{}
}

// Handles the events
type EventHandlerInterface interface {
	Handle(event EventInterface)
}

type EventDispatcherInterface interface {
	Register(eventName string, handler EventHandlerInterface) error // Register the handler for the event
	Dispatch(event EventInterface) error
	Remove(eventName string, handler EventHandlerInterface) error
	Has(eventName string, handler EventHandlerInterface) bool // Check if the event has a handler
	Clear() error                                             // Clear all events
}
