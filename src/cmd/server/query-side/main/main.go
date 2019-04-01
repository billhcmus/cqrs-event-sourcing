package main

import (
	"runtime"
	"encoding/json"
	"github.com/billhcmus/cqrs/pkg/event"
	"github.com/sirupsen/logrus"
	"github.com/nats-io/go-nats"
)

const (
	topic = "account"
)


func init() {
	natsCli, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		logrus.Fatal("NATS: ", err)
	}
	logrus.Info("Connected to ", nats.DefaultURL)

	_,_ = natsCli.Subscribe(topic, func(msg *nats.Msg) {
		data := event.Event{}
		_ = json.Unmarshal(msg.Data, &data)
		logrus.Info(data)
	})
}

func main() {
	runtime.Goexit()
}