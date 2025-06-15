package handler

import (
    "context"
    "errors"
    "time"

    "github.com/google/uuid"
    "github.com/Tao-Zzzz/GoCampus/user-service/config"
    "github.com/Tao-Zzzz/GoCampus/user-service/model"
    "github.com/Tao-Zzzz/GoCampus/user-service/pkg/jwt"
    "github.com/Tao-Zzzz/GoCampus/user-service/pkg/logger"
    "github.com/Tao-Zzzz/GoCampus/user-service/pkg/metrics"
    "github.com/Tao-Zzzz/GoCampus/user-service/proto"
    "github.com/Tao-Zzzz/GoCampus/user-service/service"
    "github.com/prometheus/client_golang/prometheus"
    "go.opentelemetry.io/otel"
    "go.opentelemetry.io/otel/attribute"
)

// 包级指标，只注册一次，避免重复注册 panic
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
    prometheus.MustRegister(requestCounter)
}

// UserHandler implements the gRPC UserServiceServer.
type UserHandler struct {
    proto.UnimplementedUserServiceServer
    userService *service.UserService
    logger      *logger.Logger
    metrics     *metrics.Metrics
    cfg         *config.Config
}

// NewUserHandler creates a new UserHandler with dependencies.
func NewUserHandler(repo service.UserRepository, cfg *config.Config, log *logger.Logger, met *metrics.Metrics) *UserHandler {
    userService := service.NewUserService(repo, cfg, log, met)
    return &UserHandler{
        userService: userService,
        logger:      log,
        metrics:     met,
        cfg:         cfg,
    }
}

// RegisterUser handles user registration requests.
func (h *UserHandler) RegisterUser(ctx context.Context, req *proto.RegisterUserRequest) (*proto.RegisterUserResponse, error) {
    tracer := otel.Tracer("user-service")
    ctx, span := tracer.Start(ctx, "UserHandler.RegisterUser")
    defer span.End()

    start := time.Now()
    defer func() {
        duration := time.Since(start).Seconds()
        h.metrics.RequestDuration().WithLabelValues("RegisterUser", "success").Observe(duration)
        requestCounter.WithLabelValues("RegisterUser", "success").Inc()
    }()

    h.logger.Info(ctx).Msgf("Received RegisterUser request for email: %s", req.Email)

    // 输入校验
    if req.Email == "" || req.Password == "" || req.Nickname == "" {
        err := errors.New("email, password, nickname are required")
        h.logger.Warn(ctx).Err(err).Msg("Invalid register input")
        h.metrics.RequestDuration().WithLabelValues("RegisterUser", "error").Observe(time.Since(start).Seconds())
        requestCounter.WithLabelValues("RegisterUser", "error").Inc()
        span.RecordError(err)
        return nil, err
    }

    user := &model.User{
        ID:       uuid.New().String(),
        Email:    req.Email,
        Password: req.Password,
        Nickname: req.Nickname,
        Avatar:   req.Avatar,
    }

    userID, err := h.userService.Register(ctx, user)
    if err != nil {
        h.logger.Error(ctx).Err(err).Msg("Failed to register user")
        h.metrics.RequestDuration().WithLabelValues("RegisterUser", "error").Observe(time.Since(start).Seconds())
        requestCounter.WithLabelValues("RegisterUser", "error").Inc()
        span.RecordError(err)
        return nil, err
    }

    h.logger.Info(ctx).Msgf("User registered: %s", userID)
    span.SetAttributes(attribute.String("user_id", userID))
    return &proto.RegisterUserResponse{UserId: userID}, nil
}

// Login handles user login requests.
func (h *UserHandler) Login(ctx context.Context, req *proto.LoginRequest) (*proto.LoginResponse, error) {
    tracer := otel.Tracer("user-service")
    ctx, span := tracer.Start(ctx, "UserHandler.Login")
    defer span.End()

    start := time.Now()
    defer func() {
        duration := time.Since(start).Seconds()
        h.metrics.RequestDuration().WithLabelValues("Login", "success").Observe(duration)
        requestCounter.WithLabelValues("Login", "success").Inc()
    }()

    h.logger.Info(ctx).Msgf("Received Login request for email: %s", req.Email)

    if req.Email == "" || req.Password == "" {
        err := errors.New("email and password are required")
        h.logger.Warn(ctx).Err(err).Msg("Invalid login input")
        h.metrics.RequestDuration().WithLabelValues("Login", "error").Observe(time.Since(start).Seconds())
        requestCounter.WithLabelValues("Login", "error").Inc()
        span.RecordError(err)
        return nil, err
    }

    token, err := h.userService.Login(ctx, req.Email, req.Password)
    if err != nil {
        h.logger.Error(ctx).Err(err).Msg("Failed to login user")
        h.metrics.RequestDuration().WithLabelValues("Login", "error").Observe(time.Since(start).Seconds())
        requestCounter.WithLabelValues("Login", "error").Inc()
        span.RecordError(err)
        return nil, err
    }

    h.logger.Info(ctx).Msg("User logged in successfully")
    return &proto.LoginResponse{Token: token}, nil
}

// GetUserInfo handles user info requests with JWT authentication.
func (h *UserHandler) GetUserInfo(ctx context.Context, req *proto.GetUserInfoRequest) (*proto.GetUserInfoResponse, error) {
    tracer := otel.Tracer("user-service")
    ctx, span := tracer.Start(ctx, "UserHandler.GetUserInfo")
    defer span.End()

    start := time.Now()
    defer func() {
        duration := time.Since(start).Seconds()
        h.metrics.RequestDuration().WithLabelValues("GetUserInfo", "success").Observe(duration)
        requestCounter.WithLabelValues("GetUserInfo", "success").Inc()
    }()

    h.logger.Info(ctx).Msg("Received GetUserInfo request")

    // Validate JWT token
    userID, err := jwt.ValidateToken(req.Token, h.cfg.JWT.Secret)
    if err != nil {
        h.logger.Error(ctx).Err(err).Msg("Invalid JWT token")
        h.metrics.RequestDuration().WithLabelValues("GetUserInfo", "error").Observe(time.Since(start).Seconds())
        requestCounter.WithLabelValues("GetUserInfo", "error").Inc()
        span.RecordError(err)
        return nil, errors.New("invalid token")
    }

    user, err := h.userService.GetUserInfo(ctx, userID)
    if err != nil {
        h.logger.Error(ctx).Err(err).Msg("Failed to get user info")
        h.metrics.RequestDuration().WithLabelValues("GetUserInfo", "error").Observe(time.Since(start).Seconds())
        requestCounter.WithLabelValues("GetUserInfo", "error").Inc()
        span.RecordError(err)
        return nil, err
    }

    h.logger.Info(ctx).Msgf("User info retrieved for ID: %s", userID)
    span.SetAttributes(attribute.String("user_id", userID))
    return &proto.GetUserInfoResponse{
        UserId:   user.ID,
        Email:    user.Email,
        Nickname: user.Nickname,
        Avatar:   user.Avatar,
    }, nil
}