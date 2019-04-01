package trans

// MoneyTransferCreated is event
type MoneyTransferCreated struct {
	Details TransferDetail
}

// CreditRecorded is event trigger when account credited
type CreditRecorded struct {
	Details TransferDetail
}

// DebitRecorded is event trigger when account debited
type DebitRecorded struct {
	Details TransferDetail
}

// FailDebitRecorded is event trigger when account debit fail
type FailDebitRecorded struct {
	Details TransferDetail
}
