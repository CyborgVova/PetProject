package main

import (
	"log"
	"military/api"
	"military/transmitter"
	"net"

	"google.golang.org/grpc"
)

func main() {
	s := grpc.NewServer()
	l, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	transmitter.RegisterTransmitterServer(s, &api.GRPCServer{})
	if err := s.Serve(l); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
