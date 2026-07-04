package auth

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	CreateUser(ctx context.Context, user *User) error
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	GetUserByUsername(ctx context.Context, username string) (*User, error)
	AdminExists(ctx context.Context) (bool, error)
}

type repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) Repository {
	return &repository{db: db}
}

func (r *repository) CreateUser(ctx context.Context, user *User) error {
	query := `INSERT INTO users (email, username, password_hash, role, reference_id)
			  VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at, updated_at`

	err := r.db.QueryRow(ctx, query,
		user.Email, user.Username, user.PasswordHash, user.Role, user.ReferenceID,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

func (r *repository) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	query := `SELECT id, email, username, password_hash, role, reference_id, created_at, updated_at
			  FROM users WHERE email = $1`

	user := &User{}
	err := r.db.QueryRow(ctx, query, email).Scan(
		&user.ID, &user.Email, &user.Username, &user.PasswordHash,
		&user.Role, &user.ReferenceID, &user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	return user, nil
}

func (r *repository) GetUserByUsername(ctx context.Context, username string) (*User, error) {
	query := `SELECT id, email, username, password_hash, role, reference_id, created_at, updated_at
			  FROM users WHERE username = $1`

	user := &User{}
	err := r.db.QueryRow(ctx, query, username).Scan(
		&user.ID, &user.Email, &user.Username, &user.PasswordHash,
		&user.Role, &user.ReferenceID, &user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	return user, nil
}

func (r *repository) AdminExists(ctx context.Context) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE role = 'admin')`

	var exists bool
	err := r.db.QueryRow(ctx, query).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check admin exists: %w", err)
	}

	return exists, nil
}
