package eventsourcing

import (
	"time"
)

// Aggregate interface được implement bởi mỗi aggregate được định nghĩa sau này. 
type Aggregate interface {
	GetID() string
	incrementVersion()
	updateUpdatedAtField(time.Time)
	AggregateType() string
	TableName() string
}

// RootAggregate la aggregate goc cac aggregate sau ebedded root vao no
type RootAggregate struct {
	ID string `json:"id" gorm:"column:id;type:uuid;primary_key:true"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at"`
	DeletedAt time.Time `json:"deleted_at" gorm:"column:deleted_at"`
	Version uint64 `json:"version" gorm:"column:version"`
}

// Implement Aggregate interface on RootAggregate

// GetID return aggregate's ID
func (agg RootAggregate) GetID() string {
	return agg.ID
}

func (agg RootAggregate) incrementVersion() {
	agg.Version++
}

func (agg RootAggregate) updateUpdatedAtField(time time.Time) {
	agg.UpdatedAt = time
}