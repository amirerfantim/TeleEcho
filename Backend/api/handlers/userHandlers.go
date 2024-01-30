package handlers

import (
	"TeleEcho/api/database"
	"TeleEcho/api/services"
	"TeleEcho/configs"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"io/ioutil"
	"net/http"
	"os"
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
				return c.String(http.StatusInternalServerError, "Unable to upload profile picture")
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
func DeleteUser(c echo.Context) error {
	userID := c.Get("id").(string)
	userIDInt, err := strconv.ParseUint(userID, 10, 0)
	if err != nil {
		fmt.Printf("Error while parsing user id:%s\n", err)
		return c.JSON(http.StatusBadRequest, "User id is wrong")
	}
	err = database.DeleteUserByUserID(uint(userIDInt))

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, map[string]interface{}{
				"error": fmt.Sprintf("User with ID %d not found", userIDInt),
			})
		} else {
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"error": "Failed to delete user.",
			})
		}
	}
	return c.NoContent(http.StatusNoContent)
}
func GetUserInformation(c echo.Context) error {
	username := c.FormValue("username")
	if username == "" {
		userID := c.Get("id").(string)
		userIDInt, err := strconv.ParseUint(userID, 10, 0)
		if err != nil {
			fmt.Printf("Error while parsing user id:%s\n", err)
			return c.JSON(http.StatusBadRequest, "User id is wrong")
		}
		user, err := database.GetUserByUserID(uint(userIDInt))
		if err != nil {
			if errors.Is(err, database.NotFoundUser) {
				return c.JSON(http.StatusNotFound, map[string]interface{}{
					"error": fmt.Sprintf("User with ID %d not found", userIDInt),
				})
			} else {
				return c.JSON(http.StatusInternalServerError, map[string]interface{}{
					"error": "Failed to retrieve user",
				})
			}
		} else {
			profilePhotoFile, err := services.DownloadS3(services.StorageSession, configs.Config.StorageServiceBucket, user.ProfilePicture)
			if err != nil {
				logrus.Println("Can not download photo:", err)
				return c.JSON(http.StatusInternalServerError, "Error while downloading photo.")
			}
			file, err := os.Open(profilePhotoFile.Name())
			if err != nil {
				logrus.Println("Can not open file", err)
				return c.JSON(http.StatusInternalServerError, "Error while processing photo.")
			}
			bytes, err := ioutil.ReadAll(file)
			if err != nil {
				logrus.Println("Can not convert photo to bytes.")
				return c.JSON(http.StatusInternalServerError, "Error while processing photo.")
			}
			base64ProfilePicture := base64.StdEncoding.EncodeToString(bytes)
			user.ProfilePicture = base64ProfilePicture

			return c.JSON(http.StatusOK, user)

		}

	} else {
		searchedUser, err := database.GetUserByUsername(username)
		if err != nil {
			if errors.Is(err, database.NotFoundUser) {
				return c.JSON(http.StatusBadRequest, fmt.Sprintf("No user found with username %s", username))
			}
			return c.JSON(http.StatusInternalServerError, "Error while finding user")
		}
		profilePhotoFile, err := services.DownloadS3(services.StorageSession, configs.Config.StorageServiceBucket, searchedUser.ProfilePicture)
		if err != nil {
			logrus.Println("Can not download photo:", err)
			return c.JSON(http.StatusInternalServerError, "Error while downloading photo.")
		}
		file, err := os.Open(profilePhotoFile.Name())
		if err != nil {
			logrus.Println("Can not open file", err)
			return c.JSON(http.StatusInternalServerError, "Error while processing photo.")
		}
		bytes, err := ioutil.ReadAll(file)
		if err != nil {
			logrus.Println("Can not convert photo to bytes.")
			return c.JSON(http.StatusInternalServerError, "Error while processing photo.")
		}
		base64ProfilePicture := base64.StdEncoding.EncodeToString(bytes)
		searchedUser.ProfilePicture = base64ProfilePicture
		return c.JSON(http.StatusOK, searchedUser)
	}

}
func UpdateUserInformation(c echo.Context) error {
	userID := c.Get("id").(string)
	userIDInt, err := strconv.ParseUint(userID, 10, 0)
	if err != nil {
		fmt.Printf("Error while parsing user id:%s\n", err)
		return c.JSON(http.StatusBadRequest, "User id is wrong")
	}
	user, err := database.GetUserByUserID(uint(userIDInt))
	if err != nil {
		if errors.Is(err, database.NotFoundUser) {
			return c.JSON(http.StatusNotFound, map[string]interface{}{
				"error": fmt.Sprintf("User with ID %d not found", userIDInt),
			})
		} else {
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"error": "Failed to retrieve user",
			})
		}
	}
	username := c.FormValue("username")
	if username != "" {
		usernameOK := database.IsUsernameDuplicate(username)
		if usernameOK {
			user.Username = username
		} else {
			return c.JSON(http.StatusBadRequest, "Your username has already been used")
		}
	}
	firstname := c.FormValue("firstname")
	if firstname != "" {
		user.Firstname = firstname
	}
	logrus.Printf(firstname)
	lastname := c.FormValue("lastname")
	if lastname != "" {
		user.Lastname = lastname
	}
	phone := c.FormValue("phone")
	if phone != "" {
		phoneNumberOK := database.IsPhoneDuplicate(phone)
		if phoneNumberOK {
			user.Phone = phone
		} else {
			return c.JSON(http.StatusBadRequest, "Your phone number has already been used")
		}
	}
	password := c.FormValue("password")
	if password != "" {
		if len(password) < 8 {
			return c.JSON(http.StatusBadRequest, "Password is too easy")
		} else {
			hashFunc := sha256.New()
			hashFunc.Write([]byte(password))
			hashPassword := hex.EncodeToString(hashFunc.Sum(nil))
			user.Password = hashPassword
		}
	}
	logrus.Infof(username)
	profilePicture, err := c.FormFile("profile")
	println(profilePicture)
	if err != nil {
		if err != http.ErrMissingFile {
			logrus.Printf("Error opening file: %v\n", err)
			return c.JSON(http.StatusBadRequest, "Unable to open file")
		}
	} else {
		profilePath, err := services.UploadS3(services.StorageSession, profilePicture, configs.Config.StorageServiceBucket, username)
		if err != nil {
			logrus.Printf("Unable to upload image\n")
			return c.String(http.StatusInternalServerError, "Unable to upload profile picture")
		}
		user.ProfilePicture = profilePath
	}

	bio := c.FormValue("bio")
	if bio != "" {
		user.Bio = bio
	}
	err = database.UpdateUserByUserID(uint(userIDInt), *user)
	if err != nil {
		logrus.Printf("Error while updating user:%s\n", err)
		if errors.Is(err, database.NotFoundUser) {
			return c.JSON(http.StatusNotFound, map[string]interface{}{
				"error": fmt.Sprintf("User with ID %d not found", userIDInt),
			})
		} else {
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"error": "Failed to update user",
			})
		}
	}
	return c.JSON(http.StatusOK, "User updated successfully")

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
