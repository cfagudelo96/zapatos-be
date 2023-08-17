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
	s := &InMemoryStore{
		lock: &sync.RWMutex{},
		db:   make(map[string]*zapato.Zapato),
	}
	s.addInitialData()
	return s
}

func (s *InMemoryStore) addInitialData() {
	ctx := context.Background()
	l := []*zapato.NewZapato{
		{
			Nombre:    "NY 90 W",
			Precio:    399.950,
			Categoria: zapato.CategoriaFemenino,
			Tags: []string{
				zapato.TagCasual, zapato.TagDeportivo,
			},
			Imagen: "https://assets.adidas.com/images/w_766,h_766,f_auto,q_auto,fl_lossy,c_fill,g_auto/1b9dff65bee740de8529aeec0091c760_9366/ny-90-w.jpg",
		},
		{
			Nombre:    "Duramo RC Shoes",
			Precio:    270.950,
			Categoria: zapato.CategoriaFemenino,
			Tags: []string{
				zapato.TagDeportivo,
			},
			Imagen: "https://assets.adidas.com/images/w_766,h_766,f_auto,q_auto,fl_lossy,c_fill,g_auto/e2b338361b024762a9b5ff8f93689853_9366/duramo-rc-shoes.jpg",
		},
		{
			Nombre:    "Response Runner Shoes",
			Precio:    229.950,
			Categoria: zapato.CategoriaFemenino,
			Tags: []string{
				zapato.TagDeportivo,
			},
			Imagen: "https://assets.adidas.com/images/w_766,h_766,f_auto,q_auto,fl_lossy,c_fill,g_auto/fbd76634f84e4dfd8154bddbd72ba7c3_9366/response-runner-shoes.jpg",
		},
		{
			Nombre:    "Yzy Knit RNR",
			Precio:    899.950,
			Categoria: zapato.CategoriaMasculino,
			Tags: []string{
				zapato.TagDeportivo,
			},
			Imagen: "https://assets.adidas.com/images/w_766,h_766,f_auto,q_auto,fl_lossy,c_fill,g_auto/39c8469426b04c17bd4da7f42e89a6c8_9366/yzy-knit-rnr.jpg",
		},
		{
			Nombre:    "Tenis Samba OG",
			Precio:    599.950,
			Categoria: zapato.CategoriaMasculino,
			Tags: []string{
				zapato.TagFormal,
			},
			Imagen: "https://assets.adidas.com/images/w_766,h_766,f_auto,q_auto,fl_lossy,c_fill,g_auto/5c7f0a2ca4c5450fbea8afb200f9f923_9366/tenis-samba-og.jpg",
		},
		{
			Nombre:    "Tenis Questar Flow NXT",
			Precio:    159.271,
			Categoria: zapato.CategoriaInfantil,
			Tags:      []string{zapato.TagCasual, zapato.TagFormal},
			Imagen:    "https://assets.adidas.com/images/w_766,h_766,f_auto,q_auto,fl_lossy,c_fill,g_auto/91ef876b56de40d281d4ac3700aacd0b_9366/tenis-questar-flow-nxt.jpg",
		},
		{
			Nombre:    "Tenis Monofit Slip-On",
			Precio:    219.50,
			Categoria: zapato.CategoriaInfantil,
			Tags:      []string{zapato.TagCasual, zapato.TagDeportivo},
			Imagen:    "https://assets.adidas.com/images/w_766,h_766,f_auto,q_auto,fl_lossy,c_fill,g_auto/ee931bc180644ea6bdc0ece249b48176_9366/tenis-monofit-slip-on.jpg",
		},
	}
	for _, nz := range l {
		z := nz.ToZapato()
		if err := s.Create(ctx, z); err != nil {
			panic(err)
		}
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

func (s *InMemoryStore) AddComment(ctx context.Context, c *zapato.Comentario) (*zapato.Zapato, error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	z, ok := s.db[c.ZapatoID]
	if !ok {
		return nil, zapato.ErrNotFound
	}

	z.Calificacion = (float64(len(z.Comentarios))*z.Calificacion + float64(c.Calificacion)) / float64(len(z.Comentarios)+1)
	z.Comentarios = append(z.Comentarios, c)
	s.db[z.ID] = z

	return z, nil
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

func (s *InMemoryStore) Delete(ctx context.Context, id string) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	if _, ok := s.db[id]; !ok {
		return zapato.ErrNotFound
	}

	delete(s.db, id)
	return nil
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
