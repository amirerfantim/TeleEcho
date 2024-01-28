package middleware

import (
	"TeleEcho/configs"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"net/http"
)

func ValidateJWT(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		tokenString := c.Request().Header.Get("Authorization")
		token, _ := jwt.ParseWithClaims(tokenString, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(configs.Config.TokenKey), nil
		})
		if claims, ok := token.Claims.(*jwt.StandardClaims); ok && token.Valid {
			c.Set("id", claims.Subject)
			return next(c)
		} else {
			return c.String(http.StatusUnauthorized, "Invalid token")
		}
	}
}
