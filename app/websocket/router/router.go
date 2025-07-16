package router

import (
	"github.com/gin-gonic/gin"
	"websocket/server/controller"
)

func Init() *gin.Engine {
	r := gin.Default()
	r.Use(CORSMiddleware())
	r.GET("/api/v1/chat/ws", controller.WebSocketMiddleware())

	return r
}
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204) // 预检请求直接返回 204
			return
		}
		c.Next()
	}
}
