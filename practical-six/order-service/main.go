package main

import (
	"log"
	"net"
	"os"

	"github.com/practical6/order-service/database"
	"github.com/practical6/order-service/grpc"
	menuv1 "github.com/practical6/proto/menu/v1"
	orderv1 "github.com/practical6/proto/order/v1"
	userv1 "github.com/practical6/proto/user/v1"
	grpcClient "google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	grpcServer "google.golang.org/grpc"
)

func main() {
	database.InitDB()

	// Connect to user service
	userConn, err := grpcClient.Dial(
		getEnv("USER_SERVICE_ADDR", "localhost:50051"),
		grpcClient.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("Failed to connect to user service: %v", err)
	}
	defer userConn.Close()
	userClient := userv1.NewUserServiceClient(userConn)

	// Connect to menu service
	menuConn, err := grpcClient.Dial(
		getEnv("MENU_SERVICE_ADDR", "localhost:50052"),
		grpcClient.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("Failed to connect to menu service: %v", err)
	}
	defer menuConn.Close()
	menuClient := menuv1.NewMenuServiceClient(menuConn)

	port := getEnv("GRPC_PORT", "50053")
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpcServer.NewServer()
	orderv1.RegisterOrderServiceServer(s, grpc.NewOrderServer(userClient, menuClient))

	log.Printf("Order service listening on port %s", port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
