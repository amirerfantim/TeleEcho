package router

import (
	"TeleEcho/api/handlers"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func New() *echo.Echo {
	e := echo.New()
	e.Use(middleware.Recover())
	userGroup := e.Group("/user")
	userGroup.POST("/register", handlers.RegisterUser)
	userGroup.POST("/login", handlers.Login)
	return e
}
