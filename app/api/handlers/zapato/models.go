package zapato

import "github.com/cfagudelo96/zapatos-be/business/zapato"

type CreateRequest struct {
	Nombre    string   `json:"nombre" validate:"required,ascii"`
	Precio    float64  `json:"precio" validate:"gt=0"`
	Categoria string   `json:"categoria" validate:"oneof=Femenino Masculino Infantil"`
	Tags      []string `json:"tags" validate:"required,dive,oneof=Deportivo Casual Formal"`
	Imagen    string   `json:"imagen" validate:"required,http_url"`
}

func (r *CreateRequest) ToNewZapato() *zapato.NewZapato {
	return &zapato.NewZapato{
		Nombre:    r.Nombre,
		Precio:    r.Precio,
		Categoria: r.Categoria,
		Tags:      r.Tags,
		Imagen:    r.Imagen,
	}
}

type Filtro struct {
	AtributoOrden string   `query:"atributo" validate:"omitempty,oneof=precio calificacion"`
	Orden         string   `query:"orden" validate:"omitempty,oneof=asc desc"`
	Categoria     string   `query:"categoria" validate:"omitempty,oneof=Femenino Masculino Infantil"`
	Tags          []string `query:"tags" validate:"omitempty,dive,oneof=Deportivo Casual Formal"`
}

func (f Filtro) ToFiltro() zapato.Filtro {
	return zapato.Filtro{
		Orden: zapato.Orden{
			Atributo: f.AtributoOrden,
			Orden:    f.Orden,
		},
		Categoria: f.Categoria,
		Tags:      f.Tags,
	}
}
