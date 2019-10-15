package commandhandler

import (
	"eventsourcing"
	"eventsourcing/repository"
	"log"
)

// StCommandHandler chua cac thong tin can thiet de quan ly command
type StCommandHandler struct {
	repository     *repository.Repository
	aggregate      eventsourcing.Aggregate
	bucket, subset string
}

// CreateNewCommandHandler like constructor of the handler
func CreateNewCommandHandler(repo *repository.Repository, aggregate eventsourcing.Aggregate, bucket, subset string) *StCommandHandler {
	return &StCommandHandler{
		repository: repo,
		aggregate:  aggregate,
		bucket:     bucket,
		subset:     subset,
	}
}

// Handle is definition of CommandHandler interface
func (handler *StCommandHandler) Handle(command eventsourcing.Command) error {
	var err error
	version := command.GetVersion()
	// if version != 0 {
	// 	if err = handler.repository.Replay(handler.aggregate, command.GetAggregateID()); err != nil {
	// 		return err
	// 	}
	// 	log.Println(handler.aggregate.GetID())
	// }

	// Do command
	log.Printf("HandleCommand call to handle command")
	if err = handler.aggregate.HandleCommand(command); err != nil {
		return err
	}

	log.Println("ID: ", handler.aggregate.GetID())

	// Save
	log.Printf("Save event and aggregate to database")
	if err = handler.repository.Save(handler.aggregate, version); err != nil {
		log.Fatal(err)
		return err
	}

	// Publish event
	log.Printf("Publish to event bus")
	if err = handler.repository.PublishEvents(handler.aggregate, handler.bucket, handler.subset); err != nil {
		log.Fatal("Error publish to event bus")
		return err
	}

	handler.aggregate.ClearUncommited()

	return nil
}
