package handler

import (
	"context"
	"errors"
	"github.com/Tao-Zzzz/GoCampus/user-service/model"
	"github.com/Tao-Zzzz/GoCampus/user-service/config"
	"github.com/Tao-Zzzz/GoCampus/user-service/proto"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"golang.org/x/crypto/bcrypt"
	"testing"
	"time"
)

// MockUserRepository for testing the service layer.
type MockUserRepository struct {
	CreateUserFunc    func(ctx context.Context, user *model.User) (string, error)
	GetUserByEmailFunc func(ctx context.Context, email string) (*model.User, error)
	GetUserByIDFunc   func(ctx context.Context, id string) (*model.User, error)
}

func (m *MockUserRepository) CreateUser(ctx context.Context, user *model.User) (string, error) {
	return m.CreateUserFunc(ctx, user)
}

func (m *MockUserRepository) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	return m.GetUserByEmailFunc(ctx, email)
}

func (m *MockUserRepository) GetUserByID(ctx context.Context, id string) (*model.User, error) {
	return m.GetUserByIDFunc(ctx, id)
}

func TestUserHandler_RegisterUser(t *testing.T) {
	mockRepo := &MockUserRepository{
		CreateUserFunc: func(ctx context.Context, user *model.User) (string, error) {
			return user.ID, nil
		},
		GetUserByEmailFunc: func(ctx context.Context, email string) (*model.User, error) {
			return nil, errors.New("user not found")
		},
	}
	cfg := &config.Config{
		JWT: config.JWTConfig{
			Secret:        "secret-key",
			DurationHours: 24,
		},
	}
	handler := NewUserHandler(mockRepo, cfg)

	tests := []struct {
		name     string
		req      *proto.RegisterRequest
		wantResp *proto.RegisterResponse
	}{
		{
			name: "Successful registration",
			req: &proto.RegisterRequest{
				Email:    "test@example.com",
				Password: "password123",
				Nickname: "TestUser",
				Avatar:   "http://example.com/avatar.png",
			},
			wantResp: &proto.RegisterResponse{
				Success: true,
				Message: "User registered successfully",
			},
		},
		{
			name: "Invalid input",
			req: &proto.RegisterRequest{
				Email:    "",
				Password: "password123",
				Nickname: "TestUser",
				Avatar:   "http://example.com/avatar.png",
			},
			wantResp: &proto.RegisterResponse{
				Success: false,
				Message: "email, password, and nickname are required",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := handler.RegisterUser(context.Background(), tt.req)
			if err != nil {
				t.Fatalf("RegisterUser() error = %v", err)
			}
			if resp.Success != tt.wantResp.Success || resp.Message != tt.wantResp.Message {
				t.Errorf("RegisterUser() = %+v, want %+v", resp, tt.wantResp)
			}
			if resp.Success && resp.UserId == "" {
				t.Errorf("RegisterUser() expected non-empty userID")
			}
			count := testutil.ToFloat64(requestCounter.WithLabelValues("RegisterUser", "success"))
			if tt.wantResp.Success && count == 0 {
				t.Errorf("Expected RegisterUser success metric to be recorded")
			}
		})
	}
}

func TestUserHandler_Login(t *testing.T) {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	mockRepo := &MockUserRepository{
		GetUserByEmailFunc: func(ctx context.Context, email string) (*model.User, error) {
			if email == "test@example.com" {
				return &model.User{
					ID:        "user123",
					Email:     email,
					Password:  string(hashedPassword),
					Nickname:  "TestUser",
					Avatar:    "http://example.com/avatar.png",
					CreatedAt: time.Now(),
				}, nil
			}
			return nil, errors.New("user not found")
		},
	}
	cfg := &config.Config{
		JWT: config.JWTConfig{
			Secret:        "secret-key",
			DurationHours: 24,
		},
	}
	handler := NewUserHandler(mockRepo, cfg)

	tests := []struct {
		name     string
		req      *proto.LoginRequest
		wantResp *proto.LoginResponse
	}{
		{
			name: "Successful login",
			req: &proto.LoginRequest{
				Email:    "test@example.com",
				Password: "password123",
			},
			wantResp: &proto.LoginResponse{
				Success: true,
				Message: "Login successful",
			},
		},
		{
			name: "Invalid credentials",
			req: &proto.LoginRequest{
				Email:    "test@example.com",
				Password: "wrongpassword",
			},
			wantResp: &proto.LoginResponse{
				Success: false,
				Message: "invalid credentials",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := handler.Login(context.Background(), tt.req)
			if err != nil {
				t.Fatalf("Login() error = %v", err)
			}
			if resp.Success != tt.wantResp.Success || resp.Message != tt.wantResp.Message {
				t.Errorf("Login() = %+v, want %+v", resp, tt.wantResp)
			}
			if tt.wantResp.Success && resp.Token == "" {
				t.Errorf("Login() expected non-empty token")
			}
		})
	}
}

func TestUserHandler_GetUserInfo(t *testing.T) {
	mockRepo := &MockUserRepository{
		GetUserByIDFunc: func(ctx context.Context, id string) (*model.User, error) {
			if id == "user123" {
				return &model.User{
					ID:        "user123",
					Email:     "test@example.com",
					Nickname:  "TestUser",
					Avatar:    "http://example.com/avatar.png",
					CreatedAt: time.Now(),
				}, nil
			}
			return nil, errors.New("user not found")
		},
	}
	cfg := &config.Config{
		JWT: config.JWTConfig{
			Secret:        "secret-key",
			DurationHours: 24,
		},
	}
	handler := NewUserHandler(mockRepo, cfg)

	tests := []struct {
		name     string
		req      *proto.GetUserInfoRequest
		wantResp *proto.GetUserInfoResponse
	}{
		{
			name: "Successful get user info",
			req: &proto.GetUserInfoRequest{
				UserId: "user123",
			},
			wantResp: &proto.GetUserInfoResponse{
				Success: true,
				Message: "User info retrieved successfully",
				User: &proto.UserInfo{
					UserId:   "user123",
					Email:    "test@example.com",
					Nickname: "TestUser",
					Avatar:   "http://example.com/avatar.png",
				},
			},
		},
		{
			name: "User not found",
			req: &proto.GetUserInfoRequest{
				UserId: "invalid",
			},
			wantResp: &proto.GetUserInfoResponse{
				Success: false,
				Message: "user not found",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := handler.GetUserInfo(context.Background(), tt.req)
			if err != nil {
				t.Fatalf("GetUserInfo() error = %v", err)
			}
			if resp.Success != tt.wantResp.Success || resp.Message != tt.wantResp.Message {
				t.Errorf("GetUserInfo() = %+v, want %+v", resp, tt.wantResp)
			}
			if resp.User != nil && tt.wantResp.User != nil {
				if resp.User.UserId != tt.wantResp.User.UserId ||
					resp.User.Email != tt.wantResp.User.Email ||
					resp.User.Nickname != tt.wantResp.User.Nickname ||
					resp.User.Avatar != tt.wantResp.User.Avatar {
					t.Errorf("GetUserInfo() User = %+v, want %+v", resp.User, tt.wantResp.User)
				}
			}
		})
	}
}