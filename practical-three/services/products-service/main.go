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

const serviceName = "products-service"
const servicePort = 50052

type Product struct {
	gorm.Model
	Name  string
	Price float64
}

type server struct {
	pb.UnimplementedProductServiceServer
	db *gorm.DB
}

func (s *server) CreateProduct(ctx context.Context, req *pb.CreateProductRequest) (*pb.ProductResponse, error) {
	product := Product{Name: req.Name, Price: req.Price}
	if result := s.db.Create(&product); result.Error != nil {
		return nil, result.Error
	}
	return &pb.ProductResponse{Product: &pb.Product{Id: fmt.Sprint(product.ID), Name: product.Name, Price: product.Price}}, nil
}

func (s *server) GetProduct(ctx context.Context, req *pb.GetProductRequest) (*pb.ProductResponse, error) {
	var product Product
	if result := s.db.First(&product, req.Id); result.Error != nil {
		return nil, result.Error
	}
	return &pb.ProductResponse{Product: &pb.Product{Id: fmt.Sprint(product.ID), Name: product.Name, Price: product.Price}}, nil
}

func main() {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		getEnv("DB_HOST", "products-db"), getEnv("DB_USER", "user"), getEnv("DB_PASSWORD", "password"), getEnv("DB_NAME", "products_db"), getEnv("DB_PORT", "5432"))

	var db *gorm.DB
	var err error
	for i := 0; i < 30; i++ {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			break
		}
		log.Printf("products-service: DB connect attempt %d failed: %v", i+1, err)
		time.Sleep(2 * time.Second)
	}
	if err != nil {
		log.Fatalf("products-service: could not connect to DB: %v", err)
	}

	db.AutoMigrate(&Product{})

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", servicePort))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterProductServiceServer(s, &server{db: db})

	if err := registerServiceWithConsul(); err != nil {
		log.Printf("products-service: consul registration failed: %v", err)
	}

	log.Printf("products-service listening on %v", lis.Addr())
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
