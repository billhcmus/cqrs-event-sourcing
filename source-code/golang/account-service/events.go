package eventsourcing

import (
	"github.com/google/uuid"
	"strconv"
	"reflect"
	"encoding/json"
	"github.com/jinzhu/gorm/dialects/postgres"
	"time"
)

var eventRegistry = map[string]reflect.Type{}

// Event is in-memory event
type Event struct {
	ID string `json:"id"`
	Timestamp time.Time `json:"timestamp"`
	AggregateID string `json:"aggregate_id"`
	AggregateType string `json:"aggregate_type"`
	Action string `json:"action"`
	Version uint64 `json:"version"`
	Type string `json:"type"`
	Data interface{} `json:"data"` // empty interface may hold values of any type
	NonPersisted interface{} `json:"-"`
}

// EventData duoc implement boi cac event duoc dinh nghia sau nay
type EventData interface {
	AggregateType() string
	Action() string
	Version() uint64
	// Apply Event cho Aggregate
	Apply(Aggregate, Event)
}

// apply function call Apply method in event's data to update Aggregate
func (event Event) apply(agg Aggregate) {
	event.Data.(EventData).Apply(agg, event)
	agg.incrementVersion()
	agg.updateUpdatedAtField(event.Timestamp)
}

// StoreEvent is struct to be serialized to eventstore or deserialized from eventstore
type StoreEvent struct {
	ID string `json:"id" gorm:"column:id;type:uuid;primary_key:true"`
	Timestamp time.Time `json:"timestamp"`
	AggregateID string `json:"aggregate_id" gorm:"column:aggregate_id;type:uuid"`
	AggregateType string `json:"aggregate_type"`
	Action string `json:"action"`
	Version uint64 `json:"version"`
	Type string `json:"type"`

	RawData postgres.Jsonb `json:"-" gorm:"type:jsonb;column:data"`
}

// Register to register event's type
func Register(events ...EventData) {
	for _, event := range events {
		eventType := event.AggregateType() + 
		"." + event.Action() + 
		"." + strconv.FormatUint(event.Version(), 10)
		eventRegistry[eventType] = reflect.TypeOf(event)
	}
}

func buildBaseEvent(eventdata EventData, nonPersisted interface{}, aggregatedID string) Event {
	event := Event{}
	uuidv4,_ := uuid.NewRandom()

	event.ID = uuidv4.String()
	event.Timestamp = time.Now().UTC()
	event.AggregateID = aggregatedID
	event.AggregateType = eventdata.AggregateType()
	event.Action = eventdata.Action()
	event.Version = eventdata.Version()
	event.Type = eventdata.AggregateType() + "." + eventdata.Action()
	event.NonPersisted = nonPersisted
	
	return event
}

// Serialize serialized Event to StoreEvent to be stored in event store
func (event Event) Serialize() (StoreEvent, error) {
	res := StoreEvent{}
	var err error

	res.ID = event.ID
	res.Timestamp = event.Timestamp
	res.AggregateID = event.AggregateID
	res.AggregateType = event.AggregateType
	res.Action = event.Action
	res.Version = event.Version
	res.Type = event.Type

	res.RawData.RawMessage, err = json.Marshal(event.Data)
	if err != nil {
		return StoreEvent{}, err
	}
	return res, nil
}

// Deserialize deserialize from StoreEvent to Event
func (event StoreEvent) Deserialize() (Event, error) {
	res := Event{}
	var err error
	
	// Xac dinh event type
	eventType := event.AggregateType + 
	"." + event.Action + 
	"." + strconv.FormatUint(event.Version, 10)

	dataPtr := reflect.New(eventRegistry[eventType])
	dataValue := dataPtr.Elem()
	i := dataValue.Interface()
	err = json.Unmarshal(event.RawData.RawMessage, &i)

	if err != nil {
		return Event{}, err
	}

	res.ID = event.ID
	res.AggregateID = event.AggregateID
	res.AggregateType = event.AggregateType
	res.Timestamp = event.Timestamp
	res.Action = event.Action
	res.Version = event.Version
	res.Type = event.Type
	res.Data = i

	return res, nil
}

// TableName used by gorm to insert and retrieve events
func (event StoreEvent) TableName() string {
	return event.AggregateType + "s_events"
}