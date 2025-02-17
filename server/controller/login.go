package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/streadway/amqp"
	"net/http"
	"rpc/User"
	"server/Jwt"
	"server/Mq"
	"server/WebSocket"
	"server/random"
	"time"
)

func UserLogin() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			req          = User.LoginRequest{}
			v            []byte
			err          error
			result       = User.LoginResponse{}
			resp         []byte
			msgs         <-chan amqp.Delivery
			corrId       string
			respCh       = make(chan []byte, 1)
			closedSignal = make(chan bool, 1)
			timer        *time.Timer
			token        string
		)
		defer func() {
			close(respCh)
			closedSignal <- true
			close(closedSignal)
		}()

		//
		err = c.BindJSON(&req)
		if err != nil {
			goto ERR
		}
		v, err = json.Marshal(&req)
		if err != nil {
			goto ERR
		}

		// 发送消息到mq
		corrId = random.RandomString(32)
		msgs, err = Mq.PublishMessage(corrId, v, "", "Login", false, false)
		if err != nil {
			goto ERR
		}
		// 开启一个goroutine监听mq的响应
		go func() {
			defer func() {
				if err := recover(); err != nil {
					fmt.Println(err)
				}
			}()
			for d := range msgs {
				fmt.Println("收到mq的响应")
				if d.CorrelationId == corrId {
					// 响应id匹配，确认消息
					err = d.Ack(false)
					if err != nil {
						panic(err)
					}
					select {
					case <-closedSignal:
						//超时，需要从websocket发送响应给用户
						exist, ws := WebSocket.ReadWebSocketMap(req.GetAccountNum())
						if !exist {
							return
						}
						var wsData = WebSocket.WebSocketData{}
						wsData.Data = string(d.Body)
						wsData.Type = "LoginResponse"
						wsData.AccountNum = req.GetAccountNum()
						wsData.Receiver = req.GetAccountNum()
						wsData.CreatedAt = time.Now().Format(time.RFC3339)
						err = ws.OutChanWrite(wsData)
						if err != nil {
							fmt.Println(err)
						}
						return
					case respCh <- d.Body:
						return
					}
				}
			}
		}()

		// 接收mq的响应,超时时间为3秒,超时则后续由websocket发送响应给用户
		timer = time.NewTimer(time.Second * 3)
		select {
		case resp = <-respCh:
			goto Correction
		case <-timer.C:
			err = errors.New("服务器繁忙，请稍后再试")
			goto ERR
		}

	Correction:
		err = json.Unmarshal(resp, &result)
		if err != nil {
			goto ERR
		}
		if result.Code != http.StatusOK {
			goto ERR
		}
		token, err = Jwt.GenToken(req.GetPassword(), req.GetAccountNum())
		if err != nil {
			goto ERR
		}
		c.Header("Authorization", token)
		fmt.Println("token:", token)
		c.JSON(http.StatusOK, gin.H{
			"code": result.GetCode(),
			"msg":  result.GetMsg(),
		})
		return
	ERR:
		fmt.Println(result.GetMsg())
		c.JSON(500, gin.H{
			"code":  500,
			"msg":   "User Login error",
			"error": result.GetMsg(),
		})
	}
}
