package handler

import (
	"context"
	"github.com/Tao-Zzzz/GoCampus/user-service/proto"
	"github.com/Tao-Zzzz/GoCampus/user-service/service"
	"github.com/Tao-Zzzz/GoCampus/user-service/pkg/jwt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// UserHandler implements the UserService gRPC interface.
type UserHandler struct {
	proto.UnimplementedUserServiceServer
	userService *service.UserService
	jwtUtil     *jwt.JWTUtil
}

// NewUserHandler creates a new UserHandler with dependencies.
func NewUserHandler(userService *service.UserService, jwtUtil *jwt.JWTUtil) *UserHandler {
	return &UserHandler{
		userService: userService,
		jwtUtil:     jwtUtil,
	}
}

// RegisterUser handles user registration requests.
func (h *UserHandler) RegisterUser(ctx context.Context, req *proto.RegisterUserRequest) (*proto.RegisterUserResponse, error) {
	userID, err := h.userService.RegisterUser(ctx, req.Email, req.Password, req.Nickname, req.Avatar)
	if err != nil {
		return &proto.RegisterUserResponse{
			Success: false,
			Message: err.Error(),
		}, status.Errorf(codes.InvalidArgument, err.Error())
	}
	return &proto.RegisterUserResponse{
		UserId:  userID,
		Success: true,
		Message: "User registered successfully",
	}, nil
}

// Login handles user login and returns a JWT token.
func (h *UserHandler) Login(ctx context.Context, req *proto.LoginRequest) (*proto.LoginResponse, error) {
	token, err := h.userService.Login(ctx, req.Email, req.Password)
	if err != nil {
		return &proto.LoginResponse{
			Success: false,
			Message: err.Error(),
		}, status.Errorf(codes.Unauthenticated, err.Error())
	}
	return &proto.LoginResponse{
		Token:   token,
		Success: true,
		Message: "Login successful",
	}, nil
}

// GetUserInfo retrieves user information, requires JWT authentication.
func (h *UserHandler) GetUserInfo(ctx context.Context, req *proto.GetUserInfoRequest) (*proto.GetUserInfoResponse, error) {
	// Validate JWT token from metadata
	userID, err := h.jwtUtil.ValidateTokenFromContext(ctx)
	if err != nil || userID != req.UserId {
		return &proto.GetUserInfoResponse{
			Success: false,
			Message: "Invalid or missing JWT token",
		}, status.Errorf(codes.Unauthenticated, "Invalid token")
	}

	user, err := h.userService.GetUserInfo(ctx, req.UserId)
	if err != nil {
		return &proto.GetUserInfoResponse{
			Success: false,
			Message: err.Error(),
		}, status.Errorf(codes.NotFound, err.Error())
	}
	return &proto.GetUserInfoResponse{
		UserId:   user.ID,
		Email:    user.Email,
		Nickname: user.Nickname,
		Avatar:   user.Avatar,
		Success:  true,
		Message:  "User info retrieved successfully",
	}, nil
}