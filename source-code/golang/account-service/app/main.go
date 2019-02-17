package main

import (
	"github.com/google/uuid"
	"fmt"
	"log"
	"context"
	"errors"
	"eventsource"
)

//////////////////////////////////////
// Aggregate definition area

// User la aggregate
type User struct {
	eventsourcing.RootAggregate
	FirstName string
	LastName string
	Balance uint64
}

// Implement lai cac method dinh nghia trong interface Aggregate

// AggregateType tra ve kieu cua aggregate
func (user *User)AggregateType() string {
	return "user"
}

// TableName la bang trong eventstore
func (user *User)TableName() string {
	return "users"
}

//////////////////////////////////////
// Events definition area

// CreateAccountEvent is event
type CreateAccountEvent struct {
	ID string `json:"id"`
	FirstName string `json:"first_name"`
	LastName string `json:"last_name"`
}

// Apply event to aggregate
func (eventData CreateAccountEvent) Apply(agg eventsourcing.Aggregate, event eventsourcing.Event) {
	user := agg.(*User)
	user.ID = eventData.ID
	user.FirstName = eventData.FirstName
	user.LastName = eventData.LastName
	user.CreatedAt = event.Timestamp
	user.Balance = 0
}

// AggregateType get type of target aggregate
func (CreateAccountEvent) AggregateType() string {
	return "user"
}

// Action is the performed action
func (CreateAccountEvent) Action() string {
	return "account_created"
}

// Version is the event version
func (CreateAccountEvent) Version() uint64 {
	return 1
}

////////////////////////////////////
// Command definition area

// ValidationError is custom vailidation error
type ValidationError error

// CreateValidationError return validationError
func CreateValidationError(message string) error {
	return errors.New(message).(ValidationError)
}

func validateFirstName(firstName string) error {
	length := len(firstName)

	if length < 3 {
		return CreateValidationError("First name to short")
	} else if length > 31 {
		return CreateValidationError("First name to long")
	}
	return nil
}

// CreateAccount is command
type CreateAccount struct {
	FirstName string
	LastName string
}

// Validate check data valid
func (create CreateAccount)Validate(ctx context.Context, tx eventsourcing.Tx, aggregate eventsourcing.Aggregate) error {
	return validateFirstName(create.FirstName)
}

// BuildEvent return the CreateAccount event
func (create CreateAccount)BuildEvent(ctx context.Context) (eventsourcing.EventData, interface{}, error) {
	uuidv4,_ := uuid.NewRandom()
	return CreateAccountEvent {
		ID: uuidv4.String(),
		FirstName: create.FirstName,
		LastName: create.LastName,
	}, nil, nil
}

// AggregateType return the target aggregate type
func (create CreateAccount)AggregateType() string {
	return "user"
}

func registerUser() {
	eventsourcing.Register(CreateAccountEvent{})
}

func main() {
	var dbURL = "postgres://postgres:1234@localhost/eventstore?sslmode=disable"
	// config db
	err := eventsourcing.Init(dbURL)
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("Config database complete")
	}

	eventsourcing.DB.LogMode(true)

	registerUser()

	var user User

	command := CreateAccount {
		FirstName: "Bill",
		LastName: "Gates",
	}

	_,err = eventsourcing.Execute(context.Background(), command, &user)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Print(user)
}