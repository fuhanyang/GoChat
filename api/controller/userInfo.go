package controller

import (
	"api/Mq"
	"api/random"
	"encoding/json"
	"errors"
	"fmt"
	"rpc/user"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/streadway/amqp"
)

type info struct {
	AccountNum string `json:"account_num"`
	Nickname   string `json:"name"`
	IP         string `json:"ip"`
	Email      string `json:"email"`
	CreateAt   string `json:"create_at"`
}

func GetUserInfo() gin.HandlerFunc {
	return func(c *gin.Context) {

		var (
			req          = user.GetUserInfoRequest{}
			v            []byte
			err          error
			result       = user.GetUserInfoResponse{}
			resp         []byte
			msgs         <-chan amqp.Delivery
			corrId       string
			respCh       = make(chan []byte, 1)
			closedSignal = make(chan bool, 1)
			timer        *time.Timer
			data         info
		)
		defer func() {
			close(respCh)
			closedSignal <- true
			close(closedSignal)
		}()
		accountNum := c.Query("account_num")
		if accountNum == "" {
			err = errors.New("account_num不能为空")
			goto ERR
		}
		fmt.Println("get user info, account_num:", accountNum)
		req.AccountNum = accountNum
		req.HandlerName = "GetUserInfo"
		v, err = json.Marshal(&req)
		if err != nil {
			goto ERR
		}
		// 发送消息到mq
		corrId = random.RandomString(32)
		msgs, err = Mq.PublishMessage(corrId, v, "", "GetUserInfo", false, false)
		if err != nil {
			goto ERR
		}

		// 开启一个goroutine监听mq的响应
		go func() {
			defer func() {
				fmt.Println("mq监听goroutine退出")
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
						// 超时，记入日志
						//TODO

						return
					case respCh <- d.Body:
						return
					}
				} else {
					err = d.Nack(false, true)
					if err != nil {
						panic(err)
					}
				}
			}
		}()

		// 接收mq的响应,超时时间为3秒,超时则后续由websocket发送响应给用户
		timer = time.NewTimer(time.Second * 3)
		defer timer.Stop()
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
		data = info{
			AccountNum: result.GetAccountNum(),
			Nickname:   result.GetUsername(),
			IP:         result.GetIp(),
			Email:      result.GetEmail(),
			CreateAt:   result.GetCreateAt(),
		}
		c.JSON(int(result.GetCode()), gin.H{
			"code": result.GetCode(),
			"msg":  result.GetMsg(),
			"data": data,
		})
		return
	ERR:
		fmt.Println(err)
		c.JSON(500, gin.H{
			"code":  500,
			"msg":   err.Error(),
			"error": err.Error(),
		})
	}
}
