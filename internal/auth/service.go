package auth

import (
	"context"
	"errors"
	"fmt"
)

var ErrDuplicateEmail = errors.New("user with this email already exists")
var ErrDuplicateUsername = errors.New("user with this username already exists")
var ErrInvalidCredentials = errors.New("invalid email or password")
var ErrNotFound = errors.New("user not found")

type Service interface {
	Login(ctx context.Context, email, username, password string) (*LoginResponse, error)
	CreateUser(ctx context.Context, email, username, password, role string, referenceID int64) (*LoginResponse, error)
	SeedAdmin(ctx context.Context, username, password, email string) error
}

type service struct {
	repo      Repository
	jwtSecret string
}

func NewService(repo Repository, jwtSecret string) Service {
	return &service{repo: repo, jwtSecret: jwtSecret}
}

func (s *service) Login(ctx context.Context, email, username, password string) (*LoginResponse, error) {
	var user *User
	var err error

	if email != "" {
		user, err = s.repo.GetUserByEmail(ctx, email)
	} else {
		user, err = s.repo.GetUserByUsername(ctx, username)
	}
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	if !checkPassword(password, user.PasswordHash) {
		return nil, ErrInvalidCredentials
	}

	token, err := generateJWT(Claims{
		UserID:      user.ID,
		Email:       user.Email,
		Role:        user.Role,
		ReferenceID: safeRefID(user.ReferenceID),
	}, s.jwtSecret)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &LoginResponse{
		Token:       token,
		UserID:      user.ID,
		Role:        user.Role,
		ReferenceID: safeRefID(user.ReferenceID),
	}, nil
}

func (s *service) CreateUser(ctx context.Context, email, username, password, role string, referenceID int64) (*LoginResponse, error) {
	existing, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil && !errors.Is(err, ErrNotFound) {
		return nil, fmt.Errorf("failed to check existing user: %w", err)
	}
	if existing != nil {
		return nil, ErrDuplicateEmail
	}

	existing, err = s.repo.GetUserByUsername(ctx, username)
	if err != nil && !errors.Is(err, ErrNotFound) {
		return nil, fmt.Errorf("failed to check existing username: %w", err)
	}
	if existing != nil {
		return nil, ErrDuplicateUsername
	}

	hash, err := hashPassword(password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	var ref *int64
	if referenceID != 0 {
		ref = &referenceID
	}

	user := &User{
		Email:        email,
		Username:     username,
		PasswordHash: hash,
		Role:         role,
		ReferenceID:  ref,
	}

	if err := s.repo.CreateUser(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	token, err := generateJWT(Claims{
		UserID:      user.ID,
		Email:       user.Email,
		Role:        user.Role,
		ReferenceID: safeRefID(user.ReferenceID),
	}, s.jwtSecret)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &LoginResponse{
		Token:       token,
		UserID:      user.ID,
		Role:        user.Role,
		ReferenceID: safeRefID(user.ReferenceID),
	}, nil
}

func (s *service) SeedAdmin(ctx context.Context, username, password, email string) error {
	exists, err := s.repo.AdminExists(ctx)
	if err != nil {
		return fmt.Errorf("failed to check admin exists: %w", err)
	}

	if exists {
		return nil
	}

	hash, err := hashPassword(password)
	if err != nil {
		return fmt.Errorf("failed to hash admin password: %w", err)
	}

	user := &User{
		Email:        email,
		Username:     username,
		PasswordHash: hash,
		Role:         "admin",
	}

	if err := s.repo.CreateUser(ctx, user); err != nil {
		return fmt.Errorf("failed to seed admin user: %w", err)
	}

	return nil
}

func safeRefID(ref *int64) int64 {
	if ref == nil {
		return 0
	}
	return *ref
}
