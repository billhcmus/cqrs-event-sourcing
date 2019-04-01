package cache

import (
	"reflect"
	"github.com/billhcmus/cqrs/pkg/aggregate"
)

// ICache is interface which contain method check state of aggregate
type ICache interface {
	SaveAggregateInstanceToCache(aggregate.IAggregate) error
	GetAggregateInstanceFromCache(ID string, aggreagateType reflect.Type) (aggregate.IAggregate, error)
}