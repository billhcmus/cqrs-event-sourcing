package transfer_service

import (
	"../transfer-command"
	"log"
)

func TransferMoney(detail transfer_command.TransferDetail) {
	// event-sourcing: create event to event store
	log.Println("Do business logic and write to DB", detail.FromAccount)
}