package main

import (
	"eventsourcing"
	"encoding/json"
	"runtime"
	"log"
	"github.com/nats-io/go-nats"
)

const (
	subject = "acc.account"
)


func init() {
	natsCli, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatal("NATS: ", err)
	}
	log.Println("Connected to ", nats.DefaultURL)

	_,_ = natsCli.Subscribe(subject, func(msg *nats.Msg) {
		data := eventsourcing.Event{}
		_ = json.Unmarshal(msg.Data, &data)
		log.Println(data)
	})
}

func main() {
	runtime.Goexit()
}