package controller

import (
	"api/Mq"
	"api/random"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/streadway/amqp"
	"net/http"
	"rpc/message"
	"time"
)

func RefreshText() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			req          = message.RefreshRequest{}
			v            []byte
			err          error
			result       = message.RefreshTextResponse{}
			resp         []byte
			msgs         <-chan amqp.Delivery
			corrId       string
			respCh       = make(chan []byte, 1)
			closedSignal = make(chan bool, 1)
			timer        *time.Timer
		)
		defer func() {
			close(respCh)
			closedSignal <- true
			close(closedSignal)
		}()

		type Text struct {
			Content  []*message.Text `json:"content"`
			Sender   string          `json:"sender"`
			Receiver string          ` json:"receiver"`
		}
		var text Text

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
		msgs, err = Mq.PublishMessage(corrId, v, "", "RefreshText", false, false)
		if err != nil {
			goto ERR
		}
		// 开启一个goroutine监听mq的响应
		// 不加入corrId会导致闭包引用问题
		go func() {
			defer func() {
				if err := recover(); err != nil {
					fmt.Println(err)
				}
			}()
			for d := range msgs {
				if d.CorrelationId == corrId {
					// 响应id匹配，确认消息
					err = d.Ack(false)
					if err != nil {
						panic(err)
					}
					select {
					case <-closedSignal:
						//超时，需要调用websocket微服务去处理
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
		if result.GetCode() != http.StatusOK {
			c.JSON(int(result.GetCode()), gin.H{
				"code":  result.GetCode(),
				"msg":   result.GetMsg(),
				"error": "RefreshText failed",
			})
			return
		}
		text = Text{
			Content:  result.GetContent(),
			Sender:   result.GetSenderAccountNum(),
			Receiver: result.GetReceiverAccountNum(),
		}
		fmt.Println(text)
		c.JSON(http.StatusOK, gin.H{
			"code": result.GetCode(),
			"msg":  result.GetMsg(),
			"text": text,
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
