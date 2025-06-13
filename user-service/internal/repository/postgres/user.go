package postgres

import (
	"context"
	"database/sql"

	"github.com/campus-trading/user-service/internal/model"
	_ "github.com/lib/pq"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(dsn string) (*UserRepository, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	return &UserRepository{db: db}, nil
}

func (r *UserRepository) Close() {
	r.db.Close()
}

func (r *UserRepository) Create(ctx context.Context, user *model.User) error {
	query := `
		INSERT INTO users (user_id, email, phone, password_hash, nickname, avatar_url, contact_info, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	_, err := r.db.ExecContext(ctx, query,
		user.UserID, user.Email, user.Phone, user.PasswordHash, user.Nickname,
		user.AvatarURL, user.ContactInfo, user.CreatedAt, user.UpdatedAt)
	return err
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	query := `SELECT user_id, email, phone, password_hash, nickname, avatar_url, contact_info, created_at, updated_at
	          FROM users WHERE email = $1 AND deleted_at IS NULL`
	user := &model.User{}
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.UserID, &user.Email, &user.Phone, &user.PasswordHash, &user.Nickname,
		&user.AvatarURL, &user.ContactInfo, &user.CreatedAt, &user.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, errors.New("user not found")
	}
	return user, err
}

func (r *UserRepository) GetByID(ctx context.Context, userID string) (*model.User, error) {
	query := `SELECT user_id, email, phone, password_hash, nickname, avatar_url, contact_info, created_at, updated_at
	          FROM users WHERE user_id = $1 AND deleted_at IS NULL`
	user := &model.User{}
	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&user.UserID, &user.Email, &user.Phone, &user.PasswordHash, &user.Nickname,
		&user.AvatarURL, &user.ContactInfo, &user.CreatedAt, &user.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, errors.New("user not found")
	}
	return user, err
}

func (r *UserRepository) Update(ctx context.Context, user *model.User) error {
	query := `
		UPDATE users SET nickname = $2, avatar_url = $3, contact_info = $4, updated_at = $5
		WHERE user_id = $1 AND deleted_at IS NULL
	`
	_, err := r.db.ExecContext(ctx, query,
		user.UserID, user.Nickname, user.AvatarURL, user.ContactInfo, user.UpdatedAt)
	return err
}