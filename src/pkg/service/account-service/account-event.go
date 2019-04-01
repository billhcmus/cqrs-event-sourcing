package account

// AccountCreated is event
type AccountCreated struct {
	Name string `json:"name"`
}

// MoneyRecharged is event
type MoneyRecharged struct {
	Amount int64 `json:"amount"`
}

// WithdrawPerformed is event
type WithdrawPerformed struct {
	Amount int64 `json:"amount"`
}

// AccountCredited is event
type AccountCredited struct {
	Amount int64 `json:"amount"`
	TransID string `json:"trans_id"`	
}

// AccountDebited is event
type AccountDebited struct {
	Amount int64 `json:"amount"`
	TransID string `json:"trans_id"`
}