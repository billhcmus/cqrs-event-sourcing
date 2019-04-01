package account

import (
	"github.com/billhcmus/cqrs/pkg/command"
)

// CreateAccount is user command
type CreateAccount struct {
	command.RootCommand
	Name string
}

// RechargeMoney is user command
type RechargeMoney struct {
	command.RootCommand
	Amount int64
}

// WithdrawMoney is user command
type WithdrawMoney struct {
	command.RootCommand
	Amount int64
}
