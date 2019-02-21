package main

import (
	"flag"
	"eventsource/app/commands"
	"github.com/google/uuid"
	"eventsource/app/user"
	"eventsource/commandhandler"
	"eventsource/eventbus"
	"eventsource/app/events"
	"eventsource"
	"log"
	"eventsource/eventstore"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func main() {
	flag.Parse()
	var dbURL = "postgres://postgres:1234@localhost/eventstore?sslmode=disable"

	// Create client for access database => can ta interface cho cac db khac nhau
	store, err := eventstore.CreateClient(dbURL, true)
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("Config database complete")
	}

	eventbus, err := eventbus.CreateClient()
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("Config eventbus complete")
	}

	repository := eventsourcing.CreateNewRepository(store, eventbus)
	
	// Register events
	reg := eventsourcing.CreateEventRegister()
	reg.Set(events.AccountCreated{})
	reg.Set(events.MoneyRecharged{})
	reg.Set(events.WithdrawlPerformed{})

	var user user.User
	handler := commandhandler.CreateNewCommandHandler(repository, &user, "account", "account service")
	
	uuidv4,_ := uuid.NewRandom()
	userID := uuidv4.String()
	var create commands.CreateAccount
	create.AggregateID = userID
	create.Name = "thuyenpt"

	handler.Handle(create)

	recharge := commands.RechargeMoney{Amount: 1000,}
	recharge.AggregateID = userID
	recharge.Version = 1
	handler.Handle(recharge)
	
	withdrawl := commands.WithdrawlMoney{Amount: 500}
	withdrawl.AggregateID = userID
	withdrawl.Version = 2
	handler.Handle(withdrawl)
}
