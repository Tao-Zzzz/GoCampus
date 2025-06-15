package service

import (
	"context"
	"errors"
	"time"

	"github.com/Tao-Zzzz/GoCampus/user-service/config"
	"github.com/Tao-Zzzz/GoCampus/user-service/model"
	"github.com/Tao-Zzzz/GoCampus/user-service/pkg/jwt"
	"github.com/Tao-Zzzz/GoCampus/user-service/pkg/logger"
	"github.com/Tao-Zzzz/GoCampus/user-service/pkg/metrics"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"golang.org/x/crypto/bcrypt"
)

// UserRepository defines the interface for data access.
type UserRepository interface {
	CreateUser(ctx context.Context, user *model.User) (string, error)
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
	GetUserByID(ctx context.Context, id string) (*model.User, error)
}

// UserService implements user-related business logic.
type UserService struct {
	repo         UserRepository
	cfg          *config.Config
	logger       *logger.Logger
	metrics      *metrics.Metrics
	tracer       otel.Tracer
	jwtKey       string
}

// NewUserService creates a new UserService instance.
func NewUserService(repo UserRepository, cfg *config.Config, log *logger.Logger, met *metrics.Metrics) *UserService {
	return &UserService{
		repo:    repo,
		cfg:     cfg,
		logger:  log,
		metrics: met,
		tracer:  otel.Tracer("user-service"),
		jwtKey:  cfg.JWT.Secret,
	}
}

// Register creates a new user with hashed password.
func (s *UserService) Register(ctx context.Context, user *model.User) (string, error) {
	ctx, span := s.tracer.Start(ctx, "UserService.Register")
	defer span.End()

	start := time.Now()
	defer func() {
		duration := time.Since(start).Seconds()
		s.metrics.RequestDuration().WithLabelValues("Register", "success").Observe(duration)
	}()

	s.logger.Info(ctx).Msgf("Registering user with email: %s", user.Email)

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
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		s.logger.Error(ctx).Err(err).Msg("Failed to hash password")
		s.metrics.RequestDuration().WithLabelValues("Register", "error").Observe(time.Since(start).Seconds())
		span.RecordError(err)
		return "", errors.New("failed to hash password")
	}
	user.Password = string(hashedPassword)

	// Create user
	userID, err := s.repo.CreateUser(ctx, user)
	if err != nil {
		s.logger.Error(ctx).Err(err).Msg("Failed to create user")
		s.metrics.RequestDuration().WithLabelValues("Register", "error").Observe(time.Since(start).Seconds())
		span.RecordError(err)
		return "", errors.New("failed to create user")
	}

	s.logger.Info(ctx).Msgf("User registered successfully: %s", userID)
	span.SetAttributes(attribute.String("user_id", userID))
	return userID, nil
}

// Login authenticates a user and generates a JWT token.
func (s *UserService) Login(ctx context.Context, email, password string) (string, error) {
	ctx, span := s.tracer.Start(ctx, "UserService.Login")
	defer span.End()

	start := time.Now()
	defer func() {
		duration := time.Since(start).Seconds()
		s.metrics.RequestDuration().WithLabelValues("Login", "success").Observe(duration)
	}()

	s.logger.Info(ctx).Msgf("Logging in user with email: %s", email)

	// Validate input
	if email == "" || password == "" {
		s.logger.Warn(ctx).Msg("Empty email or password provided")
		s.metrics.RequestDuration().WithLabelValues("Login", "error").Observe(time.Since(start).Seconds())
		span.RecordError(errors.New("empty email or password"))
		return "", errors.New("email and password are required")
	}

	// Get user by email
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		s.logger.Error(ctx).Err(err).Msg("Failed to get user by email")
		s.metrics.RequestDuration().WithLabelValues("Login", "error").Observe(time.Since(start).Seconds())
		span.RecordError(err)
		return "", errors.New("invalid credentials")
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		s.logger.Warn(ctx).Msg("Invalid password provided")
		s.metrics.RequestDuration().WithLabelValues("Login", "error").Observe(time.Since(start).Seconds())
		span.RecordError(errors.New("invalid password"))
		return "", errors.New("invalid credentials")
	}

	// Generate JWT token
	token, err := jwt.GenerateToken(user.ID, s.jwtKey, time.Duration(s.cfg.JWT.DurationHours)*time.Hour)
	if err != nil {
		s.logger.Error(ctx).Err(err).Msg("Failed to generate JWT token")
		s.metrics.RequestDuration().WithLabelValues("Login", "error").Observe(time.Since(start).Seconds())
		span.RecordError(err)
		return "", errors.New("failed to generate token")
	}

	s.logger.Info(ctx).Msgf("User logged in successfully: %s", user.ID)
	span.SetAttributes(attribute.String("user_id", user.ID))
	return token, nil
}

// GetUserInfo retrieves user information by ID.
func (s *UserService) GetUserInfo(ctx context.Context, userID string) (*model.User, error) {
	ctx, span := s.tracer.Start(ctx, "UserService.GetUserInfo")
	defer span.End()

	start := time.Now()
	defer func() {
		duration := time.Since(start).Seconds()
		s.metrics.RequestDuration().WithLabelValues("GetUserInfo", "success").Observe(duration)
	}()

	s.logger.Info(ctx).Msgf("Retrieving user info for ID: %s", userID)

	// Validate input
	if userID == "" {
		s.logger.Warn(ctx).Msg("Empty user ID provided")
		s.metrics.RequestDuration().WithLabelValues("GetUserInfo", "error").Observe(time.Since(start).Seconds())
		span.RecordError(errors.New("empty user ID"))
		return nil, errors.New("user ID is required")
	}

	// Get user by ID
	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		s.logger.Error(ctx).Err(err).Msg("Failed to get user by ID")
		s.metrics.RequestDuration().WithLabelValues("GetUserInfo", "error").Observe(time.Since(start).Seconds())
		span.RecordError(err)
		return nil, errors.New("failed to get user")
	}

	s.logger.Info(ctx).Msgf("User info retrieved successfully: %s", userID)
	span.SetAttributes(attribute.String("user_id", userID))
	return user, nil
}

// // jwtDuration returns the JWT token duration from the config.
// func (s *UserService) jwtDuration() time.Duration {
// 	return time.Duration(s.cfg.JWT.DurationHours) * time.Hour
// }