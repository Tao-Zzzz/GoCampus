package handler

import (
    "context"
    "github.com/Tao-Zzzz/GoCampus/user-service/proto"
    "github.com/Tao-Zzzz/GoCampus/user-service/config"
    "github.com/Tao-Zzzz/GoCampus/user-service/service"
    "github.com/prometheus/client_golang/prometheus"
    "go.opentelemetry.io/otel"
    "go.opentelemetry.io/otel/attribute"
    "go.opentelemetry.io/otel/trace"
)

var (
    requestCounter = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "user_service_requests_total",
            Help: "Total number of requests to user service endpoints",
        },
        []string{"method", "status"},
    )
)

func init() {
    // 只注册一次
    prometheus.MustRegister(requestCounter)
}

// UserHandler implements the gRPC UserService server.
type UserHandler struct {
    proto.UnimplementedUserServiceServer
    userService *service.UserService
}

// NewUserHandler creates a new UserHandler with dependencies.
func NewUserHandler(repo service.UserRepository, cfg *config.Config) *UserHandler {
    userService := service.NewUserService(repo, cfg)
    return &UserHandler{
        userService: userService,
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
        requestCounter.WithLabelValues("RegisterUser", "error").Inc()
        return &proto.RegisterResponse{
            Success: false,
            Message: err.Error(),
        }, nil
    }
    requestCounter.WithLabelValues("RegisterUser", "success").Inc()
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
        requestCounter.WithLabelValues("Login", "error").Inc()
        return &proto.LoginResponse{
            Success: false,
            Message: err.Error(),
        }, nil
    }
    requestCounter.WithLabelValues("Login", "success").Inc()
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
        requestCounter.WithLabelValues("GetUserInfo", "error").Inc()
        return &proto.GetUserInfoResponse{
            Success: false,
            Message: err.Error(),
        }, nil
    }
    requestCounter.WithLabelValues("GetUserInfo", "success").Inc()
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