package repository

import (
"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/Tao-Zzzz/GoCampus/user-service/config"
	"github.com/Tao-Zzzz/GoCampus/user-service/model"
	"github.com/Tao-Zzzz/GoCampus/user-service/pkg/logger"
	"github.com/mattn/go-sqlite3"
)

// setupTestDB creates the users table in the PostgreSQL database.
func setupTestDB(t *testing.T, db *sql.DB) {
	// Create users table
	_, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS users (
		id TEXT PRIMARY KEY,
		email TEXT UNIQUE NOT NULL,
		password TEXT NOT NULL,
		nickname TEXT NOT NULL,
		avatar TEXT,
		created_at TIMESTAMP
	)
	`)
	if err != nil {
		t.Fatalf("Failed to create users table: %v", err)
	}
}

func cleanupTestDB(t *testing.T, db *sql.DB) {
    _, err := db.Exec("DELETE FROM users")
    if err != nil {
        t.Fatalf("Failed to clean up users table: %v", err)
    }
}

func TestPostgresRepository_CreateUser(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}
	defer db.Close()

	// Create schema
	_, err = db.Exec("CREATE TABLE users (id TEXT, email TEXT, password TEXT, nickname TEXT, avatar TEXT, created_at TIMESTAMP)")
	if err != nil {
		t.Fatalf("Failed to create table: %v", err)
	}

	cfg := &config.Config{
		Database: config.DatabaseConfig{
			Driver: "sqlite3",
		},
		Service: config.ServiceConfig		{
			Name:     "test-service",
			Port:      8080,
			LogLevel:   "debug",
		},
	}
	log := logger.NewLogger(cfg)
	repo := &PostgresRepository{db: db, logger: log, tracer: otel.Tracer("test-postgres-repository")}

	tests := []struct {
		name    string
		user    *model.User
		wantErr bool
	}{
		{
			name: "Successful create user",
			user: &model.User{
				ID:       "user123",
				Email:    "test@example.com",
				Password: "hashedpassword",
				Nickname: "TestUser",
				Avatar:   "http://example.com/avatar.png",
				CreatedAt: time.Now(),
			},
			wantErr: false,
		},
		{
			name: "Duplicate email",
			user: &model.User{
				ID:       "user456",
				Email:    "test@example.com",
				Password: "hashedpassword",
				Nickname: "TestUser2",
				Avatar:   "http://example.com/avatar2.png",
				CreatedAt: time.Now(),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userID, err := repo.CreateUser(context.Background(), tt.user)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateUser() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && userID != tt.user.ID {
				t.Errorf("CreateUser() userID = %v, want %v", userID, tt.user.ID)
			}
		})
	}
}

func TestPostgresRepository_GetUserByEmail(t *testing.T) {
	db, _ := sql.Open("sqlite3", ":memory:")
	_, _ := db.Exec("CREATE TABLE users (id TEXT, email TEXT, password TEXT, nickname TEXT, nickname TEXT, avatar TEXT, created_at TIMESTAMP)")
	cfg := &config.Config{
		Database: config.DatabaseConfig{
			Driver: "sqlite3",
		},
		Service: config.ServiceConfig{
			Name:     "test-service",
			Port:     8080,
			LogLevel: "debug",
		},
	}
	log := logger.NewLogger(cfg)
	repo := &PostgresRepository{db: db, logger: log, tracer: otel.Tracer("test-postgres-repository")}

	// Insert a test user
	user := &model.User{
		ID:       "user123",
		Email:    "test@example.com",
		Password: "hashedpassword",
		Nickname: "TestUser",
		Avatar:   "http://example.com/avatar.png",
		CreatedAt: time.Now().UTC().Truncate(time.Millisecond),
	}
	_, err = repo.CreateUser(context.Background(), user)
	if err != nil {
		t.Fatalf("Failed to insert test user: %v", err)
	}

	tests := []struct {
		name     string
		email    string
		wantErr  bool
		wantUser *model.User
	}{
		{
			name:     "Successful get user by email",
			email:    "test@example.com",
			wantErr:  false,
			wantUser: user,
		},
		{
			name:     "User not found",
			email:    "invalid@example.com",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotUser, err := repo.GetUserByEmail(context.Background(), tt.email)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUserByEmail() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && (gotUser.ID != tt.wantUser.ID || gotUser.Email != tt.wantUser.Email ||
				gotUser.Password != tt.wantUser.Password || gotUser.Nickname != tt.wantUser.Nickname ||
				gotUser.Avatar != tt.wantUser.Avatar || !gotUser.CreatedAt.UTC().Equal(tt.wantUser.CreatedAt.UTC())) {
				t.Errorf("GetUserByEmail() user = %+v, want %+v", gotUser, tt.wantUser)
			}
		})
	}
}

func TestPostgresRepository_GetUserByID(t *testing.T) {
	cfg := &config.Config{
		Database: config.DatabaseConfig{
			Driver:   "postgres",
			Host:     "localhost",
			Port:     5432,
			User:     "postgres",       // 替换为您的数据库用户名
			Password: "123456",    // 替换为您的数据库密码
			DBName:   "users",    // 替换为您的数据库名称
			SSLMode:  "disable",
		},
	}

	db, err := sql.Open(cfg.Database.Driver, cfg.Database.GetDSN())
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close() // 测试结束后关闭数据库连接

	// Setup the database schema
	setupTestDB(t, db)
	cleanupTestDB(t, db)
	repo, err := NewPostgresRepository(cfg)
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}

	// Insert a test user
	user := &model.User{
		ID:       "user123",
		Email:    "test@example.com",
		Password: "hashedpassword",
		Nickname: "TestUser",
		Avatar:   "http://example.com/avatar.png",
		CreatedAt: time.Now().UTC().Truncate(time.Millisecond),
	}
	_, err = repo.CreateUser(context.Background(), user)
	if err != nil {
		t.Fatalf("Failed to insert test user: %v", err)
	}

	tests := []struct {
		name     string
		id       string
		wantErr  bool
		wantUser *model.User
	}{
		{
			name:     "Successful get user by ID",
			id:       "user123",
			wantErr:  false,
			wantUser: user,
		},
		{
			name:     "User not found",
			id:       "invalid",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotUser, err := repo.GetUserByID(context.Background(), tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUserByID() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && (gotUser.ID != tt.wantUser.ID || gotUser.Email != tt.wantUser.Email ||
				gotUser.Password != tt.wantUser.Password || gotUser.Nickname != tt.wantUser.Nickname ||
				gotUser.Avatar != tt.wantUser.Avatar || !gotUser.CreatedAt.UTC().Equal(tt.wantUser.CreatedAt.UTC())) {
				t.Errorf("GetUserByID() user = %+v, want %+v", gotUser, tt.wantUser)
			}
		})
	}
}
