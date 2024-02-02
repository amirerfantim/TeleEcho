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
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
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

	groupHandlersGroup := e.Group("/group", myMiddleware.ValidateJWT)
	groupHandlersGroup.POST("", handlers.CreateGroup)
	groupHandlersGroup.GET("", handlers.GetUserGroups)
	groupHandlersGroup.PATCH("", handlers.AddUserToGroup)
	groupHandlersGroup.DELETE("", handlers.RemoveUserGroup)
	groupHandlersGroup.GET("/all", handlers.GetAllUsersInGroup)

	groupChatHandlersGroup := e.Group("/group-chat", myMiddleware.ValidateJWT)
	groupChatHandlersGroup.POST("", handlers.CreateGroupChat)
	groupChatHandlersGroup.DELETE("", handlers.DeleteGroupChatHandler)
	groupChatHandlersGroup.GET("", handlers.GetGroupChatMessagesHandler)

	directChatHandlersGroup := e.Group("/chat", myMiddleware.ValidateJWT)
	directChatHandlersGroup.POST("", handlers.CreateDirectChat)
	directChatHandlersGroup.DELETE("/:chatID", handlers.DeleteDirectChatHandler)
	directChatHandlersGroup.GET("/message", handlers.GetDirectChatMessagesHandler)
	directChatHandlersGroup.DELETE("/message", handlers.DeleteMessageHandler)

	e.GET("/all-chat", handlers.GetDirectChatsWs)
	e.GET("/write-message", handlers.NewChatMessageWs)
	e.GET("/read-message", handlers.GetMessageByCountWs)

	e.GET("/all-group-chat", handlers.GetGroupChatWs)

	return e
}
