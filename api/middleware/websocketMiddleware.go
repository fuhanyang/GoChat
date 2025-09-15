package middleware

import (
	"api/rpc/client"
	"common/jwt"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	Websocket "rpc/websocket"
)

// WebSocketMiddleware
// websocket中间件，每次请求时判断是否建立了websocket连接,如果没有则创建连接并返回,如果已经建立连接则直接返回
func WebSocketMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			err error
			mc  *jwt.CustomClaims
		)
		fmt.Println(" start websocket middleware")
		token := c.Query("token")
		// parts[1]是获取到的tokenString，使用定义好的解析JWT的函数来解析它
		mc, err = jwt.ParseToken(token)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code": 2005,
				"msg":  "无效的Token",
			})
			c.Abort()
			return
		}

		if mc.AccountNum == "" {
			c.Abort()
			fmt.Println("account_num is empty")
			return
		}
		//调用微服务尝试获取一个websocket服务器
		connect, err := client.WebsocketServiceClient.Client.TryConnect(context.Background(), &Websocket.TryConnectRequest{
			AccountNum: mc.AccountNum,
			Password:   mc.Password,
		})
		if err != nil {
			c.JSON(400, gin.H{
				"code": 2005,
				"msg":  err.Error(),
			})
			c.Abort()
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"URL":   connect.Url,
			"Token": connect.Token,
		})
	}
}
