package handlers

import (
	"TeleEcho/api/database"
	"fmt"
	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

func CreateDirectChat(c echo.Context) error {
	userID := c.Get("id").(string)
	userIDInt, err := strconv.ParseUint(userID, 10, 0)
	if err != nil {
		logrus.Printf("Error while parsing user id:%s\n", err)
		return c.JSON(http.StatusBadRequest, "User id is wrong")
	}
	receiverID := c.FormValue("id")
	receiverIDInt, err := strconv.ParseUint(receiverID, 10, 0)
	if err != nil {
		logrus.Printf("Error while parsing receiver id:%s\n", err)
		return c.JSON(http.StatusBadRequest, "receiver id is wrong")
	}
	chatExist, err := database.DoesDirectChatExist(uint(userIDInt), uint(receiverIDInt))
	if err != nil {
		logrus.Printf("Error while checking chat exist:%s\n", err)
		return c.JSON(http.StatusInternalServerError, "Can not check chat")
	}
	if chatExist {
		return c.JSON(http.StatusBadRequest, fmt.Sprintf("This chat between %d id and %d id is already exist", userIDInt, receiverIDInt))
	}

	_, err = database.CreateDirectChat(uint(userIDInt), uint(receiverIDInt))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Error while creating chat")
	} else {
		return c.JSON(http.StatusCreated, "Chat created successfully.")
	}
}
