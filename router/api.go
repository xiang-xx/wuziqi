package router

import (
	"github.com/gin-gonic/gin"
	"gobang/middleware"
)

func ApiRouter(r *gin.Engine){
	r.Use(middleware.Cors())
}