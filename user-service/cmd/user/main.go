package main

import (
	"context"
	"log"

	"github.com/campus-trading/user-service/internal/config"
	"github.com/campus-trading/user-service/internal/handler"
	"github.com/campus-trading/user-service/internal/repository/postgres"
	"github.com/campus-trading/user-service/internal/repository/redis"
	"github.com/campus-trading/user-service/internal/service"
	"github.com/campus-trading/user-service/pkg/proto"
	"github.com/micro/go-micro/v3"
	"github.com/micro/go-micro/v3/registry"
	"github.com/micro/go-micro/v3/server"
)

func main() {
	// 加载配置
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 初始化数据库
	pgRepo, err := postgres.NewUserRepository(cfg.Postgres.DSN)
	if err != nil {
		log.Fatalf("Failed to connect to Postgres: %v", err)
	}
	defer pgRepo.Close()

	redisRepo, err := redis.NewSessionRepository(cfg.Redis.Addr, cfg.Redis.Password)
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	// 初始化服务
	userService := service.NewUserService(pgRepo, redisRepo, cfg)

	// 创建 Go Micro 服务
	srv := micro.NewService(
		micro.Name("user-service"),
		micro.Registry(registry.NewRegistry()),
		micro.BeforeStart(func() error {
			log.Println("User service starting...")
			return nil
		}),
	)

	// 注册 gRPC 处理程序
	proto.RegisterUserServiceHandler(srv.Server(), handler.NewUserHandler(userService))

	// 启动服务
	if err := srv.Run(); err != nil {
		log.Fatalf("Failed to run service: %v", err)
	}
}