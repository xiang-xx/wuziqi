package main

import (
	"github.com/gin-gonic/gin"
	"gobang/controller"
	"gobang/middleware"
)

func main() {
	r := gin.Default()

	r.Use(middleware.Cors())
	r.GET("/game", func(c *gin.Context) {
		controller.WsHandler(c.Writer, c.Request)
	})
	r.POST("/new_room", controller.NewRoom)

	r.Run(":8000")
}