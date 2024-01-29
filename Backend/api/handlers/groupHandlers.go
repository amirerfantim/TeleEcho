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
		fmt.Printf("Error while parsing user id:%s\n", err)
		return c.JSON(http.StatusBadRequest, "User id is wrong")
	}
	groupName := c.FormValue("name")
	groupExist, err := database.DoesGroupExist(uint(userIDInt), groupName)
	if err != nil {
		fmt.Printf("Error while checking groups:%s\n", err)
		return c.JSON(http.StatusInternalServerError, "Can not check your groups")
	}
	if !groupExist {
		return c.JSON(http.StatusBadRequest, "This group with your user id is already exist")
	}
	groupDescription := c.FormValue("description")
	profilePicture, err := c.FormFile("profile")
	if err != nil {
		logrus.Printf("Unable to open image\n")
		return c.String(http.StatusBadRequest, "Unable to open file")
	}
	profilePath, err := services.UploadS3(services.StorageSession, profilePicture, configs.Config.StorageServiceBucket, userID+"/"+groupName)
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
func GetUserGroups(c echo.Context) error {
	groupName := c.FormValue("name")
	userID := c.Get("id").(string)
	userIDInt, err := strconv.ParseUint(userID, 10, 0)
	if err != nil {
		fmt.Printf("Error while parsing user id:%s\n", err)
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
					logrus.Println("Can not open file", err)
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
