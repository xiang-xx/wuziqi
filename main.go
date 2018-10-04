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

	r.Run(":8000")
}