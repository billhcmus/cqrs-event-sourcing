package eventsourcing

import (
	"reflect"
)

// Register define method perform on some Registry, ex: event-registry or command-registry
type Register interface {
	Set(t interface{})
	Get(name string) (reflect.Type, error)
}
