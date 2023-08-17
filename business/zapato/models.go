package zapato

import "github.com/google/uuid"

const (
	AtributoPrecio       = "precio"
	AtributoCalificacion = "calificacion"

	OrdenAsc  = "asc"
	OrdenDesc = "desc"

	CategoriaFemenino  = "Femenino"
	CategoriaMasculino = "Masculino"
	CategoriaInfantil  = "Infantil"

	TagDeportivo = "Deportivo"
	TagCasual    = "Casual"
	TagFormal    = "Formal"
)

type Zapato struct {
	ID           string        `json:"id,omitempty"`
	Nombre       string        `json:"nombre,omitempty"`
	Precio       float64       `json:"precio,omitempty"`
	Categoria    string        `json:"categoria,omitempty"`
	Tags         []string      `json:"tags,omitempty"`
	Imagen       string        `json:"imagen,omitempty"`
	Calificacion float64       `json:"calificacion,omitempty"`
	Comentarios  []*Comentario `json:"comentarios,omitempty"`
}

type NewZapato struct {
	Nombre    string
	Precio    float64
	Categoria string
	Tags      []string
	Imagen    string
}

func (nz NewZapato) ToZapato() *Zapato {
	return &Zapato{
		ID:           uuid.NewString(),
		Nombre:       nz.Nombre,
		Precio:       nz.Precio,
		Categoria:    nz.Categoria,
		Tags:         nz.Tags,
		Imagen:       nz.Imagen,
		Comentarios:  nil,
		Calificacion: 0.0,
	}
}

type ListaZapatos struct {
	Data     []*Zapato `json:"data"`
	Metadata Metadata  `json:"metadata"`
}

type Metadata struct {
	Filtro Filtro `json:"filtro"`
}

type NewComentario struct {
	Calificacion uint8  `json:"calificacion,omitempty"`
	Autor        string `json:"autor,omitempty"`
	Comentario   string `json:"comentario,omitempty"`
}

func (n NewComentario) ToComentario(zID string) *Comentario {
	return &Comentario{
		ID:           uuid.NewString(),
		ZapatoID:     zID,
		Calificacion: n.Calificacion,
		Autor:        n.Autor,
		Comentario:   n.Comentario,
	}
}

type Comentario struct {
	ID           string `json:"id,omitempty"`
	ZapatoID     string `json:"zapatoID,omitempty"`
	Calificacion uint8  `json:"calificacion,omitempty"`
	Autor        string `json:"autor,omitempty"`
	Comentario   string `json:"comentario,omitempty"`
}

type Orden struct {
	Atributo string `json:"atributo,omitempty"`
	Orden    string `json:"orden,omitempty"`
}

type Filtro struct {
	Orden     Orden    `json:"orden"`
	Categoria string   `json:"categoria,omitempty"`
	Tags      []string `json:"tags,omitempty"`
}

func (f Filtro) WithDefaults() Filtro {
	if f.Orden.Orden == "" {
		f.Orden.Orden = OrdenAsc
	}
	if f.Orden.Atributo == "" {
		f.Orden.Atributo = AtributoPrecio
	}
	return f
}
