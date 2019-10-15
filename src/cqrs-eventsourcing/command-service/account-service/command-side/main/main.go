package main

import (
	"eventsourcing"
	"eventsourcing/account-service/account"
	asynccommandbus "eventsourcing/commandbus-Impl/async"
	commandhandler "eventsourcing/commandhandler-Impl"
	eventbus "eventsourcing/eventbus-Impl"
	pb "eventsourcing/proto"
	"eventsourcing/repository"
	"flag"
	"fmt"
	"github.com/nats-io/go-nats"
	"log"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc"
)

const (
	esURL = "localhost:7777"
	ebURL = nats.DefaultURL
)

func main() {
	flag.Parse()

	// Create client for access database => can ta interface cho cac db khac nhau

	conn, err := grpc.Dial(esURL, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect to eventstore: %v", err)
	}
	defer conn.Close()

	esClient := pb.NewEventStoreClient(conn)

	eventbus, err := eventbus.CreateClient(ebURL)
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("Config eventbus complete")
	}

	repository := repository.CreateNewRepository(esClient, eventbus)

	// Register events
	eventregister := eventsourcing.CreateEventRegister()
	eventregister.Set(account.AccountCreated{})
	eventregister.Set(account.MoneyRecharged{})
	eventregister.Set(account.WithdrawPerformed{})

	var acc account.Account

	cmdhandler := commandhandler.CreateNewCommandHandler(repository, &acc, "acc", "account")

	// Register command
	commandregister := eventsourcing.CreateCommandRegister()
	commandregister.Set(account.CreateAccount{}, cmdhandler)
	commandregister.Set(account.RechargeMoney{}, cmdhandler)
	commandregister.Set(account.WithdrawMoney{}, cmdhandler)

	commandbus := asynccommandbus.CreateBus(1, commandregister)
	commandbus.StartBus()

	uuidv4, _ := uuid.NewRandom()
	userID := uuidv4.String()
	var create account.CreateAccount
	create.AggregateId = userID
	create.Name = "thuyenpt"
	commandbus.HandleCommand(create)

	time.Sleep(500 * time.Millisecond)
	recharge := account.RechargeMoney{Amount: 1000}
	recharge.AggregateId = userID
	recharge.Version = 1
	commandbus.HandleCommand(recharge)

	time.Sleep(500 * time.Millisecond)
	//userID := "d7e57147-6162-400e-9bfa-09295302b394"
	withdrawl := account.WithdrawMoney{Amount: 500}
	withdrawl.AggregateId = userID
	withdrawl.Version = 2
	commandbus.HandleCommand(withdrawl)

	time.Sleep(500 * time.Millisecond)
	fmt.Println(acc)
}
