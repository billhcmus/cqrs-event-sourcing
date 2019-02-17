package main

import (
	"../pb"
	"encoding/json"
	"fmt"
	"github.com/gogo/protobuf/proto"
	"github.com/nats-io/go-nats"
	"github.com/satori/go.uuid"
	"log"
	"runtime"
)

const (
	aggregate = "Payment"
	event = "PaymentCreated"
)

var nc *nats.Conn

func publishPaymentCreated(payment *pb.PaymentCreatedCommand) {

	log.Println("Connected to " + nats.DefaultURL)

	eventData, _ := json.Marshal(payment)

	eventId, _ := uuid.NewV4()

	event := pb.Event{
		AggregateId: payment.PaymentId,
		AggregateType: aggregate,
		EventId: eventId.String(),
		EventType: event,
		EventData: string(eventData),
	}

	subject := "Payment.PaymentCreated"
	data,_ := proto.Marshal(&event)

	// Publish
	_ = nc.Publish(subject, data)
	fmt.Println("Publish to ", subject)
}

func init() {
	// connect to NATs server
	log.Println("Connecting ", nats.DefaultURL)
	var err error
	nc,err = nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatal("Error: ", err)
	}
}

func main() {
	var payment = new(pb.PaymentCreatedCommand)

	u4, err := uuid.NewV4()

	if err != nil {
		log.Fatal("something went wrong ", err)
	}

	payment.PaymentId = u4.String()
	payment.SenderId = "bill"
	payment.ReceiverId = "nic"
	payment.Status = "NEW"
	payment.Amount = 10000

	var item = new(pb.PaymentCreatedCommand_Item)
	item.ItemId = "123"
	item.ItemName = "Banh chung"
	item.UnitPrice = 1000000
	item.Quantity = 2
	listItem := []*pb.PaymentCreatedCommand_Item{item}
	payment.ListItems = listItem

	publishPaymentCreated(payment)

	runtime.Goexit()
}