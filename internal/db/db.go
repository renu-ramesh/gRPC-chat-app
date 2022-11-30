package db

import (
	"chat_app_grpc/internal/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {

	database, err := gorm.Open(sqlite.Open("Chat_App_DB.db"), &gorm.Config{})

	if err != nil {
		panic("Failed to connect to database!")
	}

	err = database.AutoMigrate(&models.User{})
	if err != nil {
		return
	}

	err = database.AutoMigrate(&models.Channel{})
	if err != nil {
		return
	}

	err = database.AutoMigrate(&models.UserChannelDetails{})
	if err != nil {
		return
	}

	DB = database
}
