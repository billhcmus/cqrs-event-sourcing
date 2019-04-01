package cmdhandler

import (
	"github.com/billhcmus/cqrs/common"
	"reflect"
	"github.com/sirupsen/logrus"
	"github.com/billhcmus/cqrs/pkg/aggregate"
	"github.com/billhcmus/cqrs/pkg/repository"
	"github.com/billhcmus/cqrs/pkg/command"
)

// ICommandHandler is the interface which handle command
type ICommandHandler interface {
	Handle(command.RootCommand) error
}

// CommandHandler is ICommandHandler implementation 
type CommandHandler struct {
	repository *repository.Repository
	aggregate aggregate.IAggregate
	topic string
}

// CreateCommandHandler like constructor of the handler
func CreateCommandHandler(repo *repository.Repository, aggregate aggregate.IAggregate, topic string) *CommandHandler {
	return &CommandHandler{
		repository: repo,
		aggregate:  aggregate,
		topic: topic,
	}
}

// Handle is definition of CommandHandler interface
func (handler *CommandHandler) Handle(command command.ICommand) error {
	var err error
	
	// check from cache, if not available then try to replay
	version := command.GetVersion()
	if version != 0 {
		aggregateType, aggregateName := common.GetTypeName(handler.aggregate)
		logrus.Infof("[Command Handler] handle command of aggregate %v", aggregateName)
		instance, err := handler.repository.GetAggregateInstanceFromCache(command.GetAggregateID(), aggregateType)
		if err != nil {
			logrus.Infof("[Command Handler] aggregate not exist in cache, try to Replay it")
			logrus.Info("[Command Handler] clear old value")
			p := reflect.ValueOf(handler.aggregate).Elem()
			p.Set(reflect.Zero(p.Type()))
			logrus.Info("[Command Handler] Replaying")
			if err = handler.repository.Replay(handler.aggregate, command.GetAggregateID()); err != nil {
			 	return err
			}
		} else {
			reflect.ValueOf(handler.aggregate).Elem().Set(reflect.ValueOf(instance).Elem())
		}
	}

	// Do command
	logrus.Infof("[Command Handler] handler %v", command.GetType())
	if err = handler.aggregate.HandleCommand(command); err != nil {
		return err
	}

	// Save and publish events
	logrus.Info("[Command Handler] Store aggregate to event store and update to cache")
	if err = handler.repository.SaveAndPublish(handler.aggregate, version, handler.topic); err != nil {
		return err
	}

	return nil
}