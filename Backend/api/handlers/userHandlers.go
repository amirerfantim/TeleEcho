package handlers

import (
	"TeleEcho/api/database"
	"TeleEcho/api/services"
	"TeleEcho/configs"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"time"
)

func RegisterUser(c echo.Context) error {
	username := c.FormValue("username")
	firstname := c.FormValue("firstname")
	lastname := c.FormValue("lastname")
	phone := c.FormValue("phone")
	password := c.FormValue("password")
	profilePicture, err := c.FormFile("profile")
	if err != nil {
		logrus.Printf("Unable to open image\n")
		return c.String(http.StatusBadRequest, "Unable to open file")
	}
	bio := c.FormValue("bio")
	usernameOK := database.IsUsernameDuplicate(username)
	phoneNumberOk := database.IsPhoneDuplicate(phone)
	if usernameOK {
		if phoneNumberOk {
			profilePath, err := services.UploadS3(services.StorageSession, profilePicture, configs.Config.StorageServiceBucket, username)
			if err != nil {
				logrus.Printf("Unable to upload image\n")
				return c.String(http.StatusBadRequest, "Unable to upload first profile picture")
			}
			hashFunc := sha256.New()
			hashFunc.Write([]byte(password))
			hashPassword := hex.EncodeToString(hashFunc.Sum(nil))

			err = database.CreateUser(username, firstname, lastname, phone, hashPassword, profilePath, bio)
			if err != nil {
				logrus.Printf("Error creating user:%s\n", err)
				return c.JSON(http.StatusInternalServerError, "Can not create user")
			}
			return c.JSON(http.StatusCreated, "User created successfully")
		} else {
			return c.JSON(http.StatusBadRequest, "Your phone number has already been used")
		}
	} else {
		return c.JSON(http.StatusBadRequest, "Your username has already been used")
	}

}
func Login(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")
	hashFunc := sha256.New()
	hashFunc.Write([]byte(password))
	hashPassword := hex.EncodeToString(hashFunc.Sum(nil))
	user, err := database.CheckPassword(username, hashPassword)
	if err != nil {
		logrus.Printf("Error while checking username and password:%s\n", err)
		if errors.Is(err, database.NotFoundUser) {
			return c.JSON(http.StatusBadRequest, "Username is wrong.")
		} else if errors.Is(err, database.IncorrectPassword) {
			return c.JSON(http.StatusBadRequest, "Password is wrong.")
		} else {
			return c.String(http.StatusInternalServerError, "Can not check username and password")
		}

	} else {
		token, err := generateJWT(user.ID)
		if err != nil {
			logrus.Printf("Error while generating token:%s", err)
			return c.JSON(http.StatusInternalServerError, "Can not create token")
		}
		return c.JSON(http.StatusOK, map[string]string{"token": token})
	}

}
func generateJWT(userID uint) (string, error) {
	claims := &jwt.StandardClaims{
		Subject:   strconv.Itoa(int(userID)),
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(configs.Config.TokenKey))
	if err != nil {
		logrus.Printf("Error sign token:%s", err)
		return "", err
	}
	return signedToken, nil
}
