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
	GetUserByUsername(ctx context.Context, email string) (dbstore.User, error)
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
	ErrUsernameTaken 	= errors.New("username already in use")
	ErrInvalidCreds 	= errors.New("invalid credentials")
	ErrUserNotFound 	= errors.New("user not found")
)

// Register creates a new user account, hashing the password and checking for duplicates.
func (s *Service) Register(ctx context.Context, req reqdto.RegisterReq) (dbstore.User, error) {
	if _, err := s.store.GetUserByUsername(ctx, req.Username); err == nil {
		return dbstore.User{}, ErrUsernameTaken
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return dbstore.User{}, err
	}

	return s.store.CreateUser(ctx, dbstore.CreateUserParams{
		Username:    req.Username,
		Password: string(hash),
	})
}

// Login gets the user by it's username, checks if the password matches, and returns the user or an error
func (s *Service) Login(ctx context.Context, req reqdto.LoginReq) (dbstore.User, error) {
	u, err := s.store.GetUserByUsername(ctx, req.Username)
	if err != nil {
			return dbstore.User{}, ErrUserNotFound
	}

	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(req.Password))
	if err != nil {
		return dbstore.User{}, ErrInvalidCreds
	}

	return u, nil
}
