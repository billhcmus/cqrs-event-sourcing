package asynccommandbus

import (
	"github.com/google/uuid"
	"log"
	"eventsourcing"
)

var workerJobPool = make(chan chan eventsourcing.Command)

// Worker struct contain information to handle command
type Worker struct {
	ID string
	WorkerChannel chan chan eventsourcing.Command
	JobChannel chan eventsourcing.Command
	CommandRegister eventsourcing.CommandHandlerRegister
}

// Start a worker
func (w *Worker) Start() {
	go func() {
		for {
			// Khi khoi dong 1 worker, worker se dua channel de chua job vao worker pool (noi giao tiep giua worker va bus)
			w.WorkerChannel <- w.JobChannel
			job := <- w.JobChannel // khi bus dua job vao thi ta lay ra
	
			handler, err := w.CommandRegister.Get(job)
			log.Printf("%v handle %v", w.ID, job)

			if err != nil {
				continue
			}
			if err := handler.Handle(job); err != nil {
				log.Fatal(err)
			}
		}
	}()
}

// NewWorker create worker
func NewWorker(reg eventsourcing.CommandHandlerRegister) {
	uuidv4, _ := uuid.NewRandom()
	worker := Worker {
		ID: uuidv4.String(),
		WorkerChannel: workerJobPool,
		JobChannel: make(chan eventsourcing.Command), 
		CommandRegister: reg,
	}
	worker.Start()
}


// Bus is struct contain information
type Bus struct {
	maxWorker int
	commandHandlerRegister eventsourcing.CommandHandlerRegister
}

// CreateBus create new command bus
func CreateBus(numWorker int, reg eventsourcing.CommandHandlerRegister) *Bus {
	return &Bus{maxWorker: numWorker, commandHandlerRegister: reg}
}

// HandleCommand is definition of CommandHandler interface
func (bus *Bus) HandleCommand(command eventsourcing.Command) {
	go func(command eventsourcing.Command) {
		JobChannel := <- workerJobPool
		JobChannel <- command
	}(command)
}

// StartBus start command bus
func (bus *Bus)StartBus() {
	for i := 0; i < bus.maxWorker; i++ {
		NewWorker(bus.commandHandlerRegister)
	}
}