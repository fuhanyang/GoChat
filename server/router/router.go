package router

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"server/WebSocket"
	"server/controller"
	"server/middware"
)

func Init() *gin.Engine {
	r := gin.Default()
	r.Use(CORSMiddleware())
	r.LoadHTMLGlob("templates/*")
	// 静态文件服务
	r.Static("/api/static", "./static")
	v1 := r.Group("/api/v1")
	v1.GET("/start", func(c *gin.Context) {
		c.HTML(http.StatusOK, "user.html", gin.H{})
	})
	{
		userGroup := v1.Group("/user")
		{
			userGroup.POST("/register", controller.UserRegister())
			userGroup.POST("/login", controller.UserLogin())
			userGroup.POST("/logoff", middware.JWTAuthMiddleware(), controller.UserLogoff())
		}
		chatGroup := v1.Group("/chat")
		{
			chatGroup.GET("/ws", WebSocket.WebSocketMiddleware())
			chatGroup.GET("", func(c *gin.Context) {
				accountNum := c.Query("account_num")
				if accountNum == "" {
					c.Redirect(http.StatusFound, "/")
					return
				}
				c.HTML(http.StatusOK, "chat.html", gin.H{
					"account_num": accountNum,
				})
			})
			chatGroup.Use(middware.JWTAuthMiddleware())
			chatGroup.POST("/text/refresh", controller.RefreshText())
		}
		friendGroup := v1.Group("/friend")
		{
			friendGroup.Use(middware.JWTAuthMiddleware())
			friendGroup.POST("/addition", controller.AddFriend())
			//friendGroup.POST("/delete", controller.DeleteFriend())
			friendGroup.POST("/list", controller.GetFriends())
		}
	}
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
