package account

import (
	"fmt"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/billhcmus/cqrs/pkg/aggregate"
	"github.com/billhcmus/cqrs/pkg/command"
	"github.com/billhcmus/cqrs/pkg/event"
	"github.com/google/uuid"
)

// Account la aggregate
type Account struct {
	aggregate.RootAggregate
	Name    string
	Balance int64
}

// Apply apply change to account
func (acc *Account) Apply(event event.Event) {
	logrus.Infof("[Account Service] applying %v", event.Type)
	switch e := event.Data.(type) {
	case *AccountCreated:
		acc.Name = e.Name
		acc.AggregateID = event.AggregateID
		acc.AggregateType = event.AggregateType
	case *MoneyRecharged:
		acc.Balance += e.Amount
	case *WithdrawPerformed:
		acc.Balance -= e.Amount
	}

}

// HandleCommand handle command operation on account
func (acc *Account) HandleCommand(command command.ICommand) error {
	uuidv4, _ := uuid.NewRandom()
	event := event.Event{
		ID:            uuidv4.String(),
		AggregateID:   command.GetAggregateID(),
		AggregateType: "Account",
		Timestamp:     time.Now().UTC(),
	}
	var err error
	switch c := command.(type) {
		case CreateAccount:
			err = create(c, &event)
		case RechargeMoney:
			err = recharge(c, &event)
		case WithdrawMoney:
			err = withdraw(acc, c, &event)
	}
	if err != nil {
		return err
	}

	acc.ApplyChangeHelper(acc, event, true)
	return nil
}

func create(c CreateAccount, e *event.Event) error {
	e.AggregateID = c.GetAggregateID()
	e.Type = "AccountCreated"
	e.Data = &AccountCreated{Name: c.Name}
	return nil
}

func withdraw(acc *Account, c WithdrawMoney, e *event.Event) error {
	var newBalance = acc.Balance - c.Amount
	logrus.Info("[Account Service] new balance ", newBalance)
	if newBalance < 0 {
		return fmt.Errorf("Withdrawl of '%v' failed as there is only '%v' in account '%v'", c.Amount, acc.Balance, acc.GetID())
	}
	e.Data = &WithdrawPerformed{Amount: c.Amount}
	e.Type = "WithdrawPerformed"
	return nil
}

func recharge(c RechargeMoney, e *event.Event) error {
	e.Data = &MoneyRecharged{Amount: c.Amount}
	e.Type = "MoneyRecharged"
	return nil
}