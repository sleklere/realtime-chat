// Package auth provides authentication services, including user
// registration, login, and related domain logic.
package auth

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/golang-jwt/jwt/v5"
	reqdto "github.com/sleklere/realtime-chat/cmd/server/internal/api/dto/request"
	resdto "github.com/sleklere/realtime-chat/cmd/server/internal/api/dto/response"
	dbstore "github.com/sleklere/realtime-chat/cmd/server/internal/store"
	"golang.org/x/crypto/bcrypt"
)

// Store defines the persistence methods required for authentication operations.
type Store interface {
	GetUserByUsername(ctx context.Context, username string) (dbstore.User, error)
	CreateUser(ctx context.Context, arg dbstore.CreateUserParams) (dbstore.User, error)
}

// Service provides authentication-related business logic using a Store.
type Service struct {
	store  Store
	logger *slog.Logger
	auth   *Config
}

// NewService creates a new Service with the given Store and logger.
func NewService(s Store, l *slog.Logger, authCfg *Config) *Service {
	return &Service{store: s, logger: l, auth: authCfg}
}

// ErrUsernameTaken indicates that the username is already registered.
// ErrInvalidCreds indicates that the provided credentials are invalid.
// ErrUserNotFound indicates that no user was found in the database.
var (
	ErrUsernameTaken = errors.New("username already in use")
	ErrInvalidCreds  = errors.New("invalid credentials")
	ErrUserNotFound  = errors.New("user not found")
)

// Register creates a new user account, hashing the password and checking for duplicates.
func (s *Service) Register(ctx context.Context, req reqdto.RegisterReq) (resdto.AuthRes, error) {
	if _, err := s.store.GetUserByUsername(ctx, req.Username); err == nil {
		return resdto.AuthRes{}, ErrUsernameTaken
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return resdto.AuthRes{}, err
	}

	u, err := s.store.CreateUser(ctx, dbstore.CreateUserParams{
		Username: req.Username,
		Password: string(hash),
	})
	if err != nil {
		return resdto.AuthRes{}, err
	}

	token, exp, err := s.generateAccessToken(u)
	if err != nil {
		return resdto.AuthRes{}, err
	}

	authRes := resdto.AuthRes{
		User: resdto.UserRes{
			ID:        u.ID,
			Username:  u.Username,
			CreatedAt: u.CreatedAt.Time,
		},
		Token:     token,
		ExpiresAt: exp.Unix(),
	}

	return authRes, nil
}

// Login gets the user by it's username, checks if the password matches, and returns the user or an error
func (s *Service) Login(ctx context.Context, req reqdto.LoginReq) (resdto.AuthRes, error) {
	u, err := s.store.GetUserByUsername(ctx, req.Username)
	if err != nil {
		return resdto.AuthRes{}, ErrUserNotFound
	}

	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(req.Password))
	if err != nil {
		return resdto.AuthRes{}, ErrInvalidCreds
	}

	token, exp, err := s.generateAccessToken(u)
	if err != nil {
		return resdto.AuthRes{}, err
	}

	authRes := resdto.AuthRes{
		User: resdto.UserRes{
			ID:        u.ID,
			Username:  u.Username,
			CreatedAt: u.CreatedAt.Time,
		},
		Token:     token,
		ExpiresAt: exp.Unix(),
	}

	return authRes, nil
}

func (s *Service) generateAccessToken(u dbstore.User) (string, time.Time, error) {
	now := time.Now().UTC()
	exp := now.Add(s.auth.AccessTTL)
	claims := Claims{
		UserID:   u.ID,
		Username: u.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    s.auth.Issuer,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(exp),
		},
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := t.SignedString(s.auth.JWTSecret)
	return token, exp, err
}
