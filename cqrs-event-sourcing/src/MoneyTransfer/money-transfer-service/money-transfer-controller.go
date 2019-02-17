package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	transferCommand "./transfer-command"
	transfer_service "./transfer-service"
	nats "github.com/nats-io/go-nats"
)

var natsServer = flag.String("nats", "nats:4222", "NATs server URI")
var natsClient *nats.Conn

func init() {
	var err error
	for i := 0; i < 5; i++ {
		natsClient, err = nats.Connect("nats://" + *natsServer)
		if err == nil {
			break
		}
		fmt.Println("Waiting before connectiong to NATS at: ", *natsServer)
		time.Sleep(1 * time.Microsecond)
	}

	if err != nil {
		log.Fatal("error establishing connection to nats: ", err)
	}

	fmt.Println("connected to nats at: ", natsClient.ConnectedUrl())
}

func main() {
	log.Println("starting command service on port 8080")
	http.HandleFunc("/api/transfers", transferHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func transferHandler(writer http.ResponseWriter, request *http.Request) {
	if request.Method == "POST" {
		insertNewTransfer(writer, request)
	}
}

func insertNewTransfer(writer http.ResponseWriter, request *http.Request) {

	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	var detail = transferCommand.TransferDetail{}
	_ = json.Unmarshal(body, &detail)
	log.Println(detail.Description)
	transfer_service.TransferMoney(detail)

	defer request.Body.Close()

	// publish to query side
	fmt.Println("publish to query side")
	_ = natsClient.Publish("transfer.inserted", body)
}
