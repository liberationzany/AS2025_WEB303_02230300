package integration

import (
	"context"
	"log"
	"net"
	"testing"

	menudatabase "github.com/practical6/menu-service/database"
	menugrpc "github.com/practical6/menu-service/grpc"
	menumodels "github.com/practical6/menu-service/models"
	orderdatabase "github.com/practical6/order-service/database"
	ordergrpc "github.com/practical6/order-service/grpc"
	ordermodels "github.com/practical6/order-service/models"
	menuv1 "github.com/practical6/proto/menu/v1"
	orderv1 "github.com/practical6/proto/order/v1"
	userv1 "github.com/practical6/proto/user/v1"
	userdatabase "github.com/practical6/user-service/database"
	usergrpc "github.com/practical6/user-service/grpc"
	usermodels "github.com/practical6/user-service/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const bufSize = 1024 * 1024

var (
	userListener  *bufconn.Listener
	menuListener  *bufconn.Listener
	orderListener *bufconn.Listener
)

func setupUserService(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	require.NoError(t, err)

	err = db.AutoMigrate(&usermodels.User{})
	require.NoError(t, err)

	userdatabase.DB = db

	userListener = bufconn.Listen(bufSize)
	s := grpc.NewServer()
	userv1.RegisterUserServiceServer(s, usergrpc.NewUserServer())

	go func() {
		if err := s.Serve(userListener); err != nil {
			log.Fatalf("Server exited: %v", err)
		}
	}()
}

func setupMenuService(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	require.NoError(t, err)

	err = db.AutoMigrate(&menumodels.MenuItem{})
	require.NoError(t, err)

	menudatabase.DB = db

	menuListener = bufconn.Listen(bufSize)
	s := grpc.NewServer()
	menuv1.RegisterMenuServiceServer(s, menugrpc.NewMenuServer())

	go func() {
		if err := s.Serve(menuListener); err != nil {
			log.Fatalf("Server exited: %v", err)
		}
	}()
}

func setupOrderService(t *testing.T, userConn, menuConn *grpc.ClientConn) {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	require.NoError(t, err)

	err = db.AutoMigrate(&ordermodels.Order{}, &ordermodels.OrderItem{})
	require.NoError(t, err)

	orderdatabase.DB = db

	userClient := userv1.NewUserServiceClient(userConn)
	menuClient := menuv1.NewMenuServiceClient(menuConn)

	orderListener = bufconn.Listen(bufSize)
	s := grpc.NewServer()
	orderv1.RegisterOrderServiceServer(s, ordergrpc.NewOrderServer(userClient, menuClient))

	go func() {
		if err := s.Serve(orderListener); err != nil {
			log.Fatalf("Server exited: %v", err)
		}
	}()
}

func bufDialer(listener *bufconn.Listener) func(context.Context, string) (net.Conn, error) {
	return func(ctx context.Context, url string) (net.Conn, error) {
		return listener.Dial()
	}
}

func TestIntegration_CreateUser(t *testing.T) {
	setupUserService(t)
	defer userListener.Close()

	ctx := context.Background()

	conn, err := grpc.DialContext(ctx, "bufnet",
		grpc.WithContextDialer(bufDialer(userListener)),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)
	defer conn.Close()

	client := userv1.NewUserServiceClient(conn)

	resp, err := client.CreateUser(ctx, &userv1.CreateUserRequest{
		Name:        "Integration User",
		Email:       "integration@test.com",
		IsCafeOwner: false,
	})

	require.NoError(t, err)
	assert.NotZero(t, resp.User.Id)
	assert.Equal(t, "Integration User", resp.User.Name)
}

func TestIntegration_CompleteOrderFlow(t *testing.T) {
	setupUserService(t)
	defer userListener.Close()

	setupMenuService(t)
	defer menuListener.Close()

	ctx := context.Background()

	// Connect to user service
	userConn, err := grpc.DialContext(ctx, "bufnet",
		grpc.WithContextDialer(bufDialer(userListener)),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)
	defer userConn.Close()

	userClient := userv1.NewUserServiceClient(userConn)

	// Connect to menu service
	menuConn, err := grpc.DialContext(ctx, "bufnet",
		grpc.WithContextDialer(bufDialer(menuListener)),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)
	defer menuConn.Close()

	menuClient := menuv1.NewMenuServiceClient(menuConn)

	// Setup order service with connections
	setupOrderService(t, userConn, menuConn)
	defer orderListener.Close()

	orderConn, err := grpc.DialContext(ctx, "bufnet",
		grpc.WithContextDialer(bufDialer(orderListener)),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)
	defer orderConn.Close()

	orderClient := orderv1.NewOrderServiceClient(orderConn)

	// Step 1: Create a user
	userResp, err := userClient.CreateUser(ctx, &userv1.CreateUserRequest{
		Name:        "Integration User",
		Email:       "integration@test.com",
		IsCafeOwner: false,
	})
	require.NoError(t, err)
	userID := userResp.User.Id

	// Step 2: Create menu items
	item1, err := menuClient.CreateMenuItem(ctx, &menuv1.CreateMenuItemRequest{
		Name:        "Coffee",
		Description: "Hot coffee",
		Price:       2.50,
	})
	require.NoError(t, err)

	item2, err := menuClient.CreateMenuItem(ctx, &menuv1.CreateMenuItemRequest{
		Name:        "Sandwich",
		Description: "Ham sandwich",
		Price:       5.00,
	})
	require.NoError(t, err)

	// Step 3: Create an order
	orderResp, err := orderClient.CreateOrder(ctx, &orderv1.CreateOrderRequest{
		UserId: userID,
		Items: []*orderv1.OrderItemRequest{
			{MenuItemId: item1.MenuItem.Id, Quantity: 2},
			{MenuItemId: item2.MenuItem.Id, Quantity: 1},
		},
	})

	require.NoError(t, err)
	assert.NotZero(t, orderResp.Order.Id)
	assert.Equal(t, userID, orderResp.Order.UserId)
	assert.Equal(t, "pending", orderResp.Order.Status)
	assert.Len(t, orderResp.Order.OrderItems, 2)

	// Verify prices were snapshotted
	assert.InDelta(t, 2.50, orderResp.Order.OrderItems[0].Price, 0.001)
	assert.InDelta(t, 5.00, orderResp.Order.OrderItems[1].Price, 0.001)

	// Step 4: Retrieve the order
	getOrderResp, err := orderClient.GetOrder(ctx, &orderv1.GetOrderRequest{
		Id: orderResp.Order.Id,
	})

	require.NoError(t, err)
	assert.Equal(t, orderResp.Order.Id, getOrderResp.Order.Id)
	assert.Len(t, getOrderResp.Order.OrderItems, 2)
}

func TestIntegration_OrderValidation(t *testing.T) {
	setupUserService(t)
	defer userListener.Close()

	setupMenuService(t)
	defer menuListener.Close()

	ctx := context.Background()

	userConn, err := grpc.DialContext(ctx, "bufnet",
		grpc.WithContextDialer(bufDialer(userListener)),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)
	defer userConn.Close()

	userClient := userv1.NewUserServiceClient(userConn)

	menuConn, err := grpc.DialContext(ctx, "bufnet",
		grpc.WithContextDialer(bufDialer(menuListener)),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)
	defer menuConn.Close()

	setupOrderService(t, userConn, menuConn)
	defer orderListener.Close()

	orderConn, err := grpc.DialContext(ctx, "bufnet",
		grpc.WithContextDialer(bufDialer(orderListener)),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)
	defer orderConn.Close()

	orderClient := orderv1.NewOrderServiceClient(orderConn)

	// Try to create order with invalid user
	t.Run("invalid user", func(t *testing.T) {
		_, err := orderClient.CreateOrder(ctx, &orderv1.CreateOrderRequest{
			UserId: 9999,
			Items: []*orderv1.OrderItemRequest{
				{MenuItemId: 1, Quantity: 1},
			},
		})

		require.Error(t, err)
		assert.Contains(t, err.Error(), "user not found")
	})

	// Create valid user
	userResp, err := userClient.CreateUser(ctx, &userv1.CreateUserRequest{
		Name: "Valid User", Email: "valid@test.com",
	})
	require.NoError(t, err)

	// Try to create order with invalid menu item
	t.Run("invalid menu item", func(t *testing.T) {
		_, err := orderClient.CreateOrder(ctx, &orderv1.CreateOrderRequest{
			UserId: userResp.User.Id,
			Items: []*orderv1.OrderItemRequest{
				{MenuItemId: 9999, Quantity: 1},
			},
		})

		require.Error(t, err)
		assert.Contains(t, err.Error(), "menu item 9999 not found")
	})
}

func TestIntegration_ConcurrentOrders(t *testing.T) {
	setupUserService(t)
	defer userListener.Close()

	setupMenuService(t)
	defer menuListener.Close()

	ctx := context.Background()

	userConn, err := grpc.DialContext(ctx, "bufnet",
		grpc.WithContextDialer(bufDialer(userListener)),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)
	defer userConn.Close()

	userClient := userv1.NewUserServiceClient(userConn)

	menuConn, err := grpc.DialContext(ctx, "bufnet",
		grpc.WithContextDialer(bufDialer(menuListener)),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)
	defer menuConn.Close()

	menuClient := menuv1.NewMenuServiceClient(menuConn)

	setupOrderService(t, userConn, menuConn)
	defer orderListener.Close()

	orderConn, err := grpc.DialContext(ctx, "bufnet",
		grpc.WithContextDialer(bufDialer(orderListener)),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)
	defer orderConn.Close()

	orderClient := orderv1.NewOrderServiceClient(orderConn)

	// Create test user and menu item
	userResp, err := userClient.CreateUser(ctx, &userv1.CreateUserRequest{
		Name: "Test User", Email: "test@example.com",
	})
	require.NoError(t, err)
	userID := userResp.User.Id

	itemResp, err := menuClient.CreateMenuItem(ctx, &menuv1.CreateMenuItemRequest{
		Name: "Test Item", Price: 10.00,
	})
	require.NoError(t, err)
	itemID := itemResp.MenuItem.Id

	// Create multiple orders concurrently
	numOrders := 10
	errChan := make(chan error, numOrders)
	respChan := make(chan *orderv1.CreateOrderResponse, numOrders)

	for i := 0; i < numOrders; i++ {
		go func() {
			resp, err := orderClient.CreateOrder(ctx, &orderv1.CreateOrderRequest{
				UserId: userID,
				Items: []*orderv1.OrderItemRequest{
					{MenuItemId: itemID, Quantity: 1},
				},
			})
			errChan <- err
			respChan <- resp
		}()
	}

	// Collect results
	for i := 0; i < numOrders; i++ {
		err := <-errChan
		resp := <-respChan
		require.NoError(t, err)
		assert.NotZero(t, resp.Order.Id)
	}

	// Verify all orders were created
	ordersResp, err := orderClient.GetOrders(ctx, &orderv1.GetOrdersRequest{})
	require.NoError(t, err)
	assert.Len(t, ordersResp.Orders, numOrders)
}
