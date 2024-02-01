package handlers

import (
	"TeleEcho/api/database"
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

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

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
func NewChatMessageWs(c echo.Context) error {
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		logrus.WithError(err).Error("Failed to upgrade to WebSocket")
		return err
	}
	defer ws.Close()

	senderID := c.Get("id").(string)
	senderIDInt, err := strconv.ParseUint(senderID, 10, 0)
	if err != nil {
		ws.WriteMessage(websocket.TextMessage, []byte("User id is wrong"))
		logrus.WithError(err).Error("Error while parsing user id")
		return err
	}

	chatID, err := strconv.ParseUint(c.QueryParam("chatID"), 10, 0)
	if err != nil {
		ws.WriteMessage(websocket.TextMessage, []byte("Chat id is wrong"))
		logrus.WithError(err).Error("Invalid chat ID")
		return err
	}

	for {
		var incomingMessage struct {
			Content string `json:"content"`
			Stat    string `json:"stat"`
		}
		if err := ws.ReadJSON(&incomingMessage); err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				logrus.WithError(err).Error("WebSocket unexpected close error")
			} else {
				logrus.WithError(err).Error("WebSocket read error")
			}
			return err
		}

		if incomingMessage.Stat == "exit" {
			break
		}

		chat, err := database.GetDirectChatByID(uint(chatID))
		if err != nil {
			if errors.Is(err, database.NotFoundChat) {
				ws.WriteMessage(websocket.TextMessage, []byte("This chat does not exist"))
				continue
			}
			ws.WriteMessage(websocket.TextMessage, []byte("Error while retrieving chat"))
			logrus.WithError(err).Error("Error while get chat")
			return err
		}

		if chat.SenderID != uint(senderIDInt) && chat.ReceiverID != uint(senderIDInt) {
			ws.WriteMessage(websocket.TextMessage, []byte("Cannot access this chat"))
			continue
		}

		if incomingMessage.Content == "" {
			ws.WriteMessage(websocket.TextMessage, []byte("Message content cannot be empty"))
			continue
		}

		_, err = database.CreateMessage(uint(chatID), model.TypeDirectChat, incomingMessage.Content)
		if err != nil {
			ws.WriteMessage(websocket.TextMessage, []byte("Can not create message"))
			logrus.WithError(err).Error("Error while creating message")
			return err
		}

		ws.WriteMessage(websocket.TextMessage, []byte("Message sent"))
	}

	return nil
}
func GetMessageByCountWs(c echo.Context) error {
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		logrus.WithError(err).Error("Failed to upgrade to WebSocket")
		return err
	}
	defer ws.Close()

	receiverID := c.Get("id").(string)
	receiverIDInt, err := strconv.ParseUint(receiverID, 10, 0)
	if err != nil {
		logrus.WithError(err).Error("Failed to validate JWT")
		ws.WriteMessage(websocket.TextMessage, []byte("Unauthorized access"))
		return echo.ErrUnauthorized
	}

	chatID, err := strconv.ParseUint(c.QueryParam("chatID"), 10, 0)
	if err != nil {
		logrus.WithError(err).Error("Invalid chat ID")
		ws.WriteMessage(websocket.TextMessage, []byte("Invalid chat ID"))
		return echo.ErrBadRequest
	}

	for {
		var requestData struct {
			Count uint   `json:"count"`
			Stat  string `json:"stat"`
		}

		if err := ws.ReadJSON(&requestData); err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				logrus.WithError(err).Error("WebSocket unexpected close error")
				return err
			}
			logrus.WithError(err).Error("WebSocket read error")
			break
		}

		if requestData.Stat == "exit" {
			break
		}

		chat, err := database.GetDirectChatByID(uint(chatID))
		if err != nil {
			logrus.WithError(err).Error("Error while retrieving chat")
			ws.WriteMessage(websocket.TextMessage, []byte("Error while retrieving chat"))
			return echo.ErrInternalServerError
		}

		if chat == nil || (chat.SenderID != uint(receiverIDInt) && chat.ReceiverID != uint(receiverIDInt)) {
			ws.WriteMessage(websocket.TextMessage, []byte("Cannot access this chat"))
			continue
		}

		messages, err := database.GetMessagesByChatID(uint(chatID), model.TypeDirectChat)
		if err != nil {
			logrus.WithError(err).Error("Error while retrieving messages")
			ws.WriteMessage(websocket.TextMessage, []byte("Error while retrieving messages"))
			return echo.ErrInternalServerError
		}

		sort.Slice(messages, func(i, j int) bool {
			return messages[i].CreatedAt.After(messages[j].CreatedAt)
		})
		count := requestData.Count
		if count > uint(len(messages)) {
			count = uint(len(messages))
		}

		if err := ws.WriteJSON(messages[:count]); err != nil {
			logrus.WithError(err).Error("Failed to write messages to WebSocket")
			return err
		}
	}

	return nil
}
func GetDirectChatsWs(c echo.Context) error {
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		logrus.WithError(err).Error("Failed to upgrade to WebSocket")
		return err
	}
	defer ws.Close()

	userID := c.Get("id").(string)
	userIDInt, err := strconv.ParseUint(userID, 10, 0)
	if err != nil {
		logrus.WithError(err).Error("Failed to validate JWT")
		ws.WriteMessage(websocket.TextMessage, []byte("Unauthorized access"))
		return echo.ErrUnauthorized
	}
	sendUpdatedChats := func() error {

		receiverChats, err := database.GetChatsByReceiverID(uint(userIDInt))
		if err != nil {
			logrus.WithError(err).Error("Failed to retrieve receiver chats")
			return echo.ErrInternalServerError
		}
		senderChats, err := database.GetChatsBySenderID(uint(userIDInt))
		if err != nil {
			logrus.WithError(err).Error("Failed to retrieve sender chats")
			return echo.ErrInternalServerError
		}

		allChats := append(senderChats, receiverChats...)

		senderUsers, err := database.GetUsersForChatsByReceiverID(uint(userIDInt))

		if err != nil {
			logrus.WithError(err).Error("Failed to retrieve sender chat users")
			return echo.ErrInternalServerError
		}
		type response struct {
			User          model.User `json:"user"`
			UnreadMessage int        `json:"unreadMessage"`
		}

		if len(allChats) == 0 {
			return nil
		}

		res := make([]response, len(receiverChats))
		countNumber := 0
		for _, chat := range receiverChats {
			count, err := database.CountUnreadMessages(chat.ID, model.TypeDirectChat)
			if err != nil {
				logrus.WithError(err).Error("Failed to count unread messages")
				return echo.ErrInternalServerError
			}
			for _, user := range senderUsers {
				if user.ID == chat.SenderID {
					res[countNumber] = response{
						User:          user,
						UnreadMessage: int(count),
					}
					countNumber++
					break
				}
			}
		}

		return ws.WriteJSON(res)
	}

	for {
		err = sendUpdatedChats()
		if err != nil {
			return err
		}

		type msg struct {
			Message string `json:"message"`
		}
		var m msg
		err = ws.ReadJSON(&m)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				logrus.WithError(err).Error("WebSocket unexpected close error")
				return err
			}
			break
		}

		if m.Message == "new" {
			err = sendUpdatedChats()
			if err != nil {
				return err
			}
		}
		if m.Message == "exit" {
			break
		}
	}

	return nil
}
func GetDirectChatMessagesHandler(c echo.Context) error {
	chatIDParam := c.QueryParam("chatID")
	chatID, err := strconv.ParseUint(chatIDParam, 10, 0)
	if err != nil {
		logrus.WithError(err).WithField("chatID", chatIDParam).Error("Invalid chat ID")
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid chat ID")
	}

	messages, err := database.GetMessagesByChatID(uint(chatID), model.TypeDirectChat)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to retrieve messages")
	}

	return c.JSON(http.StatusOK, messages)
}
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
