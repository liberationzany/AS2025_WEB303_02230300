package e2e

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var apiGatewayURL = getEnv("API_GATEWAY_URL", "http://localhost:8080")

type User struct {
	ID          uint32 `json:"id"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	IsCafeOwner bool   `json:"is_cafe_owner"`
}

type MenuItem struct {
	ID          uint32  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}

type OrderItem struct {
	ID           uint32  `json:"id"`
	MenuItemID   uint32  `json:"menu_item_id"`
	MenuItemName string  `json:"menu_item_name"`
	Quantity     uint32  `json:"quantity"`
	Price        float64 `json:"price"`
}

type Order struct {
	ID         uint32      `json:"id"`
	UserID     uint32      `json:"user_id"`
	Status     string      `json:"status"`
	OrderItems []OrderItem `json:"order_items"`
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func makeRequest(method, path string, body interface{}) (*http.Response, error) {
	var bodyReader io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		bodyReader = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequest(method, apiGatewayURL+path, bodyReader)
	if err != nil {
		return nil, err
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	client := &http.Client{Timeout: 10 * time.Second}
	return client.Do(req)
}

func TestMain(m *testing.M) {
	// Wait for services to be ready
	fmt.Println("Waiting for services...")
	maxRetries := 30
	for i := 0; i < maxRetries; i++ {
		resp, err := http.Get(apiGatewayURL + "/api/users")
		if err == nil && resp.StatusCode < 500 {
			resp.Body.Close()
			fmt.Println("Services ready!")
			break
		}
		if i == maxRetries-1 {
			fmt.Println("Services not ready")
			os.Exit(1)
		}
		time.Sleep(2 * time.Second)
	}

	code := m.Run()
	os.Exit(code)
}

func TestE2E_CompleteOrderFlow(t *testing.T) {
	// Step 1: Create a user
	userReq := map[string]interface{}{
		"name":          "E2E User",
		"email":         fmt.Sprintf("e2e-%d@test.com", time.Now().Unix()),
		"is_cafe_owner": false,
	}

	userResp, err := makeRequest("POST", "/api/users", userReq)
	require.NoError(t, err)
	defer userResp.Body.Close()

	assert.Equal(t, http.StatusCreated, userResp.StatusCode)

	var user User
	err = json.NewDecoder(userResp.Body).Decode(&user)
	require.NoError(t, err)

	// Step 2: Create menu items
	item1Req := map[string]interface{}{
		"name":        "Coffee",
		"description": "Hot coffee",
		"price":       2.50,
	}

	item1Resp, err := makeRequest("POST", "/api/menu", item1Req)
	require.NoError(t, err)
	defer item1Resp.Body.Close()

	var item1 MenuItem
	err = json.NewDecoder(item1Resp.Body).Decode(&item1)
	require.NoError(t, err)

	item2Req := map[string]interface{}{
		"name":        "Sandwich",
		"description": "Ham sandwich",
		"price":       5.00,
	}

	item2Resp, err := makeRequest("POST", "/api/menu", item2Req)
	require.NoError(t, err)
	defer item2Resp.Body.Close()

	var item2 MenuItem
	err = json.NewDecoder(item2Resp.Body).Decode(&item2)
	require.NoError(t, err)

	// Step 3: Create order
	orderReq := map[string]interface{}{
		"user_id": user.ID,
		"items": []map[string]interface{}{
			{"menu_item_id": item1.ID, "quantity": 2},
			{"menu_item_id": item2.ID, "quantity": 1},
		},
	}

	orderResp, err := makeRequest("POST", "/api/orders", orderReq)
	require.NoError(t, err)
	defer orderResp.Body.Close()

	assert.Equal(t, http.StatusCreated, orderResp.StatusCode)

	var order Order
	err = json.NewDecoder(orderResp.Body).Decode(&order)
	require.NoError(t, err)

	assert.NotZero(t, order.ID)
	assert.Equal(t, user.ID, order.UserID)
	assert.Len(t, order.OrderItems, 2)

	// Step 4: Retrieve order
	getOrderResp, err := makeRequest("GET", fmt.Sprintf("/api/orders/%d", order.ID), nil)
	require.NoError(t, err)
	defer getOrderResp.Body.Close()

	assert.Equal(t, http.StatusOK, getOrderResp.StatusCode)

	var retrievedOrder Order
	err = json.NewDecoder(getOrderResp.Body).Decode(&retrievedOrder)
	require.NoError(t, err)

	assert.Equal(t, order.ID, retrievedOrder.ID)
	assert.Len(t, retrievedOrder.OrderItems, 2)
}

func TestE2E_OrderValidation(t *testing.T) {
	t.Run("invalid user", func(t *testing.T) {
		orderReq := map[string]interface{}{
			"user_id": 999999,
			"items": []map[string]interface{}{
				{"menu_item_id": 1, "quantity": 1},
			},
		}

		resp, err := makeRequest("POST", "/api/orders", orderReq)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("invalid menu item", func(t *testing.T) {
		// Create a user first
		userReq := map[string]interface{}{
			"name":          "Test User",
			"email":         fmt.Sprintf("test-%d@test.com", time.Now().Unix()),
			"is_cafe_owner": false,
		}

		userResp, err := makeRequest("POST", "/api/users", userReq)
		require.NoError(t, err)
		defer userResp.Body.Close()

		var user User
		err = json.NewDecoder(userResp.Body).Decode(&user)
		require.NoError(t, err)

		// Try to create order with invalid menu item
		orderReq := map[string]interface{}{
			"user_id": user.ID,
			"items": []map[string]interface{}{
				{"menu_item_id": 999999, "quantity": 1},
			},
		}

		resp, err := makeRequest("POST", "/api/orders", orderReq)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}

func TestE2E_GetAllUsers(t *testing.T) {
	// Create a user
	userReq := map[string]interface{}{
		"name":          "List Test User",
		"email":         fmt.Sprintf("listtest-%d@test.com", time.Now().Unix()),
		"is_cafe_owner": false,
	}

	_, err := makeRequest("POST", "/api/users", userReq)
	require.NoError(t, err)

	// Get all users
	resp, err := makeRequest("GET", "/api/users", nil)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var users []User
	err = json.NewDecoder(resp.Body).Decode(&users)
	require.NoError(t, err)

	assert.NotEmpty(t, users)
}

func TestE2E_GetAllMenuItems(t *testing.T) {
	// Create a menu item
	itemReq := map[string]interface{}{
		"name":        "Test Item",
		"description": "Test description",
		"price":       3.50,
	}

	_, err := makeRequest("POST", "/api/menu", itemReq)
	require.NoError(t, err)

	// Get all menu items
	resp, err := makeRequest("GET", "/api/menu", nil)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var items []MenuItem
	err = json.NewDecoder(resp.Body).Decode(&items)
	require.NoError(t, err)

	assert.NotEmpty(t, items)
}
