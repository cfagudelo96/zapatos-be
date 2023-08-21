package user

import (
	"context"
	"errors"
	"fmt"
)

var (
	ErrNotFound          = errors.New("user not found")
	ErrWrongCredentials  = errors.New("wrong credentials")
	ErrAlreadyRegistered = errors.New("email already registered")
)

type Credentials struct {
	Email    string `json:"email" validate:"email"`
	Password string `json:"password" validate:"min=8"`
}

func (c Credentials) ToUser() *User {
	return &User{
		Email: c.Email,
	}
}

type User struct {
	Email string `json:"email"`
}

type Store interface {
	SignUp(context.Context, *Credentials) error
	SignIn(context.Context, *Credentials) (*User, error)
}

type Service struct {
	store Store
}

func NewService(s Store) *Service {
	return &Service{store: s}
}

func (s *Service) SignUp(ctx context.Context, c *Credentials) (*User, error) {
	if err := s.store.SignUp(ctx, c); err != nil {
		return nil, fmt.Errorf("storing user credentials: %w", err)
	}
	return c.ToUser(), nil
}

func (s *Service) SignIn(ctx context.Context, c *Credentials) (*User, error) {
	u, err := s.store.SignIn(ctx, c)
	if err != nil {
		return nil, err
	}
	return u, nil
}
