package zapato

import (
	"context"
	"errors"
	"fmt"
)

var (
	ErrIDUniquenessViolation = errors.New("id uniqueness violation")
	ErrNotFound              = errors.New("zapato not found")
)

type Store interface {
	Create(ctx context.Context, z *Zapato) error
	Get(ctx context.Context, id string) (*Zapato, error)
	List(ctx context.Context, filtro Filtro) ([]*Zapato, error)
}

type Service struct {
	store Store
}

func NewService(s Store) *Service {
	return &Service{store: s}
}

func (s *Service) Create(ctx context.Context, nz *NewZapato) (*Zapato, error) {
	z := nz.ToZapato()
	if err := s.store.Create(ctx, z); err != nil {
		return nil, fmt.Errorf("creating zapato in store: %w", err)
	}
	return z, nil
}

func (s *Service) Get(ctx context.Context, id string) (*Zapato, error) {
	z, err := s.store.Get(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("getting zapato from store: %w", err)
	}
	return z, nil
}

func (s *Service) List(ctx context.Context, f Filtro) (*ListaZapatos, error) {
	f = f.WithDefaults()
	l, err := s.store.List(ctx, f)
	if err != nil {
		return nil, fmt.Errorf("listing zapatos from store: %w", err)
	}
	return &ListaZapatos{
		Data: l,
		Metadata: Metadata{
			Filtro: f,
		},
	}, nil
}
