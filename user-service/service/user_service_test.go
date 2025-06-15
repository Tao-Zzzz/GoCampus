package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Tao-Zzzz/GoCampus/user-service/config"
	"github.com/Tao-Zzzz/GoCampus/user-service/model"
	"github.com/Tao-Zzzz/GoCampus/user-service/pkg/jwt"
	"github.com/Tao-Zzzz/GoCampus/user-service/pkg/logger"
	"github.com/Tao-Zzzz/GoCampus/user-service/pkg/metrics"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

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

func TestUserService_Register(t *testing.T) {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	mockRepo := &MockUserRepository{
		CreateUserFunc: func(ctx context.Context, user *model.User) (string, error) {
			if user.Email == "test@example.com" {
				return user.ID, nil
			}
			return "", errors.New("failed to create user")
		},
		GetUserByEmailFunc: func(ctx context.Context, email string) (*model.User, error) {
			if email == "existing@example.com" {
				return &model.User{ID: "user123", Email: email}, nil
			}
			return nil, errors.New("user not found")
		},
	}
	cfg := &config.Config{
		JWT: config.JWTConfig{
			Secret:        "secret-key",
			DurationHours: 24,
		},
		Service: config.ServiceConfig{
			Name:     "test-service",
			Port:     8080,
			LogLevel: "debug",
		},
	}
	log := logger.NewLogger(cfg)
	met := metrics.NewMetrics(cfg)
	service := NewUserService(mockRepo, cfg, log, met)


	tests := []struct {
		name       string
		email      string
		password   string
		nickname   string
		avatar     string
		wantErr    bool
		wantUserID string
	}{
		{
			name:       "Successful registration",
			email:      "test@example.com",
			password:   "password123",
			nickname:   "TestUser",
			avatar:     "http://example.com/avatar.png",
			wantErr:    false,
			wantUserID: uuid.New().String(), // UUID is generated dynamically
		},
		{
			name:     "User already exists",
			email:    "existing@example.com",
			password: "password123",
			nickname: "TestUser",
			avatar:   "http://example.com/avatar.png",
			wantErr:  true,
		},
		{
			name:     "Invalid input",
			email:    "",
			password: "password123",
			nickname: "TestUser",
			avatar:   "http://example.com/avatar.png",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userID, err := service.Register(context.Background(), tt.email, tt.password, tt.nickname, tt.avatar)
			if (err != nil) != tt.wantErr {
				t.Errorf("Register() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && userID == "" {
				t.Errorf("Register() expected non-empty userID")
			}
		})
	}
}

func TestUserService_Login(t *testing.T) {
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
		Service: config.ServiceConfig{
			Name:     "test-service",
			Port:     8080,
			LogLevel: "debug",
		},
	}
	log := logger.NewLogger(cfg)
	met := metrics.NewMetrics(cfg)
	service := NewUserService(mockRepo, cfg, log, met)

	tests := []struct {
		name     string
		email    string
		password string
		wantErr  bool
	}{
		{
			name:     "Successful login",
			email:    "test@example.com",
			password: "password123",
			wantErr:  false,
		},
		{
			name:     "Invalid credentials",
			email:    "test@example.com",
			password: "wrongpassword",
			wantErr:  true,
		},
		{
			name:     "User not found",
			email:    "invalid@example.com",
			password: "password123",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := service.Login(context.Background(), tt.email, tt.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("Login() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && token == "" {
				t.Errorf("Login() expected non-empty token")
			}
		})
	}
}

func TestUserService_GetUserInfo(t *testing.T) {
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
		Service: config.ServiceConfig{
			Name:     "test-service",
			Port:     8080,
			LogLevel: "debug",
		},
	}
	log := logger.NewLogger(cfg)
	met := metrics.NewMetrics(cfg)
	service := NewUserService(mockRepo, cfg, log, met)

	tests := []struct {
		name     string
		userID   string
		wantErr  bool
		wantUser *model.User
	}{
		{
			name:   "Successful get user info",
			userID: "user123",
			wantErr: false,
			wantUser: &model.User{
				ID:        "user123",
				Email:     "test@example.com",
				Nickname:  "TestUser",
				Avatar:    "http://example.com/avatar.png",
			},
		},
		{
			name:    "User not found",
			userID:  "invalid",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := service.GetUserInfo(context.Background(), tt.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUserInfo() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && user == nil {
				t.Errorf("GetUserInfo() expected non-nil user")
			}
			if !tt.wantErr && (user.ID != tt.wantUser.ID || user.Email != tt.wantUser.Email ||
				user.Nickname != tt.wantUser.Nickname || user.Avatar != tt.wantUser.Avatar) {
				t.Errorf("GetUserInfo() user = %+v, want %+v", user, tt.wantUser)
			}
		})
	}
}