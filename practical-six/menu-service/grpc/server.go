package grpc

import (
	"context"

	"github.com/practical6/menu-service/database"
	"github.com/practical6/menu-service/models"
	"github.com/practical6/proto/menu/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type MenuServer struct {
	menuv1.UnimplementedMenuServiceServer
}

func NewMenuServer() *MenuServer {
	return &MenuServer{}
}

func (s *MenuServer) CreateMenuItem(ctx context.Context, req *menuv1.CreateMenuItemRequest) (*menuv1.CreateMenuItemResponse, error) {
	menuItem := models.MenuItem{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
	}

	result := database.DB.Create(&menuItem)
	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, "failed to create menu item: %v", result.Error)
	}

	return &menuv1.CreateMenuItemResponse{
		MenuItem: &menuv1.MenuItem{
			Id:          uint32(menuItem.ID),
			Name:        menuItem.Name,
			Description: menuItem.Description,
			Price:       menuItem.Price,
		},
	}, nil
}

func (s *MenuServer) GetMenuItem(ctx context.Context, req *menuv1.GetMenuItemRequest) (*menuv1.GetMenuItemResponse, error) {
	var menuItem models.MenuItem
	result := database.DB.First(&menuItem, req.Id)
	if result.Error != nil {
		return nil, status.Errorf(codes.NotFound, "menu item not found")
	}

	return &menuv1.GetMenuItemResponse{
		MenuItem: &menuv1.MenuItem{
			Id:          uint32(menuItem.ID),
			Name:        menuItem.Name,
			Description: menuItem.Description,
			Price:       menuItem.Price,
		},
	}, nil
}

func (s *MenuServer) GetMenuItems(ctx context.Context, req *menuv1.GetMenuItemsRequest) (*menuv1.GetMenuItemsResponse, error) {
	var menuItems []models.MenuItem
	result := database.DB.Find(&menuItems)
	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, "failed to fetch menu items: %v", result.Error)
	}

	var protoItems []*menuv1.MenuItem
	for _, item := range menuItems {
		protoItems = append(protoItems, &menuv1.MenuItem{
			Id:          uint32(item.ID),
			Name:        item.Name,
			Description: item.Description,
			Price:       item.Price,
		})
	}

	return &menuv1.GetMenuItemsResponse{
		MenuItems: protoItems,
	}, nil
}
