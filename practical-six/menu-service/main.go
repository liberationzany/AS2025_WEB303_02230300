package main

import (
	"log"
	"net"
	"os"

	"github.com/practical6/menu-service/database"
	"github.com/practical6/menu-service/grpc"
	"github.com/practical6/proto/menu/v1"
	grpcServer "google.golang.org/grpc"
)

func main() {
	database.InitDB()

	port := os.Getenv("GRPC_PORT")
	if port == "" {
		port = "50052"
	}

	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpcServer.NewServer()
	menuv1.RegisterMenuServiceServer(s, grpc.NewMenuServer())

	log.Printf("Menu service listening on port %s", port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
