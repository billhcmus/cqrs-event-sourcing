package eventstore

import (
	"github.com/billhcmus/cqrs/pkg/event"
)

// IEventStore handle operation effect to events
type IEventStore interface {
	Save(event []event.Event, baseversion uint64) error
	Load(aggregateID string) ([]event.Event, error)
}
