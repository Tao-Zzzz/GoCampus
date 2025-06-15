package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/Tao-Zzzz/GoCampus/user-service/config"
	"github.com/Tao-Zzzz/GoCampus/user-service/model"
	"github.com/Tao-Zzzz/GoCampus/user-service/pkg/logger"
	_ "github.com/lib/pq"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

// PostgresRepository implements UserRepository using PostgreSQL.
type PostgresRepository struct {
	db     *sql.DB
	logger *logger.Logger
	tracer otel.Tracer
}

// NewPostgresRepository creates a new PostgresRepository instance.
func NewPostgresRepository(ctx context.Context,
	cfg *config.Config,
	log *logger.Logger
) (*PostgresRepository, error) {
	db, err := sql.Open(cfg.Database.Driver, cfg.Database.GetDSN())
	if err != nil {
		log.Error(ctx).Error(err).Msg("Failed to open database connection")
		return nil, err
	}
	// Verify connection
	if err := db.PingContext(ctx); err != nil {
		log.Error(ctx).Err(err).Msg("Failed to ping database")
		return nil, err
	}
	return &PostgresRepository{
		db:     db,
		logger: log,
		tracer: otel.Tracer("postgres-repository"),
	}, nil
}

// CreateUserService creates a new user in the database.
func (r *PostgresRepository) CreateUser(ctx context.Context, user *model.User) (string, error) {
	ctx, span := r.tracer.Start(ctx, "PostgresRepository.CreateUser")
	defer span.End()

	r.logger.Info(ctx).Msgf("Creating user with email: %s", user.Email)

	query := "INSERT INTO users (id, email, password, nickname, avatar, created_at) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id"
	_, err := r.db.ExecContext(ctx, query, user.ID, user.Email, user.Password, user.Nickname, user.Avatar, user.CreatedAt)
	if err != nil {
		r.logger.Error(ctx).Err(err).Msg("Failed to create user in database")
		span.RecordError(err)
		return "", errors.New("failed to create user")
	}

	r.logger.Info(ctx).Msgf("User created successfully: %s", user.ID)
	span.SetAttributes(attribute.String("user_id", user.ID))
	return user.ID, nil
}

// GetUserByEmail retrieves a user by email.
func (r *PostgresRepository) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	ctx, span := r.tracer.Start(ctx, "PostgresRepository.GetUserByEmail")
	defer span.End()

	r.logger.Info(ctx).Msgf("Retrieving user by email: %s", email)

	user := &model.User{}
	query := "SELECT id, email, password, nickname, avatar, created_at FROM users WHERE email = $1"
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.Nickname,
		&user.Avatar,
		&user.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			r.logger.Warn(ctx).Msg("User not found by email")
			return nil, errors.New("user not found")
		}
		r.logger.Error(ctx).Err(err).Msg("Failed to retrieve user by email")
		span.RecordError(err)
		return nil, errors.New("failed to get user")
	}

	r.logger.Info(ctx).Msgf("User retrieved by email: %s", email)
	span.SetAttributes(attribute.String("user_id", user.ID))
	return user, nil
}

// GetUserByID retrieves a user by ID.
func (r *PostgresRepository) GetUserByID(ctx context.Context, id string) (*model.User, error) {
	ctx, span := r.tracer.Start(ctx, "PostgresRepository.GetUserByID")
	defer span.End()

	r.logger.Info(ctx).Msgf("Retrieving user by ID: %s", id)

	user := &model.User{}
	query := "SELECT id, email, password, nickname, avatar, created_at FROM users WHERE id = $1"
	err := r.db.QueryRowContext(ctx, query, id).Scan(
			&user.ID,
		&user.Email,
		&user.Password,
		&user.Nickname,
		&user.Avatar,
		&user.CreatedAt,
		)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			r.logger.Warn(ctx).Msg("User not found by ID")
			return nil, errors.New("user not found")
		}
		r.logger.Error(ctx).Err(err).Msg("Failed to retrieve user by ID")
		span.RecordError(err)
		return nil, errors.New("failed to get user")
	}

	r.logger.Info(ctx).Msgf("User retrieved by ID: %s", id)
	span.SetAttributes(attribute.String("user_id", id))
	return user, nil
}