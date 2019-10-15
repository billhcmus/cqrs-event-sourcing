package eventsourcing

import (
	"eventsourcing/proto"
)

// Command interface la abstract cho toan bo command
type Command interface {
	GetType() string
	GetAggregateID() string
	GetAggregateType() string
	GetVersion() uint64
}

// RootCommand is struct root for all command
type RootCommand struct {
	pb.BaseCommand
}

// GetType return command's type
func (root RootCommand) GetType() string {
	return root.Type
}

// GetAggregateID return aggregate's id
func (root RootCommand) GetAggregateID() string {
	return root.AggregateId
}

// GetAggregateType return aggregate's type
func (root RootCommand) GetAggregateType() string {
	return root.AggregateType
}

// GetVersion return aggregate's version was passed by client
func (root RootCommand) GetVersion() uint64 {
	return root.Version
}
