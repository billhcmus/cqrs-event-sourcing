package eventsourcing

// RootAggregate la aggregate goc cac aggregate sau embedded root vao no
type RootAggregate struct {
	ID string
	Type string
	Version uint64
	Changes []Event
}

// Aggregate cac method thuc hien commnad 
type Aggregate interface {
	ApplyChange(Event)
	ApplyChangeHelper(Aggregate, Event, bool)
	UnCommited() []Event
	HandleCommand(Command) error
	ClearUncommited()
	GetID() string
	IncrementVersion()
}


// Implement Aggregate interface on RootAggregate

// GetID return aggregate's ID
func (root RootAggregate) GetID() string {
	return root.ID
}

// IncrementVersion increase version of aggregate
func (root *RootAggregate) IncrementVersion() {
	root.Version++
}

// ApplyChangeHelper tang version hien tai cua aggregate va thay doi aggregate
func (root *RootAggregate) ApplyChangeHelper(aggregate Aggregate, event Event, commit bool) {
	// tang version
	root.IncrementVersion()
	aggregate.ApplyChange(event)
	
	if commit {
		event.Version = root.Version
		root.Changes = append(root.Changes, event)
	}
}

// UnCommited tra ve danh sach cac event duoc save
func (root RootAggregate)UnCommited() []Event{
	return root.Changes
}

// ClearUncommited clear uncommited event
func (root *RootAggregate)ClearUncommited() {
	root.Changes = []Event{}
}