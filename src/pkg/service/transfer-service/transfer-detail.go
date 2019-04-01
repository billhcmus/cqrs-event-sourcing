package trans

import (
	"time"
)

// TransferDetail is detail transaction
type TransferDetail struct {
	FromAccount string
	ToAccount string
	amount int64
	date time.Time
}