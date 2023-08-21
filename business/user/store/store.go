package store

import (
	"context"
	"fmt"
	"sync"

	"github.com/cfagudelo96/zapatos-be/business/user"
	"golang.org/x/crypto/bcrypt"
)

type InMemoryStore struct {
	lock *sync.RWMutex
	db   map[string]*user.Credentials
}

func NewInMemoryStore() *InMemoryStore {
	s := &InMemoryStore{
		lock: &sync.RWMutex{},
		db:   make(map[string]*user.Credentials),
	}
	return s
}

func (s *InMemoryStore) SignUp(ctx context.Context, u *user.Credentials) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	_, ok := s.db[u.Email]
	if ok {
		return user.ErrAlreadyRegistered
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("hashing password: %w", err)
	}
	u.Password = string(hashedPassword)

	s.db[u.Email] = u
	return nil
}

func (s *InMemoryStore) SignIn(ctx context.Context, c *user.Credentials) (*user.User, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	u, ok := s.db[c.Email]
	if !ok {
		return nil, user.ErrNotFound
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(c.Password)); err != nil {
		return nil, user.ErrWrongCredentials
	}

	return &user.User{
		Email: u.Email,
	}, nil
}
