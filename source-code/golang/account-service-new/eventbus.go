package eventsourcing

// EventBus dinh nghia method de publish event
type EventBus interface {
	Publish(event Event, bucket, subset string) error
}