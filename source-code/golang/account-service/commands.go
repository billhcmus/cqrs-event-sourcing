package eventsourcing

import (
	"fmt"
	"reflect"
	"context"
)

// Command use to validate data, build event
type Command interface {
	BuildEvent(context.Context) (event EventData, nonPersisted interface{}, err error)
	Validate(context.Context, Tx, Aggregate) error
	AggregateType() string
}

// Execute a command to an aggregate
func Execute(ctx context.Context, command Command, aggregate Aggregate) (Event, error) {
	tx := DB.Begin()
	event, err := ExecuteTx(ctx, tx, command, aggregate)
	if err != nil {
		tx.Rollback()
		return Event{}, err
	}

	err = tx.Commit().Error
	if err != nil {
		tx.Rollback()
		return Event{}, err
	}

	return event, nil
}

// ExecuteTx execute the given command to the aggregate and return created event
func ExecuteTx(ctx context.Context, tx Tx, command Command, aggregate Aggregate) (Event, error) {
	var err error

	// verify aggregate
	rv := reflect.ValueOf(aggregate)

	if rv.Kind() != reflect.Ptr {
		return Event{}, fmt.Errorf("Aggregate must be pointer not %s ", reflect.TypeOf(aggregate))
	}

	if rv.IsNil() {
		return Event{}, fmt.Errorf("Aggregate was nil")
	}

	if command.AggregateType() != aggregate.AggregateType() {
		return Event{}, fmt.Errorf("Command aggregate's type (%s) and aggregate type (%s) are mismatch",
		 command.AggregateType(), aggregate.AggregateType())
	}

	if aggregate.GetID() != "" {
		tx.Set("gorm:query_option", "FOR UPDATE").First(aggregate)
	}

	err = command.Validate(ctx, tx, aggregate)
	if err != nil {
		return Event{}, err
	}

	eventdata, nonPersisted, err := command.BuildEvent(ctx)
	if err != nil {
		return Event{}, err
	}
	event := buildBaseEvent(eventdata, nonPersisted, aggregate.GetID())
	event.Data = eventdata
	
	event.apply(aggregate)

	event.AggregateID = aggregate.GetID()

	err = tx.Save(aggregate).Error
	if err != nil {
		return Event{}, err
	}

	EventToStore, err := event.Serialize()
	if err != nil {
		return Event{}, err
	}

	err = tx.Create(&EventToStore).Error
	if err != nil {
		return Event{}, err
	}

	// dispatch to other side

	return event, nil
}