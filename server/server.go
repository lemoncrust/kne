package main

import (
	"log"
	"net"

	"google.golang.org/grpc"
)

type Server struct{}

const (
	port = ":50051"
)

type server struct {
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
