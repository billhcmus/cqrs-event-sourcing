package storage

import (
	"log"
	"fmt"
	"eventstore/proto"
	"github.com/jinzhu/gorm/dialects/postgres"
	"time"
	"eventstore"
	"github.com/jinzhu/gorm"
)

// Tx is type alias
type Tx = *gorm.DB

// Client for access to Postgres
type Client struct {
	db *gorm.DB
}

// CreateClient initialize the db
func CreateClient(dbURL string, logmode bool) (eventstore.EventStore, error) {
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


func (c *Client)save(events []*pb.BaseEvent, version uint64) error {
	if len(events) == 0 {
		return nil
	}
	var err error
	for i,event := range events {
		eventDB := EventDB {
			ID: event.EventId,
			Type: event.EventType,
			Timestamp: time.Unix(int64(event.Timestamp), 0),
			AggregateID: event.AggregateId,
			AggregateType: event.AggregateType,
			Version: 1 + version + uint64(i),
		}

		if (event.RawData != nil) {
			eventDB.RawData.RawMessage = event.RawData
		}

		err = c.db.Create(&eventDB).Error
		
		if err != nil {
			return err
		}
	}

	aggregateID := events[0].AggregateId
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
func (c *Client)Save(events []*pb.BaseEvent, version uint64) error {
	//c.db.Exec("SET TRANSACTION ISOLATION LEVEL SERIALIZABLE")

	tx := c.db.Begin()
	//var eventsDB EventDB
	// if err := tx.Raw("SELECT * FROM events_log WHERE id = ? for update", events[0].AggregateId).Scan(&eventsDB).Error; err != nil {
	// 	tx.Rollback();
	// 	return err
	// }
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
func (c *Client)Load(aggregateID string)([]*pb.BaseEvent, error) {
	tx := c.db.Begin()
	eventsDB := []EventDB{}
	c.db.Find(&eventsDB, "aggregate_id = ?", aggregateID).Order("version")

	if len(eventsDB) == 0 {
		log.Printf("id %v", aggregateID)
		return nil, nil
	}
	res := make([]*pb.BaseEvent, len(eventsDB))
	for i,event := range(eventsDB) {
		res[i] = &pb.BaseEvent {
			AggregateId: event.AggregateID,
			AggregateType: event.AggregateType,
			Version: event.Version,
			Timestamp: uint64(event.Timestamp.Unix()),
			EventType: event.Type,
			RawData: event.RawData.RawMessage,
		}
	}
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	return res, nil
}