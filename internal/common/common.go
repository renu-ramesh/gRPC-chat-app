package common

import (
	"chat_app_grpc/internal/db"
	"chat_app_grpc/internal/models"
	"errors"
)

func ValidateUser(username string) (err error) {

	var user models.User
	if err := db.DB.Where("name = ?", username).First(&user).Error; err != nil {
		return errors.New("user not found")
	}

	return nil
}
func ValidateUserChannel(username, channel_name string) (err error) {

	var userChanData models.UserChannelDetails
	if err := db.DB.Where("user_name = ? AND channel_name = ? AND deleted_status = ?", username, channel_name, false).First(&userChanData).Error; err != nil {
		return errors.New("user can't send message to this group")
	}
	return nil
}
func CheckUserExistByName(username string) (bool, error) {
	count := int64(0)
	if err := db.DB.Model(&models.User{}).
		Where("name = ? ", username).
		Count(&count).
		Error; err != nil {

		return false, err
	}

	return count > 0, nil
}
