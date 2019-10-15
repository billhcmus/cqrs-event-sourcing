package repository

import (
	"eventsourcing"
	"fmt"
	"log"
	"reflect"
	"encoding/json"
	"eventsourcing/proto"
	"time"
	"context"
)

// Repository tao Aggregate, luu tru event va publish event
type Repository struct {
	eventStore pb.EventStoreClient
	eventBus   eventsourcing.EventBus
}

// CreateNewRepository tao mot repository moi
func CreateNewRepository(eventstore pb.EventStoreClient, eventbus eventsourcing.EventBus) *Repository {
	return &Repository{
		eventstore,
		eventbus,
	}
}

// Replay function tra ve trang thai cuoi cung cua aggregate
func (repo *Repository) Replay(aggregate eventsourcing.Aggregate, ID string) error {
	// Load events from eventstore (call grpc api Load)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	loadEventRequest := &pb.LoadEventsRequest{AggregateId: ID}
	log.Printf("replay id %v", ID)
	res, err := repo.eventStore.Load(ctx, loadEventRequest)
	if err != nil {
		return err
	}
	if res.Status.Code != 0 {
		return fmt.Errorf(res.Status.Error)
	}

	register := eventsourcing.CreateEventRegister()

	// Deserialize data
	for _, event := range res.Events {
		// convert type *pb.Event to eventsourcing.Event
		var err error
		var dataType reflect.Type
		dataType, err = register.Get(event.EventType)
		if err != nil {
			return err
		}
		data := reflect.New(dataType).Interface()
		if err = json.Unmarshal(event.RawData, data); err != nil {
			return err
		}
		e := eventsourcing.Event {
			BaseEvent: *event,
			Data: data,
		}

		aggregate.ApplyChangeHelper(aggregate, e, false) // truong hop nay se khong commit event
	}
	return nil
}

// Save event to event store
func (repo *Repository) Save(aggregate eventsourcing.Aggregate, version uint64) error {
	// Prepare event type [] *pb.Event
	//var err error
	log.Println("call save api")
	events := make([]*pb.BaseEvent, len(aggregate.UnCommited()))
	for i, event := range aggregate.UnCommited() {
		events[i] = &pb.BaseEvent {
			EventId: event.EventId,
			EventType: event.EventType,
			AggregateId: event.AggregateId,
			AggregateType: event.AggregateType,
			Timestamp: event.Timestamp,
			Version: 1 + version + uint64(i),
		}

		if event.Data != nil {
			events[i].RawData, _ = json.Marshal(event.Data)
		}
	}
	saveEventsRequest := &pb.SaveEventsRequest{
		Events:  events,
		Version: version,
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	response, err := repo.eventStore.Save(ctx, saveEventsRequest)
	if response != nil {
		log.Printf("response code %v", response.Status.Code)
	}

	if err != nil {
		log.Println("event ", events[0].EventType)
		log.Println("err ", err)
		return err
	}


	return nil
}

// PublishEvents to an eventbus
func (repo *Repository) PublishEvents(aggregate eventsourcing.Aggregate, bucket, subset string) error {
	var err error

	for _, event := range aggregate.UnCommited() {
		
		if err = repo.eventBus.Publish(event, bucket, subset); err != nil {
			return err
		}
	}

	return nil
}