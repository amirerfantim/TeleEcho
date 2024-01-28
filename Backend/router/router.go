package router

import (
	"TeleEcho/api/handlers"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func New() *echo.Echo {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:5173"},
		AllowMethods: []string{echo.GET, echo.PUT, echo.POST, echo.DELETE, echo.PATCH},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))
	userGroup := e.Group("/user")
	userGroup.POST("/register", handlers.RegisterUser)
	userGroup.POST("/login", handlers.Login)
	return e
}
