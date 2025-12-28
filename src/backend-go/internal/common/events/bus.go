package events

import "context"

// Publisher publishes events to the message bus
type Publisher interface {
	Publish(ctx context.Context, event Event) error
	Close() error
}

// Subscriber subscribes to events from the message bus
type Subscriber interface {
	Subscribe(ctx context.Context, handler Handler) error
	Close() error
}

// Handler processes incoming events
type Handler interface {
	Handle(ctx context.Context, event Event) error
	EventTypes() []string
}

// Bus combines publishing and subscribing capabilities
type Bus interface {
	Publisher
	Subscriber
}
