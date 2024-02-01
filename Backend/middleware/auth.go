package middleware

import (
	"TeleEcho/configs"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

func ValidateJWT(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		tokenString := c.Request().Header.Get("Authorization")

		if tokenString == "" {
			return c.String(http.StatusUnauthorized, "Missing token")
		}
		token, err := jwt.ParseWithClaims(tokenString, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(configs.Config.TokenKey), nil
		})

		if err != nil {
			return c.String(http.StatusUnauthorized, "Invalid token: "+err.Error())
		}
		if claims, ok := token.Claims.(*jwt.StandardClaims); ok && token.Valid {
			c.Set("id", claims.Subject)
			return next(c)
		} else {
			return c.String(http.StatusUnauthorized, "Invalid token")
		}
	}
}

func ValidateJWTToken(tokenString string) (uint64, bool) {
	if tokenString == "" {
		logrus.Error("Missing token")
		return 0, false
	}

	token, err := jwt.ParseWithClaims(tokenString, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			logrus.Errorf("Unexpected signing method: %v", token.Header["alg"])
			return nil, jwt.NewValidationError("invalid signing method", jwt.ValidationErrorSignatureInvalid)
		}
		return []byte(configs.Config.TokenKey), nil
	})

	if err != nil {
		logrus.WithError(err).Error("Invalid token")
		return 0, false
	}

	if claims, ok := token.Claims.(*jwt.StandardClaims); ok && token.Valid {
		userID, err := strconv.ParseUint(claims.Subject, 10, 0)
		if err != nil {
			logrus.WithError(err).Error("Failed to parse user ID from token")
			return 0, false
		}
		return userID, true
	} else {
		logrus.Error("Invalid token")
		return 0, false
	}
}
