package command

// ICommand is the interface that wraps the basic method for an command
type ICommand interface {
	GetType() string
	GetAggregateID() string
	GetAggregateType() string
	GetVersion() uint64
}

// RootCommand is default implement of command
type RootCommand struct {
	Type string 
	AggregateID string
	AggregateType string
	Version uint64
}


// GetType return command's type
func (root RootCommand) GetType() string {
	return root.Type
}

// GetAggregateID return aggregate's id
func (root RootCommand) GetAggregateID() string {
	return root.AggregateID
}

// GetAggregateType return aggregate's type
func (root RootCommand) GetAggregateType() string {
	return root.AggregateType
}

// GetVersion return command version
func (root RootCommand) GetVersion() uint64 {
	return root.Version
}