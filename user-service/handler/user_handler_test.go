package handler

import (
	"context"
	"errors"
	"testing"
	"github.com/Tao-Zzzz/GoCampus/user-service/proto"
	"github.com/Tao-Zzzz/GoCampus/user-service/service"
	"github.com/Tao-Zzzz/GoCampus/user-service/pkg/jwt"
	"github.com/golang/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Mock dependencies (to be generated using mockgen)
type mockUserService struct{}
type mockJWTUtil struct{}

func (m *mockUserService) RegisterUser(ctx context.Context, email, password, nickname, avatar string) (string, error) {
	if email == "test@example.com" {
		return "user123", nil
	}
	return "", errors.New("invalid email")
}

func (m *mockUserService) Login(ctx context.Context, email, password string) (string, error) {
	if email == "test@example.com" && password == "password" {
		return "valid-token", nil
	}
	return "", errors.New("invalid credentials")
}

func (m *mockUserService) GetUserInfo(ctx context.Context, userID string) (*service.User, error) {
	if userID == "user123" {
		return &service.User{
			ID:       "user123",
			Email:    "test@example.com",
			Nickname: "TestUser",
			Avatar:   "avatar.png",
		}, nil
	}
	return nil, errors.New("user not found")
}

func (m *mockJWTUtil) ValidateTokenFromContext(ctx context.Context) (string, error) {
	if ctx.Value("token") == "valid-token" {
		return "user123", nil
	}
	return "", errors.New("invalid token")
}

func TestUserHandler_RegisterUser(t *testing.T) {
	userService := &mockUserService{}
	jwtUtil := &mockJWTUtil{}
	handler := NewUserHandler(userService, jwtUtil)

	tests := []struct {
		name    string
		req     *proto.RegisterUserRequest
		wantErr bool
	}{
		{
			name: "Successful registration",
			req: &proto.RegisterUserRequest{
				Email:    "test@example.com",
				Password: "password",
				Nickname: "TestUser",
				Avatar:   "avatar.png",
			},
			wantErr: false,
		},
		{
			name: "Invalid email",
			req: &proto.RegisterUserRequest{
				Email:    "invalid",
				Password: "password",
				Nickname: "TestUser",
				Avatar:   "avatar.png",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := handler.RegisterUser(context.Background(), tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("RegisterUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && resp.UserId != "user123" {
				t.Errorf("RegisterUser() got UserId = %v, want %v", resp.UserId, "user123")
			}
		})
	}
}

// Additional tests for Login and GetUserInfo can be added similarly