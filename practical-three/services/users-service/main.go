package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	consulapi "github.com/hashicorp/consul/api"
	"google.golang.org/grpc"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	pb "practicalthree/proto/gen"
)

const serviceName = "users-service"
const servicePort = 50051

type User struct {
	gorm.Model
	Name  string
	Email string `gorm:"unique"`
}

type server struct {
	pb.UnimplementedUserServiceServer
	db *gorm.DB
}

func (s *server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.UserResponse, error) {
	user := User{Name: req.Name, Email: req.Email}
	if result := s.db.Create(&user); result.Error != nil {
		return nil, result.Error
	}
	return &pb.UserResponse{User: &pb.User{Id: fmt.Sprint(user.ID), Name: user.Name, Email: user.Email}}, nil
}

func (s *server) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.UserResponse, error) {
	var user User
	if result := s.db.First(&user, req.Id); result.Error != nil {
		return nil, result.Error
	}
	return &pb.UserResponse{User: &pb.User{Id: fmt.Sprint(user.ID), Name: user.Name, Email: user.Email}}, nil
}

func main() {
	// Retry DB connect
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		getEnv("DB_HOST", "users-db"), getEnv("DB_USER", "user"), getEnv("DB_PASSWORD", "password"), getEnv("DB_NAME", "users_db"), getEnv("DB_PORT", "5432"))

	var db *gorm.DB
	var err error
	for i := 0; i < 30; i++ {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			break
		}
		log.Printf("users-service: DB connect attempt %d failed: %v", i+1, err)
		time.Sleep(2 * time.Second)
	}
	if err != nil {
		log.Fatalf("users-service: could not connect to DB: %v", err)
	}

	db.AutoMigrate(&User{})

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", servicePort))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterUserServiceServer(s, &server{db: db})

	if err := registerServiceWithConsul(); err != nil {
		log.Printf("users-service: consul registration failed: %v", err)
	}

	log.Printf("users-service listening on %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

func registerServiceWithConsul() error {
	config := consulapi.DefaultConfig()
	if addr := os.Getenv("CONSUL_HTTP_ADDR"); addr != "" {
		config.Address = addr
	}
	consul, err := consulapi.NewClient(config)
	if err != nil {
		return err
	}

	serviceAddr := os.Getenv("SERVICE_ADDR")
	if serviceAddr == "" {
		serviceAddr = serviceName
	}

	hostname, _ := os.Hostname()

	reg := &consulapi.AgentServiceRegistration{
		ID:      fmt.Sprintf("%s-%s", serviceName, hostname),
		Name:    serviceName,
		Port:    servicePort,
		Address: serviceAddr,
	}
	return consul.Agent().ServiceRegister(reg)
}

func getEnv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
