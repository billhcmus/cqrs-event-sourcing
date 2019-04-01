package redis

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/billhcmus/cqrs/common"
	"github.com/billhcmus/cqrs/pkg/aggregate"
	"github.com/billhcmus/cqrs/pkg/cache"
	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"
)

var aggregateRegistry map[string]reflect.Type

// AggregateRegister is implementation of Aggregate interface
type AggregateRegister struct{}

// CreateAggregateRegister create register instance
func CreateAggregateRegister() IAggregateRegister {
	return &AggregateRegister{}
}

// IAggregateRegister is interface contain Get and Set method for register
type IAggregateRegister interface {
	Get(name string) (reflect.Type, error)
	Set(source interface{})
}

// Get function
func (reg *AggregateRegister) Get(name string) (reflect.Type, error) {
	rawType, ok := aggregateRegistry[name]
	if !ok {
		return nil, fmt.Errorf("Registry not contain type %v", name)
	}
	return rawType, nil
}

// Set function
func (reg *AggregateRegister) Set(source interface{}) {
	rawType, name := common.GetTypeName(source)
	aggregateRegistry[name] = rawType
}

// RedisCache is implementation of redis cache
type RedisCache struct {
	client *redis.Client
}

// CreateRedisCacheClient create redis client
func CreateRedisCacheClient(cacheURL string) (cache.ICache, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cacheURL,
		Password: "",
		DB:       0,
	})

	_, err := client.Ping().Result()
	if err != nil {
		return nil, err
	}
	return &RedisCache{
		client,
	}, nil
}

// SaveAggregateInstanceToCache save instance to cache
func (redis *RedisCache) SaveAggregateInstanceToCache(aggregate aggregate.IAggregate) error {
	rawData, err := json.Marshal(aggregate)
	if err != nil {
		return err
	}
	_, err = redis.client.Set(aggregate.GetID(), rawData, 0).Result()

	return err
}

// GetAggregateInstanceFromCache get instance of aggregate from ID
func (redis *RedisCache) GetAggregateInstanceFromCache(ID string, aggregateType reflect.Type) (aggregate.IAggregate, error) {
	rawData, err := redis.client.Get(ID).Bytes()
	if err != nil {
		return nil, err
	}

	instance := reflect.New(aggregateType).Interface().(aggregate.IAggregate)

	err = json.Unmarshal(rawData, instance)
	if err != nil {
		return nil, err
	}

	logrus.Infof("[Redis] instance %v", instance)
	return instance, nil
}
