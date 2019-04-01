package aggregate

import (
	"github.com/billhcmus/cqrs/pkg/command"
	"github.com/billhcmus/cqrs/pkg/event"
)

// RootAggregate is root of all aggregate
type RootAggregate struct {
	AggregateID string
	AggregateType string
	Version uint64
	Changes []event.Event
}

// IAggregate is interface handle method to change state of aggregate
type IAggregate interface {
	Apply(event.Event)
	ApplyChangeHelper(IAggregate, event.Event, bool)
	UnCommited() []event.Event
	HandleCommand(command.ICommand) error
	ClearUncommited()
	GetID() string
	IncrementVersion()
	GetNextVersion() uint64
}

// Implement Aggregate interface on RootAggregate

// GetID return aggregate's ID
func (root RootAggregate) GetID() string {
	return root.AggregateID
}

// IncrementVersion increase version of aggregate
func (root *RootAggregate) IncrementVersion() {
	root.Version++
}

// GetNextVersion return next version
func (root *RootAggregate) GetNextVersion() uint64 {
	return root.Version + 1
}

// ApplyChangeHelper increase current version of the aggregate and change aggregate
func (root *RootAggregate) ApplyChangeHelper(aggregate IAggregate, event event.Event, commit bool) {

	root.IncrementVersion()
	aggregate.Apply(event)

	if commit {	
		event.Version = root.Version
		root.Changes = append(root.Changes, event)
	}
}

// UnCommited return list event uncommited
func (root RootAggregate) UnCommited() []event.Event {
	return root.Changes
}

// ClearUncommited clear uncommited event
func (root *RootAggregate) ClearUncommited() {
	root.Changes = []event.Event{}
}
