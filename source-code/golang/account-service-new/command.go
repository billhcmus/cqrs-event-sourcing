package eventsourcing

// Command chua cac medthod lay cac thong tin cua aggregate
type Command interface {
	GetType() string
	GetAggregateID() string
	GetAggregateType() string
	IsValid() bool
	GetVersion() uint64
}

// BaseCommand chua cac thong tin co ban ma cac command co ve aggregate
type BaseCommand struct {
	Type          string `json:"type"`
	AggregateID   string `json:"aggregate_id"`
	AggregateType string `json:"aggregate_type"`
	Version       uint64 `json:"version"`
}

// GetType tra ve kieu cua aggregate ma command nay thuc hien
func (base BaseCommand) GetType() string {
	return base.Type
}

// GetAggregateID tra ve ID cua aggregate
func (base BaseCommand) GetAggregateID() string {
	return base.AggregateID
}

// GetAggregateType tra ve Type cua aggregate
func (base BaseCommand) GetAggregateType() string {
	return base.AggregateType
}

// IsValid check aggregate validation
func (base BaseCommand) IsValid() bool {
	return true
}

// GetVersion return current version of aggregate
func (base BaseCommand) GetVersion() uint64 {
	return base.Version
}
