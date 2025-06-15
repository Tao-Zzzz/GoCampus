package handler

import (
	"context"
	"github.com/Tao-Zzzz/GoCampus/user-service/proto"
	"github.com/Tao-Zzzz/GoCampus/user-service/service"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"sync"
)

// UserHandler implements the gRPC UserService server.
type UserHandler struct {
	proto.UnimplementedUserServiceServer
	userService *service.UserService
	requestCounter *prometheus.CounterVec
}

var (
	requestCounter     *prometheus.CounterVec
	registerMetricsOnce sync.Once
)

func initMetrics() {
	requestCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "user_service_requests_total",
			Help: "Total number of requests to user service endpoints",
		},
		[]string{"method", "status"},
	)
	prometheus.MustRegister(requestCounter)
}


// NewUserHandler creates a new UserHandler with dependencies.
func NewUserHandler(userService *service.UserService) *UserHandler {
	registerMetricsOnce.Do(initMetrics)
	return &UserHandler{
		userService:    userService,
		requestCounter: requestCounter,
	}
}

// RegisterUser handles user registration.
func (h *UserHandler) RegisterUser(ctx context.Context, req *proto.RegisterRequest) (*proto.RegisterResponse, error) {
	tracer := otel.Tracer("user-service")
	ctx, span := tracer.Start(ctx, "RegisterUser",
		trace.WithAttributes(
			attribute.String("email", req.Email),
			attribute.String("nickname", req.Nickname),
		))
	defer span.End()

	userID, err := h.userService.Register(ctx, req.Email, req.Password, req.Nickname, req.Avatar)
	if err != nil {
		h.requestCounter.WithLabelValues("RegisterUser", "error").Inc()
		return &proto.RegisterResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}
	h.requestCounter.WithLabelValues("RegisterUser", "success").Inc()
	return &proto.RegisterResponse{
		Success: true,
		Message: "User registered successfully",
		UserId:  userID,
	}, nil
}

// Login handles user authentication and JWT generation.
func (h *UserHandler) Login(ctx context.Context, req *proto.LoginRequest) (*proto.LoginResponse, error) {
	tracer := otel.Tracer("user-service")
	ctx, span := tracer.Start(ctx, "Login",
		trace.WithAttributes(
			attribute.String("email", req.Email),
		))
	defer span.End()

	token, err := h.userService.Login(ctx, req.Email, req.Password)
	if err != nil {
		h.requestCounter.WithLabelValues("Login", "error").Inc()
		return &proto.LoginResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}
	h.requestCounter.WithLabelValues("Login", "success").Inc()
	return &proto.LoginResponse{
		Success: true,
		Message: "Login successful",
		Token:   token,
	}, nil
}

// GetUserInfo retrieves user information.
func (h *UserHandler) GetUserInfo(ctx context.Context, req *proto.GetUserInfoRequest) (*proto.GetUserInfoResponse, error) {
	tracer := otel.Tracer("user-service")
	ctx, span := tracer.Start(ctx, "GetUserInfo",
		trace.WithAttributes(
			attribute.String("user_id", req.UserId),
		))
	defer span.End()

	user, err := h.userService.GetUserInfo(ctx, req.UserId)
	if err != nil {
		h.requestCounter.WithLabelValues("GetUserInfo", "error").Inc()
		return &proto.GetUserInfoResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}
	h.requestCounter.WithLabelValues("GetUserInfo", "success").Inc()
	return &proto.GetUserInfoResponse{
		Success: true,
		Message: "User info retrieved successfully",
		User: &proto.UserInfo{
			UserId:   user.ID,
			Email:    user.Email,
			Nickname: user.Nickname,
			Avatar:   user.Avatar,
		},
	}, nil
}