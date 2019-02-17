package transfer_command

type TransferDetail struct {
	FromAccount string `json:"from"`
	ToAccount string `json:"to"`
	Amount string `json:"amount"`
	Date string `json:"date"`
	Description string `json:"description"`
}


type TransferMoneyCommand interface {
	CreateTransferMoneyCommand(detail TransferDetail)
}