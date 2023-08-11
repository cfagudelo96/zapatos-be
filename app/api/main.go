package main

import (
	"github.com/cfagudelo96/zapatos-be/app/api/handlers"
	zapatohdlr "github.com/cfagudelo96/zapatos-be/app/api/handlers/zapato"
	"github.com/cfagudelo96/zapatos-be/business/zapato"
	"github.com/cfagudelo96/zapatos-be/business/zapato/store"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {
	e := echo.New()
	e.Validator = handlers.NewValidator()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	addZapatosRoutes(e)

	e.Logger.Fatal(e.Start(":1323"))
}

func addZapatosRoutes(e *echo.Echo) {
	zstore := store.NewInMemoryStore()
	zservice := zapato.NewService(zstore)
	zhandler := zapatohdlr.NewHandler(zservice)

	e.POST("/zapatos", zhandler.Create)
	e.GET("/zapatos", zhandler.List)
	e.GET("/zapatos/:id", zhandler.Get)
}
