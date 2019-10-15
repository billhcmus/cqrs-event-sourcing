package nats

import (
	"encoding/json"
	"eventsourcing"
	"log"
	nats "github.com/nats-io/go-nats"
)

// Client in order to access to event bus
type Client struct {
	Options nats.Options
}

// CreateClient return event bus client
func CreateClient(URL string) (eventsourcing.EventBus, error) {
	opts := nats.DefaultOptions
	opts.Url = URL
	
	return &Client{
		opts,
	}, nil
}

// Publish event to event bus
func (client *Client)Publish(event eventsourcing.Event, bucket, subset string) error {
	log.Printf("Published to bucket %s and subset is %s", bucket, subset)
	natsCli, err := client.Options.Connect()
	if err != nil {
		return err
	}
	defer natsCli.Close()

	msg, err := json.Marshal(event)
	if err != nil {
		return err
	}

	subj := bucket + "." + subset
	natsCli.Publish(subj, msg)
	natsCli.Flush()

	err = natsCli.LastError()
	return err
}