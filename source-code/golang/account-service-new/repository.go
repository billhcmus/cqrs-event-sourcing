package eventsourcing

// Repository tao Aggregate, luu tru event va publish event
type Repository struct {
	eventStore EventStore
	eventBus EventBus
}

// CreateNewRepository tao mot repository moi
func CreateNewRepository(eventstore EventStore, eventbus EventBus) *Repository {
	return &Repository{
		eventstore,
		eventbus,
	}
}

// Load function tra ve trang thai cuoi cung cua aggregate
func (repo *Repository)Load(aggregate Aggregate, ID string) error {
	events, err := repo.eventStore.Load(ID)

	if err != nil {
		return err
	}

	for _,event := range events {
		aggregate.ApplyChangeHelper(aggregate, event, false) // truong hop nay se khong commit event
	}
	return nil
}

// Save event to event store
func (repo *Repository)Save(aggregate Aggregate, version uint64) error {
	err := repo.eventStore.Save(aggregate.UnCommited(), version)
	if err != nil {
		return err
	}
	return nil
}

// PublishEvents to an eventbus
func (repo *Repository)PublishEvents(aggregate Aggregate, bucket, subset string) error {
	var err error

	for _,event := range aggregate.UnCommited() {
		if err = repo.eventBus.Publish(event, bucket, subset); err != nil {
			return err
		}
	}

	return nil
}