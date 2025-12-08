package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/gorilla/mux"
	consulapi "github.com/hashicorp/consul/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "practicalthree/proto/gen"
)

type ServiceDiscovery struct {
	consul *consulapi.Client
	mu     sync.RWMutex
	conns  map[string]*grpc.ClientConn
}

var sd *ServiceDiscovery

func main() {
	config := consulapi.DefaultConfig()
	if addr := os.Getenv("CONSUL_HTTP_ADDR"); addr != "" {
		config.Address = addr
	}
	consul, err := consulapi.NewClient(config)
	if err != nil {
		log.Fatalf("api-gateway: consul client failed: %v", err)
	}

	sd = &ServiceDiscovery{consul: consul, conns: make(map[string]*grpc.ClientConn)}

	// Wait for services to register with Consul
	time.Sleep(8 * time.Second)

	r := mux.NewRouter()
	r.HandleFunc("/api/users", createUserHandler).Methods("POST")
	r.HandleFunc("/api/users/{id}", getUserHandler).Methods("GET")
	r.HandleFunc("/api/products", createProductHandler).Methods("POST")
	r.HandleFunc("/api/products/{id}", getProductHandler).Methods("GET")
	r.HandleFunc("/api/purchases/user/{userId}/product/{productId}", getPurchaseDataHandler).Methods("GET")

	log.Println("api-gateway: listening on :8080")
	http.ListenAndServe(":8080", r)
}

func (sd *ServiceDiscovery) getServiceConnection(name string) (*grpc.ClientConn, error) {
	sd.mu.Lock()
	defer sd.mu.Unlock()

	if c, ok := sd.conns[name]; ok {
		return c, nil
	}

	// Try to get from Consul catalog
	cat, _, err := sd.consul.Catalog().Service(name, "", nil)
	if err != nil {
		return nil, fmt.Errorf("consul catalog query failed: %v", err)
	}
	if len(cat) == 0 {
		return nil, fmt.Errorf("no instances of %s found in consul", name)
	}

	svc := cat[0]
	addr := fmt.Sprintf("%s:%d", svc.ServiceAddress, svc.ServicePort)
	if svc.ServiceAddress == "" {
		addr = fmt.Sprintf("%s:%d", svc.Address, svc.ServicePort)
	}

	log.Printf("api-gateway: connecting to %s at %s", name, addr)
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("grpc dial failed: %v", err)
	}

	sd.conns[name] = conn
	return conn, nil
}

func getUsersClient() (pb.UserServiceClient, error) {
	conn, err := sd.getServiceConnection("users-service")
	if err != nil {
		return nil, err
	}
	return pb.NewUserServiceClient(conn), nil
}

func getProductsClient() (pb.ProductServiceClient, error) {
	conn, err := sd.getServiceConnection("products-service")
	if err != nil {
		return nil, err
	}
	return pb.NewProductServiceClient(conn), nil
}

// Handlers
func createUserHandler(w http.ResponseWriter, r *http.Request) {
	client, err := getUsersClient()
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}

	var req pb.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, fmt.Sprintf("invalid JSON: %v", err), http.StatusBadRequest)
		return
	}

	res, err := client.CreateUser(context.Background(), &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res.User)
}

func getUserHandler(w http.ResponseWriter, r *http.Request) {
	client, err := getUsersClient()
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}

	vars := mux.Vars(r)
	id := vars["id"]
	res, err := client.GetUser(context.Background(), &pb.GetUserRequest{Id: id})
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res.User)
}

func createProductHandler(w http.ResponseWriter, r *http.Request) {
	client, err := getProductsClient()
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}

	var req pb.CreateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, fmt.Sprintf("invalid JSON: %v", err), http.StatusBadRequest)
		return
	}

	res, err := client.CreateProduct(context.Background(), &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res.Product)
}

func getProductHandler(w http.ResponseWriter, r *http.Request) {
	client, err := getProductsClient()
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}

	vars := mux.Vars(r)
	id := vars["id"]
	res, err := client.GetProduct(context.Background(), &pb.GetProductRequest{Id: id})
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res.Product)
}

func getPurchaseDataHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userId := vars["userId"]
	productId := vars["productId"]

	var wg sync.WaitGroup
	var user *pb.User
	var product *pb.Product
	var userErr, productErr error

	wg.Add(2)

	go func() {
		defer wg.Done()
		client, err := getUsersClient()
		if err != nil {
			userErr = err
			return
		}
		res, err := client.GetUser(context.Background(), &pb.GetUserRequest{Id: userId})
		if err != nil {
			userErr = err
			return
		}
		user = res.User
	}()

	go func() {
		defer wg.Done()
		client, err := getProductsClient()
		if err != nil {
			productErr = err
			return
		}
		res, err := client.GetProduct(context.Background(), &pb.GetProductRequest{Id: productId})
		if err != nil {
			productErr = err
			return
		}
		product = res.Product
	}()

	wg.Wait()

	if userErr != nil {
		http.Error(w, fmt.Sprintf("failed to get user: %v", userErr), http.StatusNotFound)
		return
	}
	if productErr != nil {
		http.Error(w, fmt.Sprintf("failed to get product: %v", productErr), http.StatusNotFound)
		return
	}

	out := map[string]interface{}{"user": user, "product": product}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(out)
}
