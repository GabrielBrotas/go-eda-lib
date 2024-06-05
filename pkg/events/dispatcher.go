package events

import "errors"

var ErrHandlerAlreadyRegistered = errors.New("handler already registered")

type EventDispatcher struct {
	handlers map[string][]EventHandlerInterface
}

func NewEventDispatcher() *EventDispatcher {
	return &EventDispatcher{handlers: make(map[string][]EventHandlerInterface)}
}

// Register the handler for the event
func (ed *EventDispatcher) Register(eventName string, handler EventHandlerInterface) error {
	// Check if the event name is already registered
	if _, ok := ed.handlers[eventName]; ok {
		for _, h := range ed.handlers[eventName] {
			if h == handler {
				return ErrHandlerAlreadyRegistered
			}
		}
	}

	ed.handlers[eventName] = append(ed.handlers[eventName], handler)
	return nil
}

// Dispatch the event
func (ed *EventDispatcher) Dispatch(event EventInterface) error {
	if _, ok := ed.handlers[event.GetName()]; ok {
		for _, handler := range ed.handlers[event.GetName()] {
			// this could be done in a goroutine, figure how to update the tests
			handler.Handle(event)
		}
	}
	return nil
}

// Clear all events
func (ed *EventDispatcher) Clear() error {
	ed.handlers = make(map[string][]EventHandlerInterface)
	return nil
}

// Has checks if the event has a handler
func (ed *EventDispatcher) Has(eventName string, handler EventHandlerInterface) bool {
	if _, ok := ed.handlers[eventName]; ok {
		for _, h := range ed.handlers[eventName] {
			if h == handler {
				return true
			}
		}
	}
	return false
}

// Remove the handler for the event
func (ed *EventDispatcher) Remove(eventName string, handler EventHandlerInterface) error {
	if _, ok := ed.handlers[eventName]; ok {
		for i, h := range ed.handlers[eventName] {
			if h == handler {
				ed.handlers[eventName] = append(ed.handlers[eventName][:i], ed.handlers[eventName][i+1:]...)
				return nil
			}
		}
	}
	return nil
}