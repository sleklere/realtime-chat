// Package user contains domain logic
package user

import (
	"context"
	"errors"

	reqdto "github.com/sleklere/realtime-chat/cmd/server/internal/api/dto/request"
	dbstore "github.com/sleklere/realtime-chat/cmd/server/internal/store"
	"golang.org/x/crypto/bcrypt"
)

type Store interface {
	GetUserByEmail(ctx context.Context, email string) (dbstore.User, error)
	CreateUser(ctx context.Context, arg dbstore.CreateUserParams) (dbstore.User, error)
}

type Service struct {
	store Store
}

func NewService(s Store) *Service {
	return &Service{store: s}
}

var (
	ErrEmailTaken   = errors.New("email already in use")
	ErrInvalidCreds = errors.New("invalid credentials")
)

func (s *Service) Register(ctx context.Context, req reqdto.RegisterReq) (dbstore.User, error) {
	if _, err := s.store.GetUserByEmail(ctx, req.Email); err == nil {
		return dbstore.User{}, ErrEmailTaken
	}
	hash, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	return s.store.CreateUser(ctx, dbstore.CreateUserParams{
		Email: req.Email,
		Password: string(hash),
	})
}

