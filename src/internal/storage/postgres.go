package storage

import (
	"encoding/json"
	"fmt"
	"reflect"
	"time"
	"github.com/billhcmus/cqrs/pkg/event"
	"github.com/billhcmus/cqrs/pkg/eventstore"
	"github.com/jinzhu/gorm"
	"github.com/jinzhu/gorm/dialects/postgres"
)

// Tx is type alias
type Tx = *gorm.DB

// Client for access to Postgres
type Client struct {
	db *gorm.DB
}

// CreatePostgresClient initialize the db
func CreatePostgresClient(dbURL string, logmode bool) (eventstore.IEventStore, error) {
	db, err := gorm.Open("postgres", dbURL)
	if err != nil {
		return nil, err
	}
	if logmode {
		db.LogMode(true)
	}
	client := &Client{
		db,
	}
	return client, nil
}

// EventDB is struct to be serialized to eventstore or deserialized from eventstore
type EventDB struct {
	ID            string    `json:"id" gorm:"column:id;type:uuid;primary_key:true"`
	AggregateID   string    `json:"aggregate_id" gorm:"column:aggregate_id;type:uuid"`
	AggregateType string    `json:"aggregate_type" gorm:"column:aggregate_type"`
	Type          string    `json:"type" gorm:"column:type"`
	Timestamp     time.Time `json:"timestamp" gorm:"column:timestamp"`
	Version       uint64    `json:"version" gorm:"column:version"`

	RawData postgres.Jsonb `json:"-" gorm:"type:jsonb;column:data"`
}

// AggregateDB is struct to be serialized to eventstore or deserialized from eventstore
type AggregateDB struct {
	ID      string `json:"id" gorm:"column:id;type:uuid;primary_key;unique"`
	Type    string `json:"type" gorm:"column:type"`
	Version uint64 `json:"version" gorm:"column:version"`
}

// TableName used by gorm to insert and retrieve events
func (event EventDB) TableName() string {
	return "events_log"
}

// TableName used by gorm to insert and retrieve events
func (aggregate AggregateDB) TableName() string {
	return "aggregates"
}

func (c *Client) save(events []event.Event, baseversion uint64) error {
	if len(events) == 0 {
		return nil
	}
	var err error
	existed := true

	aggregateID := events[0].AggregateID
	aggregateType := events[0].AggregateType

	// query aggregate
	aggregate := AggregateDB{
		ID:      aggregateID,
		Type:    aggregateType,
		Version: 0,
	}

	err = c.db.First(&aggregate, "id = ?", aggregateID).Error
	if err != nil {
		existed = false
	}

	if !existed {
		err = c.db.Save(&aggregate).Error
		if err != nil {
			return err
		}
	}

	// lock row to insert events
	c.db.Set("gorm:query_option", "FOR UPDATE").First(&aggregate, "id = ?", aggregateID)

	if aggregate.Version != baseversion {
		return fmt.Errorf("Concurrent update error %s", aggregateID)
	}

	for i, event := range events {
		eventDB := EventDB{
			ID:            event.ID,
			Type:          event.Type,
			Timestamp:     event.Timestamp,
			AggregateID:   event.AggregateID,
			AggregateType: event.AggregateType,
			Version:       1 + baseversion + uint64(i),
		}

		if event.Data != nil {
			eventDB.RawData.RawMessage, err = json.Marshal(event.Data)
			if err != nil {
				return err
			}
		}

		err = c.db.Create(&eventDB).Error

		if err != nil {
			return err
		}
	}

	// update aggregate
	aggregate.Version += uint64(len(events))
	err = c.db.Save(&aggregate).Error
	if err != nil {
		return err
	}

	return nil
}

// Save store and update event, aggregate to db
func (c *Client) Save(events []event.Event, baseversion uint64) error {
	tx := c.db.Begin()

	err := c.save(events, baseversion)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit().Error
	if err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

// Load load events of aggregate from eventstore
func (c *Client) Load(aggregateID string) ([]event.Event, error) {
	var err error
	var dataType reflect.Type
	eventsDB := []EventDB{}
	c.db.Find(&eventsDB, "aggregate_id = ?", aggregateID).Order("version")
	if len(eventsDB) == 0 {
		return []event.Event{}, fmt.Errorf("aggregate does not exists")
	}

	res := make([]event.Event, len(eventsDB))
	register := event.CreateEventRegister()

	for i, e := range eventsDB {

		// Deserialize event's data
		// Get Datatype
		dataType, err = register.Get(e.Type)
		if err != nil {
			return []event.Event{}, err
		}
		data := reflect.New(dataType).Interface()

		if err = json.Unmarshal(e.RawData.RawMessage, &data); err != nil {
			return []event.Event{}, err
		}

		res[i].ID = e.ID
		res[i].AggregateID = e.AggregateID
		res[i].AggregateType = e.AggregateType
		res[i].Version = e.Version
		res[i].Timestamp = e.Timestamp
		res[i].Type = e.Type
		res[i].Data = data
	}
	return res, nil
}
