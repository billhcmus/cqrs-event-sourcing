package natsbus

import (
	"encoding/json"
	"github.com/billhcmus/cqrs/pkg/event"
	"github.com/billhcmus/cqrs/pkg/bus"
	nats "github.com/nats-io/go-nats"
)

// Client in order to access to event bus
type Client struct {
	Options nats.Options
}

// CreateClient return event bus client
func CreateClient(URL string) (bus.IEventBus, error) {
	opts := nats.DefaultOptions
	opts.Url = URL
	
	return &Client{
		opts,
	}, nil
}

// Publish event to event bus
func (client *Client)Publish(event event.Event, topic string) error {
	natsCli, err := client.Options.Connect()
	if err != nil {
		return err
	}
	defer natsCli.Close()

	msg, err := json.Marshal(event)
	if err != nil {
		return err
	}

	subj := topic
	natsCli.Publish(subj, msg)
	natsCli.Flush()

	err = natsCli.LastError()
	return err
}