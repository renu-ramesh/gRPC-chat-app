package controllers

import (
	"chat_app_grpc/internal/common"
	"chat_app_grpc/internal/db"
	"chat_app_grpc/internal/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Create a new user
func CreateUser(c *gin.Context) {

	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Check user already exist
	if exists, _ := common.CheckUserExistByName(user.Name); exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User already exist!"})
		return
	}
	// Create User
	userData := models.User{Name: user.Name}
	db.DB.Create(&userData)

	c.JSON(http.StatusOK, gin.H{})
}

// List all users
func FindUsers(c *gin.Context) {
	var user []models.User
	db.DB.Find(&user)

	c.JSON(http.StatusOK, gin.H{"data": user})
}

// create new group/channel
func CreateChannel(c *gin.Context) {

	var channel models.Channel
	if err := c.ShouldBindJSON(&channel); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Check channel already exist
	count := int64(0)
	if err := db.DB.Model(&models.Channel{}).
		Where("name = ? ", channel.Name).
		Count(&count).
		Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if exists := count > 0; exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Channel already exist!"})
		return
	}
	// Create channel
	ChData := models.Channel{Name: channel.Name, Type: channel.Type}
	db.DB.Create(&ChData)

	c.JSON(http.StatusOK, gin.H{})
}

// List all groups/channels
func FindChannels(c *gin.Context) {

	var channels []models.Channel
	db.DB.Where("deleted_status = ?", false).Find(&channels)

	c.JSON(http.StatusOK, gin.H{"data": channels})
}

// Delete a channel
func DeleteChannel(c *gin.Context) {

	// check channel exist
	var channel models.Channel
	if err := db.DB.Where("id = ?", c.Param("id")).First(&channel).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}

	ChData := models.Channel{UpdatedAt: time.Now(), DeletedStatus: true}
	db.DB.Model(&ChData).Where("id = ?", c.Param("id")).Updates(ChData)

	c.JSON(http.StatusOK, gin.H{})
}

// new user join to a channel/group
func JoinChannel(c *gin.Context) {
	// Validate input
	var inputData models.UserChannelDetails
	if err := c.ShouldBindJSON(&inputData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Check user already exist
	count := int64(0)
	var user models.User
	if err := db.DB.Where("id = ?", c.Param("id")).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User Record not found!"})
		return
	}
	// check channel exist
	var channel models.Channel
	if err := db.DB.Where("name = ?", inputData.ChannelName).First(&channel).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "channel Record not found!"})
		return
	}

	// check user already joined the channel
	if err := db.DB.Model(&models.UserChannelDetails{}).
		Where("user_name = ? AND channel_name = ?", user.Name, inputData.ChannelName).
		Count(&count).
		Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if channel_exists := count > 0; channel_exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User already joined"})
		return
	}

	ChData := models.UserChannelDetails{UserName: user.Name, ChannelName: inputData.ChannelName}
	db.DB.Create(&ChData)

	c.JSON(http.StatusOK, gin.H{})
}

// Left from a channel/group
func LeftChannel(c *gin.Context) {

	var inputData models.UserChannelDetails
	if err := c.ShouldBindJSON(&inputData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := db.DB.Where("id = ?", c.Param("id")).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User Record not found!"})
		return
	}
	var userDetails models.UserChannelDetails
	if err := db.DB.Where("user_name = ? AND channel_name = ?", user.Name, inputData.ChannelName).First(&userDetails).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User Record not found!"})
		return
	}

	var UserChanels models.UserChannelDetails
	db.DB.Model(&UserChanels).Where("id = ?", userDetails.ID).Updates(models.UserChannelDetails{DeletedStatus: true})

	// Delete a channel if all the users are left
	count := int64(0)
	if err := db.DB.Model(&models.UserChannelDetails{}).
		Where("deleted_status = ? ", false).
		Count(&count).
		Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if users_exists := count > 0; !users_exists {
		ChData := models.Channel{UpdatedAt: time.Now(), DeletedStatus: true}
		db.DB.Model(&ChData).Where("id = ?", c.Param("id")).Updates(ChData)
	}

	c.JSON(http.StatusOK, gin.H{})
}

// List users's joined channel details
func UsersChannels(c *gin.Context) {

	var channels []models.UserChannelDetails
	db.DB.Where("deleted_status = ?", false).Find(&channels)

	c.JSON(http.StatusOK, gin.H{"data": channels})
}
