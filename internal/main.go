package main

import (
	"chat_app_grpc/internal/controllers"
	"chat_app_grpc/internal/db"

	"github.com/gin-gonic/gin"
)

func main() {

	r := gin.Default()
	db.ConnectDatabase()

	// API-Version1
	v1 := r.Group("/api/v1")

	// Channel Methods
	channel := v1.Group("/channel")
	channel.POST("/", controllers.CreateChannel)
	channel.GET("/", controllers.FindChannels)
	channel.DELETE(":id/", controllers.DeleteChannel)

	// User Methods
	users := v1.Group("/users")
	users.POST("/", controllers.CreateUser)
	users.GET("/", controllers.FindUsers)
	users.POST(":id/join", controllers.JoinChannel)
	users.PUT(":id/left", controllers.LeftChannel)
	users.GET("/channels", controllers.UsersChannels)

	err := r.Run()
	if err != nil {
		return
	}
}
