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
		//var (
		//	err  error
		//	conn *websocket.Conn
		//	ws   *wsConn
		//	mc   *jwt.CustomClaims
		//)
		//fmt.Println(" start websocket middleware")
		//token := c.Query("token")
		//// parts[1]是获取到的tokenString，我们使用之前定义好的解析JWT的函数来解析它
		//mc, err = jwt.ParseToken(token)
		//if err != nil {
		//	c.JSON(http.StatusOK, gin.H{
		//		"code": 2005,
		//		"msg":  "无效的Token",
		//	})
		//	c.Abort()
		//	return
		//}
		//// 将当前请求的username信息保存到请求的上下文c上
		//sender := mc.AccountNum
		//fmt.Println("sender:", sender, "password:", mc.Password)
		//if sender == "" {
		//	c.Abort()
		//	fmt.Println("account_num is empty")
		//	return
		//}
		//
		//// 已经建立连接则直接返回
		//if exist, _ := ReadWebSocketMap(sender); exist {
		//	c.Abort()
		//	c.JSON(200, gin.H{
		//		"code": 200,
		//		"msg":  "account_num is exist",
		//	})
		//	return
		//}
		//
		//if conn, err = upgrade.Upgrade(c.Writer, c.Request, nil); err != nil {
		//	fmt.Println(err)
		//	return
		//}
		//if ws, err = InitWebSocket(conn, sender); err != nil {
		//	fmt.Println(err)
		//	return
		//}
		//// 使得inChan和outChan耦合起来
		//// 管道关闭则关闭WebSocket连接
		//go func() {
		//	for {
		//		var data WebSocketData
		//		if data, err = ws.InChanRead(); err != nil {
		//			fmt.Println(err)
		//			goto ERR
		//		}
		//		if err = ws.OutChanWrite(data); err != nil {
		//			fmt.Println(err)
		//			goto ERR
		//		}
		//	}
		//ERR:
		//	ws.CloseConn()
		//}()
		//return
	}
}
