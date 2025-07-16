package router

import (
	"api/controller"
	"api/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Init() *gin.Engine {
	r := gin.Default()
	r.Use(CORSMiddleware())
	r.LoadHTMLGlob("templates/*")
	// 静态文件服务
	r.Static("/api/static", "./static")
	v1 := r.Group("/api/v1")
	v1.GET("/start", controller.LoadStart())
	v1.Use(middleware.LimitMiddleware())
	{
		userGroup := v1.Group("/user")
		{
			userGroup.POST("/register", controller.UserRegister())
			userGroup.POST("/login", controller.UserLogin())
			userGroup.POST("/logoff", middleware.JWTAuthMiddleware(), controller.UserLogoff())
			userGroup.GET("/info", middleware.JWTAuthMiddleware(), controller.GetUserInfo())
			userGroup.POST("/upload-avatar", middleware.JWTAuthMiddleware(), controller.UploadAvatar)
			userGroup.GET("/info/start", func(c *gin.Context) {
				c.HTML(http.StatusOK, "userInfo.html", gin.H{})
			})
		}
		chatGroup := v1.Group("/chat")
		{
			chatGroup.GET("/ws", middleware.WebSocketMiddleware())
			chatGroup.GET("", controller.LoadChat())
			chatGroup.Use(middleware.JWTAuthMiddleware())
			chatGroup.POST("/text/refresh", controller.RefreshText())
		}
		friendGroup := v1.Group("/friend")
		{
			friendGroup.Use(middleware.JWTAuthMiddleware())
			friendGroup.POST("/addition", controller.AddFriend())
			friendGroup.POST("/addition/with_account_num", controller.AddFriendWithAccountNum())
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
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204) // 预检请求直接返回 204
			return
		}
		c.Next()
	}
}
