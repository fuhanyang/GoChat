package controller

import (
	"common/jwt"
	"common/redis"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
	"websocket/DAO/Redis"
	"websocket/Logic"
	websocket2 "websocket/Logic/websocket"
	"websocket/Models"
)

// WebSocketMiddleware
// websocket中间件，每次请求时判断是否建立了websocket连接,如果没有则创建连接并返回,如果已经建立连接则直接返回
func WebSocketMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			err  error
			conn *websocket.Conn
			ws   *Models.WsConn
			mc   *jwt.CustomClaims
		)
		token := c.Query("token")
		// parts[1]是获取到的tokenString，我们使用之前定义好的解析JWT的函数来解析它
		mc, err = jwt.ParseToken(token)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code": 2005,
				"msg":  "无效的Token",
			})
			c.Abort()
			return
		}
		// 将当前请求的username信息保存到请求的上下文c上
		sender := mc.AccountNum

		if sender == "" {
			c.Abort()
			fmt.Println("account_num is empty")
			return
		}

		if conn, err = websocket2.Upgrade.Upgrade(c.Writer, c.Request, nil); err != nil {
			fmt.Println(err)
			return
		}
		if ws, err = websocket2.InitWebSocket(conn, sender); err != nil {
			fmt.Println(err)
			return
		}
		websocket2.WriteWebSocketMap(sender, ws)

		// 存入redis
		rconn := Redis.RedisPoolGet()
		redis.HmsetWithExpireScript.Do(rconn, fmt.Sprintf("websocket_%s", ws.AccountNum), 360, "url", Logic.MachineURL)
		Redis.RedisPoolPut(rconn)
		//Redis.RedisDo(Redis2.HmsetWithExpireScriptString, fmt.Sprintf("websocket_%s", ws.AccountNum), 360, "url", Logic.MachineURL)
		// 使得inChan和outChan耦合起来
		// 管道关闭则关闭WebSocket连接
		go func() {
			for {
				var data Models.WebSocketData
				if data, err = ws.InChanRead(); err != nil {
					fmt.Println(err)
					goto ERR
				}
				if err = ws.OutChanWrite(data); err != nil {
					fmt.Println(err)
					goto ERR
				}
			}
		ERR:
			ws.CloseConn()
		}()
		return
	}
}
