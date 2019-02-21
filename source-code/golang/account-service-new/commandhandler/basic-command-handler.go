package commandhandler

import (
	"log"
	"fmt"
	"eventsource"
)

// Handler chua cac thong tin can thiet de quan ly command
type Handler struct {
	repository *eventsourcing.Repository
	aggregate eventsourcing.Aggregate
	bucket, subset string
}

// CreateNewCommandHandler like constructor of the handler
func CreateNewCommandHandler(repo *eventsourcing.Repository, aggregate eventsourcing.Aggregate, bucket, subset string) *Handler {
	return &Handler {
		repository: repo,
		aggregate: aggregate,
		bucket: bucket,
		subset: subset,
	}
}

// Handle is definition of CommandHandler interface
func (handler *Handler)Handle(command eventsourcing.Command) error {
	var err error
	version := command.GetVersion()
	
	log.Printf("With version: %d", version)
	if version != 0 {
		log.Printf("Try to load aggregate id %s", command.GetAggregateID())
		if err = handler.repository.Load(handler.aggregate, command.GetAggregateID()); err != nil {
			return err
		}
		log.Println(handler.aggregate.GetID())
	}

	// Do command
	log.Printf("HandleCommand call to handle command")
	if err = handler.aggregate.HandleCommand(command); err != nil {
		return err
	}


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

	events := handler.aggregate.UnCommited()
	fmt.Print("len event to be saved:", len(events))

	handler.aggregate.ClearUncommited()

	return nil
}