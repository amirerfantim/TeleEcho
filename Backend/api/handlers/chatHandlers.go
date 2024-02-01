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
	receiverID := c.FormValue("receiverID")
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
func CreateGroupChat(c echo.Context) error {
	groupID := c.FormValue("groupID")
	groupIDInt, err := strconv.ParseUint(groupID, 10, 0)
	if err != nil {
		logrus.Printf("Error while parsing group id:%s\n", err)
		return c.JSON(http.StatusBadRequest, "Group id is wrong")
	}
	groupChatExist, err := database.DoesGroupChatExist(uint(groupIDInt))
	if err != nil {
		logrus.Printf("Error while checking group chat exist:%s\n", err)
		return c.JSON(http.StatusInternalServerError, "Can not check group chat")
	}
	if groupChatExist {
		return c.JSON(http.StatusBadRequest, fmt.Sprintf("This group chat for this group id %d is already exist", groupIDInt))
	}
	_, err = database.CreateGroupChat(uint(groupIDInt))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Error while creating chat")
	} else {
		return c.JSON(http.StatusCreated, "Group Chat created successfully.")
	}
}
func DeleteDirectChatHandler(c echo.Context) error {
	chatIDParam := c.Param("chatID")
	chatID, err := strconv.ParseUint(chatIDParam, 10, 0)
	if err != nil {
		logrus.WithError(err).WithField("chatID", chatIDParam).Error("Invalid chat ID")
		return c.JSON(http.StatusBadRequest, "Invalid chat ID")
	}

	if err := database.DeleteDirectChatAndMessages(uint(chatID)); err != nil {
		return c.JSON(http.StatusInternalServerError, "Failed to delete chat")
	}

	return c.NoContent(http.StatusOK)
}
func DeleteGroupChatHandler(c echo.Context) error {
	userID := c.Get("id").(string)
	userIDInt, err := strconv.ParseUint(userID, 10, 0)
	if err != nil {
		logrus.Printf("Error while parsing user id:%s\n", err)
		return c.JSON(http.StatusBadRequest, "User id is wrong")
	}
	groupIDParam := c.Param("groupID")
	groupIDInt, err := strconv.ParseUint(groupIDParam, 10, 0)
	if err != nil {
		logrus.WithError(err).WithField("groupID", groupIDParam).Error("Invalid group ID")
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid group ID")
	}

	groupExist, err := database.DoesGroupExistByID(uint(userIDInt), uint(groupIDInt))
	if err != nil {
		logrus.Printf("Error while checking groups:%s\n", err)
		return c.JSON(http.StatusInternalServerError, "Can not check your groups")
	}
	if !groupExist {
		return c.JSON(http.StatusBadRequest, "You can not delete chat for this group")
	}

	if err := database.DeleteGroupChatAndMessages(uint(groupIDInt)); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to delete group chat")
	}

	return c.NoContent(http.StatusOK)
}
