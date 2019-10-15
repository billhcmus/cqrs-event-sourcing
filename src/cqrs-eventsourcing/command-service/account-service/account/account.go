package account

import (
	"eventsourcing"
	pb "eventsourcing/proto"
	"github.com/google/uuid"
	"log"
	"time"
)

// Account la aggregate
type Account struct {
	eventsourcing.RootAggregate
	Name    string
	Balance uint64
}

// User's command

// CreateAccount is user command
type CreateAccount struct {
	eventsourcing.RootCommand
	Name string
}

// RechargeMoney is user command
type RechargeMoney struct {
	eventsourcing.RootCommand
	Amount uint64
}

// WithdrawMoney is user command
type WithdrawMoney struct {
	eventsourcing.RootCommand
	Amount uint64
}

// User's event

// AccountCreated is event
type AccountCreated struct {
	Name string `json:"name"`
}

// MoneyRecharged is event
type MoneyRecharged struct {
	Amount uint64 `json:"amount"`
}

// WithdrawPerformed is event
type WithdrawPerformed struct {
	Amount uint64 `json:"amount"`
}

// Implement lai cac method dinh nghia trong interface Aggregate

// ApplyChange tra ve kieu cua aggregate
func (acc *Account) ApplyChange(event eventsourcing.Event) {
	switch e := event.Data.(type) {
	case *AccountCreated:
		acc.Name = e.Name
		acc.AggregateId = event.AggregateId
		acc.AggregateType = event.AggregateType
		acc.Version = event.Version
	case *MoneyRecharged:
		acc.Balance += e.Amount
	case *WithdrawPerformed:
		acc.Balance -= e.Amount
	default:
		log.Println("abc ", e)
	}
}

// HandleCommand handle command of user
func (acc *Account) HandleCommand(command eventsourcing.Command) error {
	uuidv4, _ := uuid.NewRandom()
	event := eventsourcing.Event{
		BaseEvent: pb.BaseEvent{
			EventId:       uuidv4.String(),
			AggregateId:   acc.AggregateId,
			AggregateType: "Account",
			Version:       command.GetVersion(),
			Timestamp:     uint64(time.Now().Unix()),
		},
	}

	switch c := command.(type) {
	case CreateAccount:
		event.AggregateId = command.GetAggregateID()
		event.EventType = "AccountCreated"
		event.Data = &AccountCreated{Name: c.Name}
	case RechargeMoney:
		event.Data = &MoneyRecharged{Amount: c.Amount}
		event.EventType = "MoneyRecharged"
	case WithdrawMoney:
		event.Data = &WithdrawPerformed{Amount: c.Amount}
		event.EventType = "WithdrawPerformed"
	}

	acc.ApplyChangeHelper(acc, event, true)
	return nil
}
