package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	nats "github.com/nats-io/go-nats"
)

type TransferDetail struct {
	FromAccount string `json:"from"`
	ToAccount   string `json:"to"`
	Amount      int    `json:"amount"`
	Date        string `json:"date"`
	Description string `json:"description"`
}

var natsClient *nats.Conn

var natsServer = flag.String("nats", "nats", "NATs server URI")

var prevData TransferDetail

func init() {
	flag.Parse()
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
	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		http.HandleFunc("/api/details", getDetailTransfer)
		log.Println("server listening on port 8081")
		log.Fatal(http.ListenAndServe(":8081", nil))
	}()

	go func() {
		if _, err := natsClient.Subscribe("transfer.inserted", handleTransferInserted); err != nil {
			log.Fatal(err)
		}
	}()

	wg.Wait()
}

func getDetailTransfer(writer http.ResponseWriter, request *http.Request) {
	log.Println("handle get")
	dataResponse, err := json.Marshal(prevData)

	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
	}

	_, _ = fmt.Fprint(writer, string(dataResponse))
}

func handleTransferInserted(msg *nats.Msg) {
	log.Println("receive event")
	prevData = TransferDetail{}

	if err := json.Unmarshal(msg.Data, &prevData); err != nil {
		log.Println("Unmarshal error")
	}
}
