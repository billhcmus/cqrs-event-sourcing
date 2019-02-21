package events

// AccountCreated is event
type AccountCreated struct {
	Name string `json:"name"`
}

// MoneyRecharged is event
type MoneyRecharged struct {
	Amount uint64 `json:"amount"`
}

// WithdrawlPerformed is event
type WithdrawlPerformed struct {
	Amount uint64 `json:"amount"`
}