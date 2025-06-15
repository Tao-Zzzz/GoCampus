package service

import (
	"context"
	"errors"
	"time"

	"github.com/Tao-Zzzz/GoCampus/user-service/model"
	"github.com/Tao-Zzzz/GoCampus/user-service/pkg/jwt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// UserRepository defines the interface for data access.
type UserRepository interface {
	CreateUser(ctx context.Context, user *model.User) (string, error)
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
	GetUserByID(ctx context.Context, id string) (*model.User, error)
}

// UserService handles user-related business logic.
type UserService struct {
	repo   UserRepository
	jwtKey string
}

// NewUserService creates a new UserService instance.
func NewUserService(repo UserRepository, jwtKey string) *UserService {
	return &UserService{
		repo:   repo,
		jwtKey: jwtKey,
	}
}

// Register creates a new user with hashed password.
func (s *UserService) Register(ctx context.Context, email, password, nickname, avatar string) (string, error) {
	// Validate input
	if email == "" || password == "" || nickname == "" {
		return "", errors.New("email, password, and nickname are required")
	}

	// Check if user already exists
	_, err := s.repo.GetUserByEmail(ctx, email)
	if err == nil {
		return "", errors.New("user already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", errors.New("failed to hash password")
	}

	// Create user
	user := &model.User{
		ID:        uuid.New().String(),
		Email:     email,
		Password:  string(hashedPassword),
		Nickname:  nickname,
		Avatar:    avatar,
		CreatedAt: time.Now(),
	}

	userID, err := s.repo.CreateUser(ctx, user)
	if err != nil {
		return "", err
	}
	return userID, nil
}

// Login authenticates a user and generates a JWT token.
func (s *UserService) Login(ctx context.Context, email, password string) (string, error) {
	// Validate input
	if email == "" || password == "" {
		return "", errors.New("email and password are required")
	}

	// Get user by email
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", errors.New("invalid credentials")
	}

	// Generate JWT token
	token, err := jwt.GenerateToken(user.ID, s.jwtKey, 24*time.Hour)
	if err != nil {
		return "", errors.New("failed to generate token")
	}
	return token, nil
}

// GetUserInfo retrieves user information by ID.
func (s *UserService) GetUserInfo(ctx context.Context, userID string) (*model.User, error) {
	if userID == "" {
		return nil, errors.New("user ID is required")
	}

	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, errors.New("user not found")
	}
	return user, nil
}