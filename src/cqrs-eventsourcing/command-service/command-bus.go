package eventsourcing

// CommandBus is a interface have CommandHandle method
type CommandBus interface {
	HandleCommand(Command)
}