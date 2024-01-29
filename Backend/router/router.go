package router

import (
	"TeleEcho/api/handlers"
	myMiddleware "TeleEcho/middleware"
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
	e.POST("/register", handlers.RegisterUser)
	e.POST("/login", handlers.Login)
	userGroup := e.Group("/users", myMiddleware.ValidateJWT)
	userGroup.DELETE("", handlers.DeleteUser)
	userGroup.GET("", handlers.GetUserInformation)
	userGroup.PATCH("", handlers.UpdateUserInformation)

	contactGroup := e.Group("/contacts", myMiddleware.ValidateJWT)
	contactGroup.POST("", handlers.CreateContact)
	contactGroup.GET("", handlers.GetUserContacts)
	contactGroup.DELETE("", handlers.DeleteContact)
	contactGroup.PATCH("", handlers.ChangeContentStatus)

	groupHandlersGroup := e.Group("/Group", myMiddleware.ValidateJWT)
	groupHandlersGroup.POST("", handlers.CreateGroup)
	groupHandlersGroup.GET("", handlers.GetUserGroups)

	return e
}
