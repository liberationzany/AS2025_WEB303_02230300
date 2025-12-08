package grpc

import (
	"context"

	"github.com/practical6/proto/user/v1"
	"github.com/practical6/user-service/database"
	"github.com/practical6/user-service/models"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserServer struct {
	userv1.UnimplementedUserServiceServer
}

func NewUserServer() *UserServer {
	return &UserServer{}
}

func (s *UserServer) CreateUser(ctx context.Context, req *userv1.CreateUserRequest) (*userv1.CreateUserResponse, error) {
	user := models.User{
		Name:        req.Name,
		Email:       req.Email,
		IsCafeOwner: req.IsCafeOwner,
	}

	result := database.DB.Create(&user)
	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, "failed to create user: %v", result.Error)
	}

	return &userv1.CreateUserResponse{
		User: &userv1.User{
			Id:          uint32(user.ID),
			Name:        user.Name,
			Email:       user.Email,
			IsCafeOwner: user.IsCafeOwner,
		},
	}, nil
}

func (s *UserServer) GetUser(ctx context.Context, req *userv1.GetUserRequest) (*userv1.GetUserResponse, error) {
	var user models.User
	result := database.DB.First(&user, req.Id)
	if result.Error != nil {
		return nil, status.Errorf(codes.NotFound, "user not found")
	}

	return &userv1.GetUserResponse{
		User: &userv1.User{
			Id:          uint32(user.ID),
			Name:        user.Name,
			Email:       user.Email,
			IsCafeOwner: user.IsCafeOwner,
		},
	}, nil
}

func (s *UserServer) GetUsers(ctx context.Context, req *userv1.GetUsersRequest) (*userv1.GetUsersResponse, error) {
	var users []models.User
	result := database.DB.Find(&users)
	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, "failed to fetch users: %v", result.Error)
	}

	var protoUsers []*userv1.User
	for _, user := range users {
		protoUsers = append(protoUsers, &userv1.User{
			Id:          uint32(user.ID),
			Name:        user.Name,
			Email:       user.Email,
			IsCafeOwner: user.IsCafeOwner,
		})
	}

	return &userv1.GetUsersResponse{
		Users: protoUsers,
	}, nil
}
