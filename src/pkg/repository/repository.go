package repository

import (
	"reflect"
	"github.com/billhcmus/cqrs/pkg/cache"
	"github.com/sirupsen/logrus"
	"github.com/billhcmus/cqrs/pkg/aggregate"
	"github.com/billhcmus/cqrs/pkg/bus"
	"github.com/billhcmus/cqrs/pkg/eventstore"
)

// Repository manage aggregate instance
type Repository struct {
	es eventstore.IEventStore
	eb bus.IEventBus
	ch cache.ICache
}

// CreateNewRepository tao mot repository moi
func CreateNewRepository(eventstore eventstore.IEventStore, eventbus bus.IEventBus, ch cache.ICache) *Repository {
	return &Repository{
		eventstore,
		eventbus,
		ch,
	}
}

// GetAggregateInstanceFromCache save aggregate instance to cache 
func (repo *Repository) GetAggregateInstanceFromCache(ID string, aggregateType reflect.Type) (aggregate.IAggregate, error) {
	return repo.ch.GetAggregateInstanceFromCache(ID, aggregateType)
}
// Replay return current state of the aggregate
func (repo *Repository) Replay(aggregate aggregate.IAggregate, ID string) error {
	events, err := repo.es.Load(ID)
	if err != nil {
		return err
	}

	for _,event := range events {
		logrus.Infof("[Repository] replaying %v", event)
		aggregate.ApplyChangeHelper(aggregate, event, false) // truong hop nay se khong commit event
	}
	return nil
}

// SaveAndPublish event to event store
func (repo *Repository) SaveAndPublish(aggregate aggregate.IAggregate, baseversion uint64, topic string) error {

	err := repo.es.Save(aggregate.UnCommited(), baseversion)
	if err != nil {
		return err
	}

	err = repo.PublishEvents(aggregate, topic)
	if err != nil {
		return err
	}
	
	aggregate.ClearUncommited()

	err = repo.ch.SaveAggregateInstanceToCache(aggregate)
	if err != nil {
		return err
	}

	return nil
}

// PublishEvents to an eventbus
func (repo *Repository) PublishEvents(aggregate aggregate.IAggregate, topic string) error {
	var err error

	for _,event := range aggregate.UnCommited() {
		if err = repo.eb.Publish(event,topic); err != nil {
			return err
		}
	}

	return nil
}