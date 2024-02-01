package handlers

import (
	"TeleEcho/api/database"
	"TeleEcho/model"
	"fmt"
	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
	"net/http"
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

//func GetGroupDataWs(c echo.Context) error {
//	groupID, err := strconv.ParseUint(c.Param("groupid"), 10, 64)
//	if err != nil {
//		return echo.ErrBadRequest
//	}
//
//	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
//	if err != nil {
//		return err
//	}
//	defer ws.Close()
//
//	for {
//		var requestData struct {
//			Stat string `json:"stat"`
//		}
//
//		err = ws.ReadJSON(&requestData)
//		if err != nil {
//			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
//				return err
//			}
//			break
//		}
//
//		if requestData.Stat == "exit" {
//			break
//		}
//
//		groupID := groupID
//
//		group, users, err := g.userGroupRepo.GetGroupWithUserGroups(c.Request().Context(), groupID)
//		if err != nil {
//			return echo.ErrInternalServerError
//		}
//
//		if len(group) == 0 {
//			ws.WriteMessage(websocket.TextMessage, []byte("No groups found"))
//			continue
//		}
//
//		err = ws.WriteJSON(echo.Map{
//			"group": group,
//			"users": users,
//		})
//		if err != nil {
//			return err
//		}
//	}
//
//	return nil
//}
//func GetGroupChatWs(c echo.Context) error {
//	groupID, err := strconv.ParseUint(c.Param("group_id"), 10, 0)
//	if err != nil {
//		logrus.WithError(err).Error("Failed to parse group ID")
//		return echo.ErrBadRequest
//	}
//
//	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
//	if err != nil {
//		logrus.WithError(err).Error("Failed to upgrade WebSocket connection")
//		return err
//	}
//	defer ws.Close()
//
//	for {
//		var requestData struct {
//			Stat string `json:"stat"`
//		}
//
//		err = ws.ReadJSON(&requestData)
//		if err != nil {
//			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
//				logrus.WithError(err).Error("Unexpected WebSocket close error")
//				return err
//			}
//			break
//		}
//
//		if requestData.Stat == "exit" {
//			break
//		}
//		users, err := database.GetAllUsersInGroup(uint(groupID))
//		group, users, err := g.userGroupRepo.GetGroupWithUserGroups(c.Request().Context(), groupID)
//		if err != nil {
//			logrus.WithError(err).Error("Failed to fetch group data from the repository")
//			return echo.ErrInternalServerError
//		}
//
//		if len(group) == 0 {
//			err := ws.WriteMessage(websocket.TextMessage, []byte("No groups found"))
//			if err != nil {
//				logrus.WithError(err).Error("Failed to send 'No groups found' message over WebSocket")
//				return err
//			}
//			continue
//		}
//
//		err = ws.WriteJSON(echo.Map{
//			"group": group,
//			"users": users,
//		})
//		if err != nil {
//			logrus.WithError(err).Error("Failed to send group data over WebSocket")
//			return err
//		}
//	}
//
//	return nil
//}
