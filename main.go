package main

import (
	"os"

	"github.com/cfagudelo96/zapatos-be/app/api/handlers"
	userhdlr "github.com/cfagudelo96/zapatos-be/app/api/handlers/user"
	zapatohdlr "github.com/cfagudelo96/zapatos-be/app/api/handlers/zapato"
	"github.com/cfagudelo96/zapatos-be/business/user"
	userstore "github.com/cfagudelo96/zapatos-be/business/user/store"
	"github.com/cfagudelo96/zapatos-be/business/zapato"
	zapatostore "github.com/cfagudelo96/zapatos-be/business/zapato/store"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const secret = "93A4BA417F16957C8191C32B593E6" // This should be stored in an env variable.

type server struct {
	e             *echo.Echo
	jwtSecret     string
	jwtMiddleware echo.MiddlewareFunc
}

func newServer() *server {
	e := echo.New()
	e.Validator = handlers.NewValidator()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	s := &server{
		e:             e,
		jwtSecret:     secret,
		jwtMiddleware: echojwt.JWT([]byte(secret)),
	}
	s.addUserRoutes()
	s.addZapatosRoutes()

	return s
}

func main() {
	s := newServer()

	port := "8080"
	if p := os.Getenv("PORT"); p != "" {
		port = p
	}

	s.e.Logger.Fatal(s.e.Start(":" + port))
}

func (s *server) addUserRoutes() {
	ustore := userstore.NewInMemoryStore()
	uservice := user.NewService(ustore)
	uhandler := userhdlr.NewHandler(uservice, s.jwtSecret)

	s.e.POST("sign-in", uhandler.SignIn)
	s.e.POST("sign-up", uhandler.SignUp)
}

func (s *server) addZapatosRoutes() {
	zstore := zapatostore.NewInMemoryStore()
	zservice := zapato.NewService(zstore)
	zhandler := zapatohdlr.NewHandler(zservice)

	g := s.e.Group("/zapatos")
	g.GET("", zhandler.List)
	g.GET("/:id", zhandler.Get)
	g.POST("", zhandler.Create, s.jwtMiddleware)
	g.DELETE("/:id", zhandler.Delete, s.jwtMiddleware)
	g.POST("/:id/comment", zhandler.AddComment, s.jwtMiddleware)
}
