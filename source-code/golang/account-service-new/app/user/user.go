package user

import (
	"time"
	"eventsource/app/events"
	"github.com/google/uuid"
	"eventsource/app/commands"
	"eventsource"
)

// User la aggregate
type User struct {
	eventsourcing.RootAggregate
	Name string
	Balance uint64
}

// Implement lai cac method dinh nghia trong interface Aggregate

// ApplyChange tra ve kieu cua aggregate
func (user *User)ApplyChange(event eventsourcing.Event) {
	 switch e := event.Data.(type) {
	 case events.AccountCreated:
		user.Name = e.Name
		user.ID = event.AggregateID
	 case events.MoneyRecharged:
		user.Balance += e.Amount
	 case events.WithdrawlPerformed:
		user.Balance -= e.Amount
	 }
}

// HandleCommand la bang trong eventstore
func (user *User)HandleCommand(command eventsourcing.Command) error {
	uuidv4, _ := uuid.NewRandom()
	event := eventsourcing.Event {
		ID: uuidv4.String(),
		AggregateID: user.ID,
		AggregateType: "Account",
		Timestamp: time.Now().UTC(),
		Version: command.GetVersion(),
	}

	switch c := command.(type) {
	case commands.CreateAccount:
		event.AggregateID = c.AggregateID
		event.Data = events.AccountCreated{Name: c.Name}
		event.Type = "AccountCreated"
	case commands.RechargeMoney:
		event.Data = events.MoneyRecharged{Amount: c.Amount}
		event.Type = "MoneyRecharged"
	case commands.WithdrawlMoney:
		event.Data = events.WithdrawlPerformed{Amount: c.Amount}
		event.Type = "WithdrawlPerformed"
	}

	user.ApplyChangeHelper(user, event, true)
	return nil
}