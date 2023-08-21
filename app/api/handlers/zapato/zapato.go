package zapato

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/cfagudelo96/zapatos-be/business/zapato"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	service *zapato.Service
}

func NewHandler(s *zapato.Service) *Handler {
	return &Handler{service: s}
}

func (h *Handler) Create(c echo.Context) error {
	ctx := c.Request().Context()
	var r CreateRequest
	if err := c.Bind(&r); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err := c.Validate(&r); err != nil {
		return err
	}
	z, err := h.service.Create(ctx, r.ToNewZapato())
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusCreated, z)
}

func (h *Handler) AddComment(c echo.Context) error {
	ctx := c.Request().Context()
	id := c.Param("id")
	var r AddCommentRequest
	if err := c.Bind(&r); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err := c.Validate(&r); err != nil {
		return err
	}
	z, err := h.service.AddComment(ctx, id, r.ToNewComentario())
	if err != nil {
		if errors.Is(err, zapato.ErrNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("El zapato con el ID %q no fue encontrado", id))
		}
		return err
	}
	return c.JSON(http.StatusOK, z)
}

func (h *Handler) Get(c echo.Context) error {
	ctx := c.Request().Context()
	id := c.Param("id")
	z, err := h.service.Get(ctx, id)
	if err != nil {
		if errors.Is(err, zapato.ErrNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("El zapato con el ID %q no fue encontrado", id))
		}
		return err
	}
	return c.JSON(http.StatusOK, z)
}

func (h *Handler) List(c echo.Context) error {
	ctx := c.Request().Context()
	var f Filtro
	if err := c.Bind(&f); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err := c.Validate(&f); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	l, err := h.service.List(ctx, f.ToFiltro())
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, l)
}

func (h *Handler) Update(c echo.Context) error {
	return nil
}

func (h *Handler) Delete(c echo.Context) error {
	ctx := c.Request().Context()
	id := c.Param("id")
	if err := h.service.Delete(ctx, id); err != nil {
		if errors.Is(err, zapato.ErrNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("El zapato con el ID %q no fue encontrado", id))
		}
		return err
	}
	return c.NoContent(http.StatusNoContent)
}
