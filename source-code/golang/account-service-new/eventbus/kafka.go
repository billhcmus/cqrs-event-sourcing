package eventbus

import (
	"log"
	"eventsource"
)

// Client in order to access to event bus
type Client struct {

}

// CreateClient return event bus client
func CreateClient() (eventsourcing.EventBus, error) {
	return &Client{}, nil
}

// Publish event to event bus
func (client *Client)Publish(event eventsourcing.Event, bucket, subset string) error {
	log.Printf("Published to bucket %s and subset is %s", bucket, subset)
	return nil
}