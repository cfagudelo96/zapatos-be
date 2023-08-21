package user

import (
	"net/http"
	"time"

	"github.com/cfagudelo96/zapatos-be/business/user"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type Claims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

type Token struct {
	Token string `json:"token"`
}

type Handler struct {
	jwtSecret string
	service   *user.Service
}

func NewHandler(s *user.Service, secret string) *Handler {
	return &Handler{jwtSecret: secret, service: s}
}

func (h *Handler) SignIn(c echo.Context) error {
	ctx := c.Request().Context()
	var r user.Credentials
	if err := c.Bind(&r); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err := c.Validate(&r); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	u, err := h.service.SignIn(ctx, &r)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	claims := &Claims{
		Email: u.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString([]byte(h.jwtSecret))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, &Token{Token: t})
}

func (h *Handler) SignUp(c echo.Context) error {
	ctx := c.Request().Context()
	var r user.Credentials
	if err := c.Bind(&r); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err := c.Validate(&r); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	u, err := h.service.SignUp(ctx, &r)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusCreated, u)
}
