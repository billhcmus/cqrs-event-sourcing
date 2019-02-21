package commands

import (
	"eventsource"
)

// CreateAccount is command
type CreateAccount struct {
	eventsourcing.BaseCommand
	Name string
}

// RechargeMoney is command
type RechargeMoney struct {
	eventsourcing.BaseCommand
	Amount uint64
}

// WithdrawlMoney is command
type WithdrawlMoney struct {
	eventsourcing.BaseCommand
	Amount uint64
}