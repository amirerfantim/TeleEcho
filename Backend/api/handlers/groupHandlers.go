package handlers

import (
	"TeleEcho/api/database"
	"TeleEcho/api/services"
	"TeleEcho/configs"
	"TeleEcho/model"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
)

func CreateGroup(c echo.Context) error {
	userID := c.Get("id").(string)
	userIDInt, err := strconv.ParseUint(userID, 10, 0)
	if err != nil {
		logrus.Printf("Error while parsing user id:%s\n", err)
		return c.JSON(http.StatusBadRequest, "User id is wrong")
	}
	groupName := c.FormValue("name")
	groupExist, err := database.DoesGroupExist(uint(userIDInt), groupName)
	if err != nil {
		logrus.Printf("Error while checking groups:%s\n", err)
		return c.JSON(http.StatusInternalServerError, "Can not check your groups")
	}
	if groupExist {
		return c.JSON(http.StatusBadRequest, "This group with your user id is already exist")
	}
	groupDescription := c.FormValue("description")
	profilePicture, err := c.FormFile("profile")
	if err != nil {
		logrus.Printf("Unable to open image\n")
		return c.String(http.StatusBadRequest, "Unable to open file")
	}
	profilePath, err := services.UploadS3(services.StorageSession, profilePicture, configs.Config.StorageServiceBucket, userID+""+groupName)
	if err != nil {
		logrus.Printf("Unable to upload image\n")
		return c.String(http.StatusInternalServerError, "Unable to upload profile picture")
	}

	err = database.CreateGroup(uint(userIDInt), groupName, groupDescription, profilePath)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Error while creating group")
	} else {
		return c.JSON(http.StatusCreated, "Group created successfully.")
	}
}
func AddUserToGroup(c echo.Context) error {
	userID := c.Get("id").(string)
	userIDInt, err := strconv.ParseUint(userID, 10, 0)
	if err != nil {
		logrus.Printf("Error while parsing user id:%s\n", err)
		return c.JSON(http.StatusBadRequest, "User id is wrong")
	}
	groupID := c.FormValue("groupID")
	username := c.FormValue("username")
	groupIDInt, err := strconv.ParseUint(groupID, 10, 0)
	if err != nil {
		logrus.Printf("Error while parsing group id:%s\n", err)
		return c.JSON(http.StatusBadRequest, "Group id is wrong")
	}
	isUserGroup, err := database.IsUserInGroup(uint(userIDInt), uint(groupIDInt))
	if !isUserGroup {
		return c.JSON(http.StatusBadRequest, fmt.Sprintf("You don't have group with group id %d", groupIDInt))
	}
	searchedUser, err := database.GetUserByUsername(username)
	if err != nil {
		if errors.Is(err, database.NotFoundUser) {
			return c.JSON(http.StatusBadRequest, fmt.Sprintf("No user found with username %s", username))
		}
		return c.JSON(http.StatusInternalServerError, "Error while finding user")
	}
	if searchedUser.ID == uint(userIDInt) {
		return c.JSON(http.StatusBadRequest, fmt.Sprintf("This user id is yours: %d", userIDInt))
	}
	isUserGroup, err = database.IsUserInGroup(searchedUser.ID, uint(groupIDInt))
	if err != nil {
		logrus.Printf("Error while checking user group :%e", err)
		return c.JSON(http.StatusInternalServerError, "Can not check user group")
	}
	if isUserGroup {
		return c.JSON(http.StatusBadRequest, fmt.Sprintf("This user already have group with group id %d", groupIDInt))
	}
	err = database.AddUserToGroup(searchedUser.ID, uint(groupIDInt))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": "Failed to add user to group",
		})
	}
	return c.JSON(http.StatusCreated, "User added to the group successfully.")
}
func GetUserGroups(c echo.Context) error {
	groupName := c.QueryParam("name")
	userID := c.Get("id").(string)
	userIDInt, err := strconv.ParseUint(userID, 10, 0)
	if err != nil {
		logrus.Printf("Error while parsing user id:%s\n", err)
		return c.JSON(http.StatusBadRequest, "User id is wrong")
	}
	groups, err := database.GetUserGroups(uint(userIDInt))
	if groupName != "" {
		if err != nil {
			if errors.Is(err, database.NotFoundGroup) {
				return c.JSON(http.StatusBadRequest, "This group doesnt exist in your groups")
			} else {
				return c.JSON(http.StatusInternalServerError, map[string]interface{}{
					"error": "Failed to retrieve groups",
				})
			}
		}
		for _, group := range groups {
			if group.Name == groupName {
				profilePhotoFile, err := services.DownloadS3(services.StorageSession, configs.Config.StorageServiceBucket, group.GroupProfilePicture)
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
				group.GroupProfilePicture = base64ProfilePicture
				return c.JSON(http.StatusOK, group)
			}
		}
		return c.JSON(http.StatusBadRequest, "This group doesnt exist in your groups")

	} else {
		if err != nil {
			if errors.Is(err, database.NotFoundGroup) {
				return c.JSON(http.StatusOK, []model.Group{})
			} else {
				return c.JSON(http.StatusInternalServerError, map[string]interface{}{
					"error": "Failed to retrieve groups",
				})
			}
		}
		return c.JSON(http.StatusOK, groups)
	}
}
func RemoveUserGroup(c echo.Context) error {
	userID := c.Get("id").(string)
	userIDInt, err := strconv.ParseUint(userID, 10, 0)
	if err != nil {
		logrus.Printf("Error while parsing user id:%s\n", err)
		return c.JSON(http.StatusBadRequest, "User id is wrong")
	}

	groupID := c.FormValue("groupID")
	username := c.FormValue("username")
	groupIDInt, err := strconv.ParseUint(groupID, 10, 0)
	if err != nil {
		logrus.Printf("Error while parsing group id:%s\n", err)
		return c.JSON(http.StatusBadRequest, "Group id is wrong")
	}

	if username != "" {
		groupExist, err := database.DoesGroupExistByID(uint(userIDInt), uint(groupIDInt))
		if err != nil {
			logrus.Printf("Error while checking groups:%s\n", err)
			return c.JSON(http.StatusInternalServerError, "Can not check your groups")
		}
		if !groupExist {
			return c.JSON(http.StatusBadRequest, "You can not remove user from this group")
		}

		searchedUser, err := database.GetUserByUsername(username)
		if err != nil {
			if errors.Is(err, database.NotFoundUser) {
				return c.JSON(http.StatusBadRequest, fmt.Sprintf("No user found with username %s", username))
			}
			return c.JSON(http.StatusInternalServerError, "Error while finding user")
		}
		if searchedUser.ID == uint(userIDInt) {
			return c.JSON(http.StatusBadRequest, fmt.Sprintf("This user id is yours: %d", userIDInt))
		}

		isUserInGroup, err := database.IsUserInGroup(searchedUser.ID, uint(groupIDInt))
		if err != nil {
			logrus.Printf("Error while checking user group :%e", err)
			return c.JSON(http.StatusInternalServerError, "Can not check user group")
		}
		if !isUserInGroup {
			return c.JSON(http.StatusBadRequest, fmt.Sprintf("This user doesn't have group with group id %d", groupIDInt))
		}
		err = database.RemoveUserFromGroup(searchedUser.ID, uint(groupIDInt))
		if err != nil {
			if errors.Is(err, database.NotUserInGroup) {
				return c.JSON(http.StatusBadRequest, fmt.Sprintf("This user doesn't have group with group id %d", groupIDInt))
			}
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"error": "Failed to add user to group",
			})
		}
		return c.JSON(http.StatusCreated, "User removed from the group successfully.")
	} else {
		isUserInGroup, err := database.IsUserInGroup(uint(userIDInt), uint(groupIDInt))
		if err != nil {
			logrus.Printf("Error while checking user group :%e", err)
			return c.JSON(http.StatusInternalServerError, "Can not check user group")
		}
		if !isUserInGroup {
			return c.JSON(http.StatusBadRequest, fmt.Sprintf("You don't have group with group id %d", groupIDInt))
		}
		err = database.RemoveUserFromGroup(uint(userIDInt), uint(groupIDInt))
		if err != nil {
			if errors.Is(err, database.NotUserInGroup) {
				return c.JSON(http.StatusBadRequest, fmt.Sprintf("You don't have group with group id d %d", groupIDInt))
			}
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"error": "Failed to add user to group",
			})
		}
		return c.JSON(http.StatusCreated, "The group deleted successfully.")
	}

}
func GetAllUsersInGroup(c echo.Context) error {
	groupIDParam := c.QueryParam("groupID")
	groupID, err := strconv.ParseUint(groupIDParam, 10, 0)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"groupID": groupIDParam,
			"error":   err,
		}).Error("Invalid group ID")
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid group ID")
	}

	users, err := database.GetAllUsersInGroup(uint(groupID))
	if err != nil {
		logrus.Printf("Error while retriving users of a group")
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to retrieve users")
	}

	return c.JSON(http.StatusOK, users)
}
