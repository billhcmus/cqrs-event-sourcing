package eventstore

import (
	"log"
	"reflect"
	"fmt"
	"encoding/json"
	"github.com/jinzhu/gorm/dialects/postgres"
	"time"
	"eventsource"
	"github.com/jinzhu/gorm"
)

// Tx is a alias of *gorm.DB
type Tx = *gorm.DB

// DB to access to DB in app
var DB *gorm.DB

// Client for access to Postgres
type Client struct {
	db *gorm.DB
}

// CreateClient initialize the db
func CreateClient(dbURL string, logmode bool) (eventsourcing.EventStore, error) {
	db, err := gorm.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
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
	ID string `json:"id" gorm:"column:id;type:uuid;primary_key:true"`
	AggregateID string `json:"aggregate_id" gorm:"column:aggregate_id;type:uuid"`
	AggregateType string `json:"aggregate_type" gorm:"column:aggregate_type"`
	Type string `json:"type" gorm:"column:type"`
	Timestamp time.Time `json:"timestamp" gorm:"column:timestamp"`
	Version uint64 `json:"version" gorm:"column:version"`

	RawData postgres.Jsonb `json:"-" gorm:"type:jsonb;column:data"`
}

// AggregateDB is struct to be serialized to eventstore or deserialized from eventstore
type AggregateDB struct {
	ID string `json:"id" gorm:"column:id;type:uuid;primary_key;unique"`
	Type string `json:"type" gorm:"column:type"`
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


func (c *Client)save(events []eventsourcing.Event, version uint64) error {
	if len(events) == 0 {
		return nil
	}
	var err error

	for i,event := range events {
		eventDB := EventDB {
			ID: event.ID,
			Type: event.Type,
			Timestamp: event.Timestamp,
			AggregateID: event.AggregateID,
			AggregateType: event.AggregateType,
			Version: 1 + version + uint64(i),
		}

		if (event.Data != nil) {
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

	aggregateID := events[0].AggregateID
	aggregateType := events[0].AggregateType
	// version = 0 => new aggregate
	if version == 0 {
		aggregate := AggregateDB {
			ID: aggregateID,
			Type: aggregateType,
			Version: uint64(len(events)),
		}
		err = c.db.Save(&aggregate).Error
		if err != nil {
			return err
		}
		return nil
	}

	// query aggregate
	var aggregate AggregateDB
	err = c.db.First(&aggregate, "id = ?", aggregateID).Error

	if err != nil {
		return err
	}

	if aggregate.Version != version {
		return fmt.Errorf("Concurrent update error %s", aggregateID)
	}

	// update it
	aggregate.Version += uint64(len(events))
	err = c.db.Save(aggregate).Error
	if err != nil {
		return err
	}

	return nil
}

// Save store and update event, aggregate to db
func (c *Client)Save(events []eventsourcing.Event, version uint64) error {
	tx := c.db.Begin()

	err := c.save(events, version)
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
func (c *Client)Load(aggregateID string)([]eventsourcing.Event, error) {
	var err error
	var dataType reflect.Type
	eventsDB := []EventDB{}
	c.db.Find(&eventsDB, "aggregate_id = ?", aggregateID).Order("version")
	res := make([]eventsourcing.Event, len(eventsDB))
	register := eventsourcing.CreateEventRegister()

	for i,event := range(eventsDB) {
		
		// Deserialize event's data
		// Get Datatype
		dataType, err = register.Get(event.Type)
		if err != nil {
			return []eventsourcing.Event{}, err
		}
		data := reflect.New(dataType).Elem().Interface()

		if err = json.Unmarshal(event.RawData.RawMessage, &data); err != nil {
			return []eventsourcing.Event{}, err
		}

		res[i].ID = event.ID
		res[i].AggregateID = event.AggregateID
		res[i].AggregateType = event.AggregateType
		res[i].Version = event.Version
		res[i].Timestamp = event.Timestamp
		res[i].Type = event.Type
		res[i].Data = data
	}
	return res, nil
}