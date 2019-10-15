package eventsourcing

import (
	"fmt"
	"reflect"
)


// CommandHandler handle command
type CommandHandler interface {
	Handle(Command) error
}

// CommandHandlerRegister contain method to access commandhandler registry
type CommandHandlerRegister interface {
	Set(command interface{}, handler CommandHandler)
	Get(command interface{}) (CommandHandler, error)
}

// CommandRegister is instance of Register
type CommandRegister struct{}

// CreateCommandRegister get a Register instance
func CreateCommandRegister() CommandHandlerRegister {
	return &CommandRegister{}
}

var commandRegistry = map[string]CommandHandler{}

// Set method register type to registry
func (reg *CommandRegister) Set(command interface{}, handler CommandHandler) {
	name := reflect.TypeOf(command).String()
	commandRegistry[name] = handler
}

// Get method register type to registry
func (reg *CommandRegister) Get(command interface{}) (CommandHandler, error) {
	rawType := reflect.TypeOf(command)
	name := rawType.String()
	handler, ok := commandRegistry[name]
	if !ok {
		return nil, fmt.Errorf("%s doesn't exist in command registry", name)
	}
	return handler, nil
}
