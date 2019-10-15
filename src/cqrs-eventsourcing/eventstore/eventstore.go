package eventstore

import (
	"eventstore/proto"
)

// EventStore luu tru event tu aggregate
type EventStore interface {
	Save(event []*pb.BaseEvent, version uint64) error
	Load(aggregateID string) ([]*pb.BaseEvent, error)
}
