package eventsourcing

// CommandHandler de handle command
type CommandHandler interface {
	Handle(command Command) error
}