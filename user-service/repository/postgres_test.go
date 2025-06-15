package repository

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/Tao-Zzzz/GoCampus/user-service/model"
	_ "github.com/mattn/go-sqlite3"
)

func setupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}

	// Create users table
	_, err = db.Exec(`
		CREATE TABLE users (
			id TEXT PRIMARY KEY,
			email TEXT UNIQUE NOT NULL,
			password TEXT NOT NULL,
			nickname TEXT NOT NULL,
			avatar TEXT,
			created_at DATETIME
		)
	`)
	if err != nil {
		t.Fatalf("Failed to create users table: %v", err)
	}

	return db
}

func TestPostgresRepository_CreateUser(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	repo := NewPostgresRepository(db)

	tests := []struct {
		name    string
		user    *model.User
		wantErr bool
	}{
		{
			name: "Successful create user",
			user: &model.User{
				ID:        "user123",
				Email:     "test@example.com",
				Password:  "hashedpassword",
				Nickname:  "TestUser",
				Avatar:    "http://example.com/avatar.png",
				CreatedAt: time.Now(),
			},
			wantErr: false,
		},
		{
			name: "Duplicate email",
			user: &model.User{
				ID:        "user456",
				Email:     "test@example.com",
				Password:  "hashedpassword",
				Nickname:  "TestUser2",
				Avatar:    "http://example.com/avatar2.png",
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
	db := setupTestDB(t)
	defer db.Close()
	repo := NewPostgresRepository(db)

	// Insert a test user
	user := &model.User{
		ID:        "user123",
		Email:     "test@example.com",
		Password:  "hashedpassword",
		Nickname:  "TestUser",
		Avatar:    "http://example.com/avatar.png",
		CreatedAt: time.Now().Truncate(time.Millisecond),
	}
	_, err := repo.CreateUser(context.Background(), user)
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
			name:    "Successful get user by email",
			email:   "test@example.com",
			wantErr: false,
			wantUser: user,
		},
		{
			name:    "User not found",
			email:   "invalid@example.com",
			wantErr: true,
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
				gotUser.Avatar != tt.wantUser.Avatar || gotUser.CreatedAt.Unix() != tt.wantUser.CreatedAt.Unix()) {
				t.Errorf("GetUserByEmail() user = %+v, want %+v", gotUser, tt.wantUser)
			}
		})
	}
}

func TestPostgresRepository_GetUserByID(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	repo := NewPostgresRepository(db)

	// Insert a test user
	user := &model.User{
		ID:        "user123",
		Email:     "test@example.com",
		Password:  "hashedpassword",
		Nickname:  "TestUser",
		Avatar:    "http://example.com/avatar.png",
		CreatedAt: time.Now().Truncate(time.Millisecond),
	}
	_, err := repo.CreateUser(context.Background(), user)
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
			name:    "Successful get user by ID",
			id:      "user123",
			wantErr: false,
			wantUser: user,
		},
		{
			name:    "User not found",
			id:      "invalid",
			wantErr: true,
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
				gotUser.Avatar != tt.wantUser.Avatar || gotUser.CreatedAt.Unix() != tt.wantUser.CreatedAt.Unix()) {
				t.Errorf("GetUserByID() user = %+v, want %+v", gotUser, tt.wantUser)
			}
		})
	}
}