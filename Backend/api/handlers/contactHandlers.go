package handlers

import (
	"TeleEcho/api/database"
	"TeleEcho/model"
	"errors"
	"fmt"
	"github.com/labstack/echo"
	"net/http"
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
	contacts, err := database.GetUserContacts(uint(userIDInt))
	if err != nil {
		if errors.Is(err, database.NotFoundContact) {
			return c.JSON(http.StatusOK, []model.Contact{})
		}
		return c.JSON(http.StatusInternalServerError, "Error while querying contact")
	} else {
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
	newStatus := c.FormValue("status")
	var status model.Status
	if newStatus == "blocked" {
		status.IsBlocked = true
	}
	if newStatus == "hide profile" {
		status.ProfilePictureHide = true
	}
	if newStatus == "hide phone" {
		status.PhoneNumberHide = true
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
