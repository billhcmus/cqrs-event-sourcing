package eventsourcing

import (
	"eventsourcing/proto"
	"eventsourcing/utils"
	"fmt"
	"reflect"
)

// eventRegistry chua thong tin type cua Event
var eventRegistry = map[string]reflect.Type{}

// Event is in-memory event store data for every event
type Event struct {
	pb.BaseEvent
	Data interface{}
}

// Register define method perform on some Registry, ex: event-registry or command-registry
type Register interface {
	Set(t interface{})
	Get(name string) (reflect.Type, error)
}

// EventRegister that implement method in Register interface
type EventRegister struct {
}

// CreateEventRegister get a Register instance
func CreateEventRegister() Register {
	return &EventRegister{}
}

// Set event to event-registry
func (reg *EventRegister) Set(event interface{}) {
	rawType, name := utils.GetTypeName(event)
	eventRegistry[name] = rawType
}

// Get event instance form event-registry
func (reg *EventRegister) Get(name string) (reflect.Type, error) {
	rawType, ok := eventRegistry[name]
	if !ok {
		return nil, fmt.Errorf("%s doesn't exist in event registry", name)
	}
	return rawType, nil
}
