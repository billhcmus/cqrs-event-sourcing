package main

import (
	"eventstore"
	"eventstore/storage"
	"eventstore/proto"
	"google.golang.org/grpc"
	"net"
	"log"
	"context"
)

type server struct {
	eventstore.EventStore
}

func (s *server) Save(ctx context.Context, in *pb.SaveEventsRequest) (*pb.SaveEventsResponse, error) {
	err := s.EventStore.Save(in.Events, in.Version)
	if err != nil {
		return &pb.SaveEventsResponse {
			Status: &pb.Status {
				Code: -1,
				Error: err.Error(),
			},
		}, err
	}
	return &pb.SaveEventsResponse {
		Status: &pb.Status{Code : 0, Error : ""},
	}, nil
}

func (s *server) Load(ctx context.Context, in *pb.LoadEventsRequest) (*pb.LoadEventsResponse, error) {
	var events []*pb.BaseEvent
	var err error
	events, err = s.EventStore.Load(in.AggregateId)
	if events == nil {
		return &pb.LoadEventsResponse {
			Events: events,
			Status: &pb.Status {Code: -1, Error: "Not found"},
		}, err
	}
	return &pb.LoadEventsResponse {
		Events: events,
		Status: &pb.Status {Code: 0, Error: ""},
	}, err
}

const port = ":7777"

func main() {
	var dbURL = "postgres://postgres:1234@localhost/eventstore?sslmode=disable"

	postgresClient, err := storage.CreateClient(dbURL, true)

	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("Config database completed")
	}
	log.Printf("event store listen on %v", port)
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("event store fail to listen %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterEventStoreServer(s, &server{postgresClient})
	s.Serve(lis)
}