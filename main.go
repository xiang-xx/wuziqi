package main

import (
	"github.com/gin-gonic/gin"
	"gobang/controller"
	"gobang/middleware"
	"fmt"
)

func main() {
	r := gin.Default()

	r.Use(middleware.Cors())
	r.GET("/game", func(c *gin.Context) {
		fmt.Println("hear")
		controller.WsHandler(c.Writer, c.Request)
	})
	r.POST("/new_room", controller.NewRoom)

	r.Run(":8000")
}