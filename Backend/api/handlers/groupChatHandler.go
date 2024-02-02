package handlers

import (
	"TeleEcho/api/database"
	"TeleEcho/middleware"
	"TeleEcho/model"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
	"net/http"
	"sort"
	"strconv"
)

func GetGroupChatMessagesHandler(c echo.Context) error {
	chatIDParam := c.QueryParam("chatID")
	chatID, err := strconv.ParseUint(chatIDParam, 10, 0)
	if err != nil {
		logrus.WithError(err).WithField("chatID", chatIDParam).Error("Invalid chat ID")
		return c.JSON(http.StatusBadRequest, "Invalid chat ID")
	}

	messages, err := database.GetMessagesByChatID(uint(chatID), model.TypeGroupChat)
	if err != nil {
		logrus.Printf("Can not rerive messages: %e\n", err)
		return c.JSON(http.StatusInternalServerError, "Failed to retrieve messages")
	}

	return c.JSON(http.StatusOK, messages)
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

func GetGroupChatWs(c echo.Context) error {
	groupID, err := strconv.ParseUint(c.QueryParam("groupID"), 10, 0)
	if err != nil {
		logrus.WithError(err).Error("Failed to parse group ID")
		return echo.ErrBadRequest
	}

	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		logrus.WithError(err).Error("Failed to upgrade WebSocket connection")
		return err
	}
	defer ws.Close()

	for {
		var requestData struct {
			Stat string `json:"stat"`
		}

		err = ws.ReadJSON(&requestData)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				logrus.WithError(err).Error("Unexpected WebSocket close error")
				return err
			}
			return echo.ErrInternalServerError
		}

		if requestData.Stat == "exit" {
			break
		}
		group, err := database.GetGroupByID(uint(groupID))
		if err != nil {
			if errors.Is(err, database.NotFoundGroup) {
				err := ws.WriteMessage(websocket.TextMessage, []byte("No groups found"))
				if err != nil {
					logrus.WithError(err).Error("Failed to send 'No groups found' message over WebSocket")
					return err
				}
				continue
			}
			return echo.ErrInternalServerError
		}
		users, err := database.GetAllUsersInGroup(uint(groupID))
		if err != nil {
			return echo.ErrInternalServerError
		}

		err = ws.WriteJSON(echo.Map{
			"group": group,
			"users": users,
		})
		if err != nil {
			logrus.WithError(err).Error("Failed to send group data over WebSocket")
			return err
		}
	}
	return nil
}

func NewGroupMessageWs(c echo.Context) error {
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		logrus.WithError(err).Error("Failed to upgrade WebSocket connection")
		return err
	}
	defer ws.Close()

	userIDInt, checkJWT := middleware.ValidateJWTToken(c.QueryParam("token"))
	if !checkJWT {
		logrus.WithError(err).Error("Failed to validate JWT")
		ws.WriteMessage(websocket.TextMessage, []byte("Unauthorized access"))
		return echo.ErrUnauthorized
	}

	groupID, err := strconv.ParseUint(c.QueryParam("groupID"), 10, 64)
	if err != nil {
		logrus.WithError(err).Error("Failed to parse group ID")
		return echo.ErrBadRequest
	}

	for {
		var incomingMessage struct {
			UserID         uint   `json:"userid"`
			MessageContent string `json:"content"`
			Stat           string `json:"stat"`
		}
		err = ws.ReadJSON(&incomingMessage)
		logrus.Println(incomingMessage)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				logrus.WithError(err).Error("Unexpected WebSocket close error")
				return err
			}
			return echo.ErrInternalServerError
		}

		if incomingMessage.Stat == "exit" {
			break
		}
		userID := incomingMessage.UserID
		messageContent := incomingMessage.MessageContent

		if userID != uint(userIDInt) {
			logrus.WithField("UserID", userID).Warn("Unauthorized user")
			ws.WriteMessage(websocket.TextMessage, []byte("Unauthorized"))
			continue
		}
		check, err := database.IsUserInGroup(uint(userIDInt), uint(groupID))
		if err != nil {
			logrus.WithError(err).Error("Failed to fetch user group data")
			return echo.ErrInternalServerError
		}

		if !check {
			logrus.Warn("User is not part of this group")
			ws.WriteMessage(websocket.TextMessage, []byte("You are not part of this group"))
			continue
		}

		if messageContent == "" {
			logrus.Warn("Message body is empty")
			ws.WriteMessage(websocket.TextMessage, []byte("Message body cannot be empty"))
			continue
		}

		_, err = database.CreateMessage(uint(groupID), model.TypeGroupChat, messageContent)

		if err != nil {
			logrus.WithError(err).Error("Failed to create a new message")
			return echo.ErrInternalServerError
		}

		ws.WriteMessage(websocket.TextMessage, []byte("Message sent"))
	}

	return nil
}
func GetGroupMessagesWs(c echo.Context) error {
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		logrus.WithError(err).Error("WebSocket upgrade failed")
		return err
	}
	defer ws.Close()
	userIDInt, checkJWT := middleware.ValidateJWTToken(c.QueryParam("token"))
	if !checkJWT {
		logrus.WithError(err).Error("Failed to validate JWT")
		ws.WriteMessage(websocket.TextMessage, []byte("Unauthorized access"))
		return echo.ErrUnauthorized
	}

	groupIDParam := c.QueryParam("groupID")
	groupIDInt, err := strconv.ParseUint(groupIDParam, 10, 0)
	if err != nil {
		logrus.WithError(err).Error("Invalid group ID")
		return echo.ErrBadRequest
	}
	for {
		var requestData struct {
			Count uint64 `json:"count"`
			Stat  string `json:"stat"`
		}

		if err := ws.ReadJSON(&requestData); err != nil {
			logrus.WithError(err).Error("Failed to read message from WebSocket")
			break
		}

		if requestData.Stat == "exit" {
			break
		}
		isMember, err := database.IsUserInGroup(uint(userIDInt), uint(groupIDInt))
		if err != nil {
			logrus.WithError(err).Error("Can not check user in group")
			return echo.ErrInternalServerError
		}
		if !isMember {
			ws.WriteMessage(websocket.TextMessage, []byte("You are not part of this group"))
			continue
		}
		groupChat, err := database.FindGroupChatByGroupID(uint(groupIDInt))
		if err != nil {
			logrus.WithError(err).Error("Failed to retrieve group chat")
			return echo.ErrInternalServerError
		}
		messages, err := database.GetMessagesByChatIDAndType(groupChat.ID, model.TypeGroupChat)
		if err != nil {
			logrus.WithError(err).Error("Failed to retrieve group messages")
			return echo.ErrInternalServerError
		}
		sort.Slice(messages, func(i, j int) bool {
			return messages[i].CreatedAt.After(messages[j].CreatedAt)
		})
		count := requestData.Count
		if count > uint64(len(messages)) {
			count = uint64(len(messages))
		}
		if err := ws.WriteJSON(messages[:count]); err != nil {
			logrus.WithError(err).Error("Failed to write messages to WebSocket")
			return err
		}
	}

	return nil
}
