package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/Tao-Zzzz/GoCampus/user-service/config"
	"github.com/Tao-Zzzz/GoCampus/user-service/model"
	_ "github.com/lib/pq"
)

// PostgresRepository implements the UserRepository interface for PostgreSQL.
type PostgresRepository struct {
	db *sql.DB
}

// NewPostgresRepository creates a new PostgresRepository instance.
func NewPostgresRepository(cfg *config.Config) (*PostgresRepository, error) {
	db, err := sql.Open(cfg.Database.Driver, cfg.Database.GetDSN())
	if err != nil {
		return nil, err
	}
	// Verify connection
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return &PostgresRepository{db: db}, nil
}


// CreateUser inserts a new user into the database.
func (r *PostgresRepository) CreateUser(ctx context.Context, user *model.User) (string, error) {
	query := `
		INSERT INTO users (id, email, password, nickname, avatar, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`
	var userID string
	err := r.db.QueryRowContext(ctx, query, user.ID, user.Email, user.Password, user.Nickname, user.Avatar, user.CreatedAt).Scan(&userID)
	if err != nil {
		if err.Error() == "pq: duplicate key value violates unique constraint \"users_email_key\"" {
			return "", errors.New("user already exists")
		}
		return "", err
	}
	return userID, nil
}

// GetUserByEmail retrieves a user by their email address.
func (r *PostgresRepository) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	query := `
		SELECT id, email, password, nickname, avatar, created_at
		FROM users
		WHERE email = $1
	`
	user := &model.User{}
	err := r.db.QueryRowContext(ctx, query, email).Scan(&user.ID, &user.Email, &user.Password, &user.Nickname, &user.Avatar, &user.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, errors.New("user not found")
	}
	if err != nil {
		return nil, err
	}
	return user, nil
}

// GetUserByID retrieves a user by their ID.
func (r *PostgresRepository) GetUserByID(ctx context.Context, id string) (*model.User, error) {
	query := `
		SELECT id, email, password, nickname, avatar, created_at
		FROM users
		WHERE id = $1
	`
	user := &model.User{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(&user.ID, &user.Email, &user.Password, &user.Nickname, &user.Avatar, &user.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, errors.New("user not found")
	}
	if err != nil {
		return nil, err
	}
	return user, nil
}