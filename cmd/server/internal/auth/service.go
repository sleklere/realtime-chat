// Package auth provides authentication services, including user
// registration, login, and related domain logic.
package auth

import (
	"context"
	"errors"
	"log/slog"

	reqdto "github.com/sleklere/realtime-chat/cmd/server/internal/api/dto/request"
	dbstore "github.com/sleklere/realtime-chat/cmd/server/internal/store"
	"golang.org/x/crypto/bcrypt"
)

// Store defines the persistence methods required for authentication operations.
type Store interface {
	GetUserByEmail(ctx context.Context, email string) (dbstore.User, error)
	CreateUser(ctx context.Context, arg dbstore.CreateUserParams) (dbstore.User, error)
}

// Service provides authentication-related business logic using a Store.
type Service struct {
	store  Store
	logger *slog.Logger
}

// NewService creates a new Service with the given Store and logger.
func NewService(s Store, l *slog.Logger) *Service {
	return &Service{store: s, logger: l}
}

// ErrEmailTaken indicates that the email is already registered.
// ErrInvalidCreds indicates that the provided credentials are invalid.
var (
	ErrEmailTaken   = errors.New("email already in use")
	ErrInvalidCreds = errors.New("invalid credentials")
)

// Register creates a new user account, hashing the password and checking for duplicates.
func (s *Service) Register(ctx context.Context, req reqdto.RegisterReq) (dbstore.User, error) {
	if _, err := s.store.GetUserByEmail(ctx, req.Email); err == nil {
		return dbstore.User{}, ErrEmailTaken
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return dbstore.User{}, err
	}

	return s.store.CreateUser(ctx, dbstore.CreateUserParams{
		Email:    req.Email,
		Password: string(hash),
	})
}
