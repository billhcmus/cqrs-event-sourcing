package main

import (
	"../pb"
	"encoding/json"
	"fmt"
	"github.com/gogo/protobuf/proto"
	"github.com/nats-io/go-nats"
	"log"
	"runtime"
)

const (
	queue = "PaymentQueue"
)

func init() {
	// connect to nats server
	log.Println("Connecting")
	nc, err := nats.Connect(nats.DefaultURL)

	if err != nil {
		log.Fatal("Error: ", err)
	}
	log.Println("Connected\nSubscribe")

	_, _ = nc.QueueSubscribe(subject, queue, func(msg *nats.Msg) {
		event := pb.Event{}
		_ = proto.Unmarshal(msg.Data, &event)
		data, _ := json.Marshal(event)
		fmt.Println(string(data))
	})
}

func main() {
	runtime.Goexit()
}