package main

import (
	"log"
	"net"
	"os"

	"github.com/practical6/proto/user/v1"
	"github.com/practical6/user-service/database"
	"github.com/practical6/user-service/grpc"
	grpcServer "google.golang.org/grpc"
)

func main() {
	database.InitDB()

	port := os.Getenv("GRPC_PORT")
	if port == "" {
		port = "50051"
	}

	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpcServer.NewServer()
	userv1.RegisterUserServiceServer(s, grpc.NewUserServer())

	log.Printf("User service listening on port %s", port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
