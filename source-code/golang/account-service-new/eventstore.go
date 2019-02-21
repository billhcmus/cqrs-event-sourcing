package eventsourcing

// EventStore luu tru event tu aggregate
type EventStore interface {
	Save(event []Event, version uint64) error
	Load(aggregateID string) ([]Event, error)
}
