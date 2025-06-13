package service

import (
	"context"
	"time"

	"github.com/campus-trading/user-service/internal/config"
	"github.com/campus-trading/user-service/internal/model"
	"github.com/campus-trading/user-service/internal/repository/postgres"
	"github.com/campus-trading/user-service/internal/repository/redis"
	"github.com/campus-trading/user-service/pkg/util"
	"github.com/google/uuid"
)

type UserService struct {
	pgRepo    *postgres.UserRepository
	redisRepo *redis.SessionRepository
	cfg       *config.Config
}

func NewUserService(pgRepo *postgres.UserRepository, redisRepo *redis.SessionRepository, cfg *config.Config) *UserService {
	return &UserService{
		pgRepo:    pgRepo,
		redisRepo: redisRepo,
		cfg:       cfg,
	}
}

func (s *UserService) Register(ctx context.Context, email, phone, password, nickname string) (string, error) {
	// 加密密码
	hash, err := util.HashPassword(password)
	if err != nil {
		return "", err
	}

	user := &model.User{
		UserID:      uuid.New().String(),
		Email:       email,
		Phone:       phone,
		PasswordHash: hash,
		Nickname:    nickname,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.pgRepo.Create(ctx, user); err != nil {
		return "", err
	}

	// 异步发送欢迎消息（通过 NATS）
	// TODO: 实现 NATS 发布
	return user.UserID, nil
}

func (s *UserService) Login(ctx context.Context, email, password string) (string, string, error) {
	user, err := s.pgRepo.GetByEmail(ctx, email)
	if err != nil {
		return "", "", err
	}

	if !util.CheckPasswordHash(password, user.PasswordHash) {
		return "", "", errors.New("invalid password")
	}

	token, err := util.GenerateJWT(user.UserID, s.cfg.JWT.Secret)
	if err != nil {
		return "", "", err
	}

	// 存储会话到 Redis
	if err := s.redisRepo.SaveSession(ctx, user.UserID, token, 24*time.Hour); err != nil {
		return "", "", err
	}

	return token, user.UserID, nil
}

func (s *UserService) UpdateProfile(ctx context.Context, userID, nickname, avatarURL, contactInfo string) error {
	user := &model.User{
		UserID:      userID,
		Nickname:    nickname,
		AvatarURL:   avatarURL,
		ContactInfo: contactInfo,
		UpdatedAt:   time.Now(),
	}
	return s.pgRepo.Update(ctx, user)
}

func (s *UserService) GetProfile(ctx context.Context, userID string) (*model.User, error) {
	return s.pgRepo.GetByID(ctx, userID)
}