package main

import (
	"TeleEcho/api/database"
	"TeleEcho/api/services"
	"TeleEcho/configs"
	"TeleEcho/router"
	"crypto/rand"
	"encoding/base64"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
	"net/http"
)

func main() {
	err := configs.ParseConfig()
	if err != nil {
		logrus.Printf("Can not read config")
	}
	err = database.ConnectDB()
	if err != nil {
		logrus.Printf("err:%s", err)
	}
	err = services.ConnectS3()
	if err != nil {
		logrus.Printf("err:%s", err)
	}
	e := router.New()
	err = e.Start(configs.Config.Address + ":" + configs.Config.Port)
	if err != nil {
		logrus.Printf("err:%s", err)
	}
}

var Store = sessions.NewCookieStore([]byte("your-secret-key"))

func CSRF(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		session, _ := Store.Get(c.Request(), "session-name")

		if c.Request().Method == http.MethodGet {
			token := generateCSRFToken()
			session.Values["csrf"] = token
			session.Save(c.Request(), c.Response().Writer)
			c.Response().Header().Set("X-CSRF-Token", token)
		} else {
			requestToken := c.Request().Header.Get("X-CSRF-Token")
			sessionToken := session.Values["csrf"]
			if requestToken == "" || sessionToken == nil || requestToken != sessionToken {
				return echo.NewHTTPError(http.StatusForbidden, "CSRF token mismatch")
			}
		}

		return next(c)
	}
}

func generateCSRFToken() string {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return ""
	}
	return base64.StdEncoding.EncodeToString(b)
}
