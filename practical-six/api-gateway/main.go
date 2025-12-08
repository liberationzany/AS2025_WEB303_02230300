package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	menuv1 "github.com/practical6/proto/menu/v1"
	orderv1 "github.com/practical6/proto/order/v1"
	userv1 "github.com/practical6/proto/user/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	userClient  userv1.UserServiceClient
	menuClient  menuv1.MenuServiceClient
	orderClient orderv1.OrderServiceClient
)

func main() {
	// Connect to user service
	userConn, err := grpc.Dial(
		getEnv("USER_SERVICE_ADDR", "localhost:50051"),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("Failed to connect to user service: %v", err)
	}
	defer userConn.Close()
	userClient = userv1.NewUserServiceClient(userConn)

	// Connect to menu service
	menuConn, err := grpc.Dial(
		getEnv("MENU_SERVICE_ADDR", "localhost:50052"),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("Failed to connect to menu service: %v", err)
	}
	defer menuConn.Close()
	menuClient = menuv1.NewMenuServiceClient(menuConn)

	// Connect to order service
	orderConn, err := grpc.Dial(
		getEnv("ORDER_SERVICE_ADDR", "localhost:50053"),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("Failed to connect to order service: %v", err)
	}
	defer orderConn.Close()
	orderClient = orderv1.NewOrderServiceClient(orderConn)

	router := mux.NewRouter()

	// User endpoints
	router.HandleFunc("/api/users", createUserHandler).Methods("POST")
	router.HandleFunc("/api/users/{id}", getUserHandler).Methods("GET")
	router.HandleFunc("/api/users", getUsersHandler).Methods("GET")

	// Menu endpoints
	router.HandleFunc("/api/menu", createMenuItemHandler).Methods("POST")
	router.HandleFunc("/api/menu/{id}", getMenuItemHandler).Methods("GET")
	router.HandleFunc("/api/menu", getMenuItemsHandler).Methods("GET")

	// Order endpoints
	router.HandleFunc("/api/orders", createOrderHandler).Methods("POST")
	router.HandleFunc("/api/orders/{id}", getOrderHandler).Methods("GET")
	router.HandleFunc("/api/orders", getOrdersHandler).Methods("GET")

	port := getEnv("PORT", "8080")
	log.Printf("API Gateway listening on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}

// User handlers
func createUserHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name        string `json:"name"`
		Email       string `json:"email"`
		IsCafeOwner bool   `json:"is_cafe_owner"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp, err := userClient.CreateUser(r.Context(), &userv1.CreateUserRequest{
		Name:        req.Name,
		Email:       req.Email,
		IsCafeOwner: req.IsCafeOwner,
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":            resp.User.Id,
		"name":          resp.User.Name,
		"email":         resp.User.Email,
		"is_cafe_owner": resp.User.IsCafeOwner,
	})
}

func getUserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	resp, err := userClient.GetUser(r.Context(), &userv1.GetUserRequest{Id: uint32(id)})
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":            resp.User.Id,
		"name":          resp.User.Name,
		"email":         resp.User.Email,
		"is_cafe_owner": resp.User.IsCafeOwner,
	})
}

func getUsersHandler(w http.ResponseWriter, r *http.Request) {
	resp, err := userClient.GetUsers(r.Context(), &userv1.GetUsersRequest{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var users []map[string]interface{}
	for _, user := range resp.Users {
		users = append(users, map[string]interface{}{
			"id":            user.Id,
			"name":          user.Name,
			"email":         user.Email,
			"is_cafe_owner": user.IsCafeOwner,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

// Menu handlers
func createMenuItemHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name        string  `json:"name"`
		Description string  `json:"description"`
		Price       float64 `json:"price"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp, err := menuClient.CreateMenuItem(r.Context(), &menuv1.CreateMenuItemRequest{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":          resp.MenuItem.Id,
		"name":        resp.MenuItem.Name,
		"description": resp.MenuItem.Description,
		"price":       resp.MenuItem.Price,
	})
}

func getMenuItemHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	resp, err := menuClient.GetMenuItem(r.Context(), &menuv1.GetMenuItemRequest{Id: uint32(id)})
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":          resp.MenuItem.Id,
		"name":        resp.MenuItem.Name,
		"description": resp.MenuItem.Description,
		"price":       resp.MenuItem.Price,
	})
}

func getMenuItemsHandler(w http.ResponseWriter, r *http.Request) {
	resp, err := menuClient.GetMenuItems(r.Context(), &menuv1.GetMenuItemsRequest{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var items []map[string]interface{}
	for _, item := range resp.MenuItems {
		items = append(items, map[string]interface{}{
			"id":          item.Id,
			"name":        item.Name,
			"description": item.Description,
			"price":       item.Price,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}

// Order handlers
func createOrderHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		UserID uint32 `json:"user_id"`
		Items  []struct {
			MenuItemID uint32 `json:"menu_item_id"`
			Quantity   uint32 `json:"quantity"`
		} `json:"items"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var items []*orderv1.OrderItemRequest
	for _, item := range req.Items {
		items = append(items, &orderv1.OrderItemRequest{
			MenuItemId: item.MenuItemID,
			Quantity:   item.Quantity,
		})
	}

	resp, err := orderClient.CreateOrder(r.Context(), &orderv1.CreateOrderRequest{
		UserId: req.UserID,
		Items:  items,
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var orderItems []map[string]interface{}
	for _, item := range resp.Order.OrderItems {
		orderItems = append(orderItems, map[string]interface{}{
			"id":             item.Id,
			"menu_item_id":   item.MenuItemId,
			"menu_item_name": item.MenuItemName,
			"quantity":       item.Quantity,
			"price":          item.Price,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":          resp.Order.Id,
		"user_id":     resp.Order.UserId,
		"status":      resp.Order.Status,
		"order_items": orderItems,
	})
}

func getOrderHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	resp, err := orderClient.GetOrder(r.Context(), &orderv1.GetOrderRequest{Id: uint32(id)})
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	var orderItems []map[string]interface{}
	for _, item := range resp.Order.OrderItems {
		orderItems = append(orderItems, map[string]interface{}{
			"id":             item.Id,
			"menu_item_id":   item.MenuItemId,
			"menu_item_name": item.MenuItemName,
			"quantity":       item.Quantity,
			"price":          item.Price,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":          resp.Order.Id,
		"user_id":     resp.Order.UserId,
		"status":      resp.Order.Status,
		"order_items": orderItems,
	})
}

func getOrdersHandler(w http.ResponseWriter, r *http.Request) {
	resp, err := orderClient.GetOrders(r.Context(), &orderv1.GetOrdersRequest{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var orders []map[string]interface{}
	for _, order := range resp.Orders {
		var orderItems []map[string]interface{}
		for _, item := range order.OrderItems {
			orderItems = append(orderItems, map[string]interface{}{
				"id":             item.Id,
				"menu_item_id":   item.MenuItemId,
				"menu_item_name": item.MenuItemName,
				"quantity":       item.Quantity,
				"price":          item.Price,
			})
		}

		orders = append(orders, map[string]interface{}{
			"id":          order.Id,
			"user_id":     order.UserId,
			"status":      order.Status,
			"order_items": orderItems,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orders)
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
