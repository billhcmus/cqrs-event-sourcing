package main

import (
	"github.com/jinzhu/gorm"
	"github.com/billhcmus/cqrs/pkg/cache/redis"
	"github.com/billhcmus/cqrs/internal/storage"
	natsbus "github.com/billhcmus/cqrs/pkg/bus/nats"
	cmdhandler "github.com/billhcmus/cqrs/pkg/command-handler"
	"github.com/billhcmus/cqrs/pkg/event"
	"github.com/billhcmus/cqrs/pkg/repository"
	"github.com/billhcmus/cqrs/pkg/service/account-service"
	"github.com/google/uuid"
	"github.com/nats-io/go-nats"
	"github.com/sirupsen/logrus"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

const (
	ebURL = nats.DefaultURL
	dbURL = "postgres://postgres:1234@localhost/eventstore?sslmode=disable"
	cacheURL = "localhost:6379"
)

func init() {
	db, err := gorm.Open("postgres", "postgres://postgres:1234@localhost/eventstore?sslmode=disable")

	if err != nil {
		panic(err)
	} else {
		logrus.Info("[Init] connect to database successfully")
	}
	defer db.Close()

	err = db.Exec(`
	 	CREATE TABLE IF NOT EXISTS aggregates (
			id UUID NOT NULL,
			type TEXT NOT NULL,
			version BIGINT NOT NULL,

			PRIMARY KEY (id)
		 );
	  `).Error

	if err != nil {
		panic(err)
	} else {
		logrus.Info("[Init] create table aggregates successfully")
	}

	err = db.Exec(`
		CREATE TABLE IF NOT EXISTS events_log (
			id UUID NOT NULL,
			aggregate_id UUID NOT NULL,
			aggregate_type TEXT NOT NULL,
			type TEXT NOT NULL,
			timestamp TIMESTAMP WITH TIME ZONE NOT NULL,
			version BIGINT NOT NULL,
			data JSONB NOT NULL,
			
			PRIMARY KEY (id)
		)
	`).Error

	if err != nil {
		panic(err)
	} else {
		logrus.Info("[Init] create table events_log successfully")
	}
}

func main() {
	eb, err := natsbus.CreateClient(ebURL)

	if err != nil {
		logrus.Fatal(err)
	}
	logrus.Info("[Main] config eventbus complete")
	
	es, err := storage.CreatePostgresClient(dbURL, true)
	if err != nil {
		logrus.Fatal(err)
	}
	logrus.Info("[Main] config event store complete")

	ch, err := redis.CreateRedisCacheClient(cacheURL)
	if err != nil {
		logrus.Fatal(err)
	}
	logrus.Info("[Main] config cache complete")
	
	repo := repository.CreateNewRepository(es, eb, ch)

	// Register events
	eventregister := event.CreateEventRegister()
	eventregister.Set(account.AccountCreated{})
	eventregister.Set(account.MoneyRecharged{})
	eventregister.Set(account.WithdrawPerformed{})

	var acc account.Account

	handler := cmdhandler.CreateCommandHandler(repo, &acc, "account")

	uuidv4, _ := uuid.NewRandom()
	userID := uuidv4.String()
	var create account.CreateAccount
	create.AggregateID = userID
	create.Name = "thuyenpt"
	create.Type = "create"
	create.Version = 0
	err = handler.Handle(create)
	if err != nil {
		logrus.Fatal(err)
	}

	var recharge account.RechargeMoney
	recharge.AggregateID = userID
	recharge.Amount = 1000
	recharge.Version = 1
	recharge.Type = "recharge"
	err = handler.Handle(recharge)
	if err != nil {
		logrus.Fatal(err)
	}

	var withdraw account.WithdrawMoney
	withdraw.AggregateID = userID
	withdraw.Amount = 1000
	withdraw.Version = 2
	withdraw.Type = "withdraw"
	err = handler.Handle(withdraw)
	if err != nil {
		logrus.Fatal(err)
	}

	logrus.Info(acc)
}
