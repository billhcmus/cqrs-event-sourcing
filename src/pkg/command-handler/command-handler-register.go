package cmdhandler

import (
	"fmt"
	"reflect"
)

// ICommandHandlerRegister contain method to access commandhandler registry
type ICommandHandlerRegister interface {
	Set(command interface{}, handler ICommandHandler)
	Get(command interface{}) (ICommandHandler, error)
}

// CommandHandlerRegister is instance of ICommandRegister
type CommandHandlerRegister struct{}

// CreateCommandHandlerRegister get a Register instance
func CreateCommandHandlerRegister() ICommandHandlerRegister {
	return &CommandHandlerRegister{}
}

var commandHandlerRegistry = map[string]ICommandHandler{}

// Set method register type to registry
func (reg *CommandHandlerRegister) Set(command interface{}, handler ICommandHandler) {
	name := reflect.TypeOf(command).String()
	commandHandlerRegistry[name] = handler
}

// Get method register type to registry
func (reg *CommandHandlerRegister) Get(command interface{}) (ICommandHandler, error) {
	rawType := reflect.TypeOf(command)
	name := rawType.String()
	handler, ok := commandHandlerRegistry[name]
	if !ok {
		return nil, fmt.Errorf("%s doesn't exist in command registry", name)
	}
	return handler, nil
}