package store

import (
	"context"
	"sort"
	"sync"

	"github.com/cfagudelo96/zapatos-be/business/zapato"
)

type InMemoryStore struct {
	lock *sync.RWMutex
	db   map[string]*zapato.Zapato
}

func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{
		lock: &sync.RWMutex{},
		db:   make(map[string]*zapato.Zapato),
	}
}

func (s *InMemoryStore) Create(ctx context.Context, z *zapato.Zapato) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	if _, ok := s.db[z.ID]; ok {
		return zapato.ErrIDUniquenessViolation
	}

	s.db[z.ID] = z
	return nil
}

func (s *InMemoryStore) Get(ctx context.Context, id string) (*zapato.Zapato, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	z, ok := s.db[id]
	if !ok {
		return nil, zapato.ErrNotFound
	}

	return z, nil
}

func (s *InMemoryStore) List(ctx context.Context, filtro zapato.Filtro) ([]*zapato.Zapato, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	var l []*zapato.Zapato
	for _, z := range s.db {
		if cumpleFiltro(z, filtro) {
			l = append(l, z)
		}
	}

	sort.Slice(l, func(i, j int) bool {
		switch filtro.Orden.Atributo {
		case zapato.AtributoCalificacion:
			if filtro.Orden.Orden == zapato.OrdenAsc {
				return l[i].Calificacion < l[j].Calificacion
			}
			return l[i].Calificacion > l[j].Calificacion
		default:
			if filtro.Orden.Orden == zapato.OrdenAsc {
				return l[i].Precio < l[j].Precio
			}
			return l[i].Precio > l[j].Precio
		}
	})

	return l, nil
}

func cumpleFiltro(z *zapato.Zapato, filtro zapato.Filtro) bool {
	if filtro.Categoria != "" && z.Categoria != filtro.Categoria {
		return false
	}
	if len(filtro.Tags) > 0 && !zapatoTieneAlgunTag(z.Tags, filtro.Tags) {
		return false
	}
	return true
}

func zapatoTieneAlgunTag(tz, tf []string) bool {
	for _, t := range tf {
		for _, z := range tz {
			if t == z {
				return true
			}
		}
	}
	return false
}
