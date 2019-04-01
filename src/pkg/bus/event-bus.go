package bus

import (
	"github.com/billhcmus/cqrs/pkg/event"
)

// IEventBus is interface handle publish method
type IEventBus interface {
	Publish(event event.Event, topic string) error
}