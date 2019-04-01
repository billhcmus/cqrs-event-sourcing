package event

import (
	"github.com/billhcmus/cqrs/common"
	"fmt"
	"time"
	"reflect"
)

// EventRegistry store event type
var EventRegistry = map[string]reflect.Type{}

// Event is in-memory event store data for every event
type Event struct {
	ID            string
	AggregateID   string
	AggregateType string
	Version       uint64
	Timestamp     time.Time
	Type          string
	Data interface{}
}

// IEventRegister is interface which register event
type IEventRegister interface {
	Set(t interface{})
	Get(name string) (reflect.Type, error)
}

// EventRegister that implement event register
type EventRegister struct {}

// CreateEventRegister get a Register instance
func CreateEventRegister() IEventRegister {
	return &EventRegister{}
}

// Set event to event-registry
func (reg *EventRegister) Set(event interface{}) {
	rawType, name := common.GetTypeName(event)
	EventRegistry[name] = rawType
}

// Get event instance form event-registry
func (reg *EventRegister) Get(name string) (reflect.Type, error) {
	rawType, ok := EventRegistry[name]
	if !ok {
		return nil, fmt.Errorf("%s doesn't exist in event registry", name)
	}
	return rawType, nil
}