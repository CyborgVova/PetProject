package service

import (
	"log"
	"net"

	"google.golang.org/grpc"
	"gorm.io/gorm"

	"shortlink/internal"
	pb "shortlink/proto"
)

func RunGrpc(storage string, db *gorm.DB) {
	lis, err := net.Listen("tcp", ":8090")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterLinkBuilderServer(s, &internal.Server{Storage: storage, DB: db})
	log.Printf("GRPC server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		panic(err)
	}
}
