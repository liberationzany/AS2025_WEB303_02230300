package grpc

import (
	"context"

	"github.com/practical6/order-service/database"
	"github.com/practical6/order-service/models"
	menuv1 "github.com/practical6/proto/menu/v1"
	orderv1 "github.com/practical6/proto/order/v1"
	userv1 "github.com/practical6/proto/user/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type OrderServer struct {
	orderv1.UnimplementedOrderServiceServer
	UserClient userv1.UserServiceClient
	MenuClient menuv1.MenuServiceClient
}

func NewOrderServer(userClient userv1.UserServiceClient, menuClient menuv1.MenuServiceClient) *OrderServer {
	return &OrderServer{
		UserClient: userClient,
		MenuClient: menuClient,
	}
}

func (s *OrderServer) CreateOrder(ctx context.Context, req *orderv1.CreateOrderRequest) (*orderv1.CreateOrderResponse, error) {
	// Validate user exists
	_, err := s.UserClient.GetUser(ctx, &userv1.GetUserRequest{Id: req.UserId})
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "user not found: %v", err)
	}

	if len(req.Items) == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "order must have at least one item")
	}

	// Create order
	order := models.Order{
		UserID: uint(req.UserId),
		Status: "pending",
	}

	// Validate menu items and create order items
	var orderItems []models.OrderItem
	for _, item := range req.Items {
		if item.Quantity == 0 {
			return nil, status.Errorf(codes.InvalidArgument, "quantity must be greater than 0")
		}

		// Get menu item details
		menuResp, err := s.MenuClient.GetMenuItem(ctx, &menuv1.GetMenuItemRequest{Id: item.MenuItemId})
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "menu item %d not found: %v", item.MenuItemId, err)
		}

		orderItem := models.OrderItem{
			MenuItemID:   uint(item.MenuItemId),
			MenuItemName: menuResp.MenuItem.Name,
			Quantity:     item.Quantity,
			Price:        menuResp.MenuItem.Price,
		}
		orderItems = append(orderItems, orderItem)
	}

	order.OrderItems = orderItems

	result := database.DB.Create(&order)
	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, "failed to create order: %v", result.Error)
	}

	// Build response
	var protoItems []*orderv1.OrderItem
	for _, item := range order.OrderItems {
		protoItems = append(protoItems, &orderv1.OrderItem{
			Id:           uint32(item.ID),
			MenuItemId:   uint32(item.MenuItemID),
			MenuItemName: item.MenuItemName,
			Quantity:     item.Quantity,
			Price:        item.Price,
		})
	}

	return &orderv1.CreateOrderResponse{
		Order: &orderv1.Order{
			Id:         uint32(order.ID),
			UserId:     uint32(order.UserID),
			Status:     order.Status,
			OrderItems: protoItems,
		},
	}, nil
}

func (s *OrderServer) GetOrder(ctx context.Context, req *orderv1.GetOrderRequest) (*orderv1.GetOrderResponse, error) {
	var order models.Order
	result := database.DB.Preload("OrderItems").First(&order, req.Id)
	if result.Error != nil {
		return nil, status.Errorf(codes.NotFound, "order not found")
	}

	var protoItems []*orderv1.OrderItem
	for _, item := range order.OrderItems {
		protoItems = append(protoItems, &orderv1.OrderItem{
			Id:           uint32(item.ID),
			MenuItemId:   uint32(item.MenuItemID),
			MenuItemName: item.MenuItemName,
			Quantity:     item.Quantity,
			Price:        item.Price,
		})
	}

	return &orderv1.GetOrderResponse{
		Order: &orderv1.Order{
			Id:         uint32(order.ID),
			UserId:     uint32(order.UserID),
			Status:     order.Status,
			OrderItems: protoItems,
		},
	}, nil
}

func (s *OrderServer) GetOrders(ctx context.Context, req *orderv1.GetOrdersRequest) (*orderv1.GetOrdersResponse, error) {
	var orders []models.Order
	result := database.DB.Preload("OrderItems").Find(&orders)
	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, "failed to fetch orders: %v", result.Error)
	}

	var protoOrders []*orderv1.Order
	for _, order := range orders {
		var protoItems []*orderv1.OrderItem
		for _, item := range order.OrderItems {
			protoItems = append(protoItems, &orderv1.OrderItem{
				Id:           uint32(item.ID),
				MenuItemId:   uint32(item.MenuItemID),
				MenuItemName: item.MenuItemName,
				Quantity:     item.Quantity,
				Price:        item.Price,
			})
		}

		protoOrders = append(protoOrders, &orderv1.Order{
			Id:         uint32(order.ID),
			UserId:     uint32(order.UserID),
			Status:     order.Status,
			OrderItems: protoItems,
		})
	}

	return &orderv1.GetOrdersResponse{
		Orders: protoOrders,
	}, nil
}
