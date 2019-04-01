package trans

import (
	"github.com/billhcmus/cqrs/pkg/command"
)

// CreateTransaction create new payment
type CreateTransaction struct {
	command.RootCommand
	Details TransferDetail
}

// RecordDebit is command execute record debit
type RecordDebit struct {
	command.RootCommand
}

// RecordCredit is command execute record credit
type RecordCredit struct {
	command.RootCommand
}

// RecordDebitFailed is command execute debit fail
type RecordDebitFailed struct {
	command.RootCommand
}