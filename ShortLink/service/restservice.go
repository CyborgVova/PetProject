package service

import (
	"context"
	"log"
	"net/http"

	pb "shortlink/proto"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func RunRest() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	err := pb.RegisterLinkBuilderHandlerFromEndpoint(ctx, mux, "localhost:8090", opts)
	if err != nil {
		panic(err)
	}
	log.Printf("HTTP server listening at 8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		panic(err)
	}
}
