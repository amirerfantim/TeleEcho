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
