package handlers

import (
	"TeleEcho/api/database"
	"TeleEcho/api/services"
	"TeleEcho/configs"
	"TeleEcho/model"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
)

func CreateContact(c echo.Context) error {
	userID := c.Get("id").(string)
	userIDInt, err := strconv.ParseUint(userID, 10, 0)
	if err != nil {
		fmt.Printf("Error while parsing user id:%s\n", err)
		return c.JSON(http.StatusBadRequest, "User id is wrong")
	}
	phone := c.FormValue("phone")
	contactUser, err := database.GetUserByPhone(phone)
	if err != nil {
		if errors.Is(err, database.NotFoundUser) {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"error": fmt.Sprintf("User with phone %s not found", phone),
			})
		}
		return c.JSON(http.StatusInternalServerError, "Error while creating contact")
	}
	if uint(userIDInt) == contactUser.ID {
		return c.JSON(http.StatusBadRequest, "Phone number is yours")
	}
	doesUserHaveContact, err := database.DoesUserHaveContact(uint(userIDInt), contactUser.ID)
	if err != nil {
		logrus.Printf("Error while checking contact for user")
		return c.JSON(http.StatusInternalServerError, "Can not check contact")

	}
	if doesUserHaveContact {
		return c.JSON(http.StatusBadRequest, fmt.Sprintf("You have already this contact with with user id %d", userIDInt))
	}
	status := model.Status{ProfilePictureHide: false, PhoneNumberHide: false, IsBlocked: false}
	err = database.CreateContact(uint(userIDInt), contactUser.ID, status)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Error while creating contact")
	} else {
		return c.JSON(http.StatusCreated, "Contact created successfully")
	}
}
func GetUserContacts(c echo.Context) error {
	userID := c.Get("id").(string)
	userIDInt, err := strconv.ParseUint(userID, 10, 0)
	if err != nil {
		fmt.Printf("Error while parsing user id:%s\n", err)
		return c.JSON(http.StatusBadRequest, "User id is wrong")
	}
	contactUserID := c.QueryParam("id")

	contacts, err := database.GetUserContactsInfo(uint(userIDInt))
	if err != nil {
		if errors.Is(err, database.NotFoundContact) {
			return c.JSON(http.StatusOK, []model.Contact{})
		}
		return c.JSON(http.StatusInternalServerError, "Error while querying contact")
	} else {
		if contactUserID != "" {
			contactUserIDInt, err := strconv.ParseUint(contactUserID, 10, 0)
			if err != nil {
				fmt.Printf("Error while parsing contact user id:%s\n", err)
				return c.JSON(http.StatusBadRequest, "Contact user id is wrong")
			}
			for _, contact := range contacts {
				if contact.ID == uint(contactUserIDInt) {
					profilePhotoFile, err := services.DownloadS3(services.StorageSession, configs.Config.StorageServiceBucket, contact.ProfilePicture)
					if err != nil {
						logrus.Println("Can not download photo:", err)
						return c.JSON(http.StatusInternalServerError, "Error while downloading photo.")
					}
					file, err := os.Open(profilePhotoFile.Name())
					if err != nil {
						logrus.Println("Can not open image file", err)
						return c.JSON(http.StatusInternalServerError, "Error while processing photo.")
					}
					bytes, err := ioutil.ReadAll(file)
					if err != nil {
						logrus.Println("Can not convert photo to bytes.")
						return c.JSON(http.StatusInternalServerError, "Error while processing photo.")
					}
					base64ProfilePicture := base64.StdEncoding.EncodeToString(bytes)
					contact.ProfilePicture = base64ProfilePicture
					return c.JSON(http.StatusOK, contact)
				}
			}
		}
		return c.JSON(http.StatusOK, contacts)
	}
}
func DeleteContact(c echo.Context) error {
	userID := c.Get("id").(string)
	userIDInt, err := strconv.ParseUint(userID, 10, 0)
	if err != nil {
		fmt.Printf("Error while parsing user id:%s\n", err)
		return c.JSON(http.StatusBadRequest, "User id is wrong")
	}
	username := c.FormValue("username")
	contact, err := database.GetUserByUsername(username)
	if err != nil {
		if errors.Is(err, database.NotFoundUser) {
			return c.JSON(http.StatusBadRequest, fmt.Sprintf("No user found with username %s", username))
		}
		return c.JSON(http.StatusInternalServerError, "Error while finding user")
	}
	err = database.DeleteUserContact(uint(userIDInt), contact.ID)
	if err != nil {
		if errors.Is(err, database.NotFoundContact) {
			return c.JSON(http.StatusBadRequest, fmt.Sprintf("No contact found with ContactUserID %d for user ID %d", contact.ID, userIDInt))
		}
		return c.JSON(http.StatusInternalServerError, "Error while deleting contact")
	} else {
		return c.JSON(http.StatusNoContent, "")
	}
}
func ChangeContentStatus(c echo.Context) error {
	userID := c.Get("id").(string)
	userIDInt, err := strconv.ParseUint(userID, 10, 0)
	if err != nil {
		fmt.Printf("Error while parsing user id:%s\n", err)
		return c.JSON(http.StatusBadRequest, "User id is wrong")
	}
	username := c.FormValue("username")
	contact, err := database.GetUserByUsername(username)
	if err != nil {
		if errors.Is(err, database.NotFoundUser) {
			return c.JSON(http.StatusBadRequest, fmt.Sprintf("No user found with username %s", username))
		}
		return c.JSON(http.StatusInternalServerError, "Error while finding user")
	}
	newStatusJSON := c.FormValue("status")
	var status model.Status

	err = json.Unmarshal([]byte(newStatusJSON), &status)
	if err != nil {
		logrus.Printf("Error parsing status JSON: %s", err)
		return c.JSON(http.StatusBadRequest, "Can not parse data.")
	}

	err = database.UpdateContactStatus(uint(userIDInt), contact.ID, status)
	if err != nil {
		if errors.Is(err, database.NotFoundContact) {
			return c.JSON(http.StatusBadRequest, fmt.Sprintf("No contact found with username %s for user ID %d", contact.ID, userIDInt))
		}
		return c.JSON(http.StatusInternalServerError, "Error while deleting contact")
	} else {
		return c.JSON(http.StatusNoContent, "")
	}
}
