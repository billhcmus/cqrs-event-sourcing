package eventsourcing

import (
	"eventsource/utils"
	"fmt"
	"reflect"
	"time"
)

// eventRegistry chua thong tin type cua Event
var eventRegistry = map[string]reflect.Type{}

// Event is in-memory event store data for every event
type Event struct {
	ID            string
	AggregateID   string
	AggregateType string
	Version       uint64
	Timestamp     time.Time
	Type          string
	Data          interface{}
}

// EventType that implement method in Register interface
type EventType struct {
}

// CreateEventRegister get a Register instance
func CreateEventRegister() Register {
	return &EventType{}
}

// Set event to event-registry
func (eventType *EventType) Set(event interface{}) {
	rawType, name := utils.GetTypeName(event)
	eventRegistry[name] = rawType
}

// Get event instance form event-registry
func (eventType *EventType) Get(name string) (reflect.Type, error) {
	rawType, ok := eventRegistry[name]
	if !ok {
		return nil, fmt.Errorf("%s doesn't exist in event registry", name)
	}
	return rawType, nil
}