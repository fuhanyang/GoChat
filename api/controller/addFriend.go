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
	"rpc/friend"
	"time"
)

func AddFriend() gin.HandlerFunc {
	return func(c *gin.Context) {
		type _friend struct {
			AccountNum string `json:"accountNum"`
			Name       string `json:"name"`
		}
		var (
			req          = friend.AddFriendRequest{}
			v            []byte
			err          error
			result       = friend.AddFriendResponse{}
			resp         []byte
			msgs         <-chan amqp.Delivery
			corrId       string
			respCh       = make(chan []byte, 1)
			closedSignal = make(chan bool, 1)
			timer        *time.Timer
			friendData   = _friend{}
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
		msgs, err = Mq.PublishMessage(corrId, v, "", req.GetHandlerName(), false, false)
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

		// 接收mq的响应,超时时间为3秒,超时则返回错误信息
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
		if result.Code != http.StatusOK {
			err = errors.New(result.GetMsg())
			goto ERR
		}
		friendData.AccountNum = result.GetFriend().GetAccountNum()
		friendData.Name = result.GetFriend().GetName()
		c.JSON(http.StatusOK, gin.H{
			"code": result.GetCode(),
			"msg":  result.GetMsg(),
			"data": friendData,
		})
		return
	ERR:
		c.JSON(500, gin.H{
			"code":  500,
			"msg":   fmt.Sprintf("Add _friend error :%s", result.GetMsg()),
			"error": err.Error(),
		})
	}
}

func AddFriendWithAccountNum() gin.HandlerFunc {
	return func(c *gin.Context) {
		type _friend struct {
			AccountNum string `json:"accountNum"`
			Name       string `json:"name"`
		}
		var (
			req          = friend.AddFriendWithAccountNumRequest{}
			v            []byte
			err          error
			result       = friend.AddFriendResponse{}
			resp         []byte
			msgs         <-chan amqp.Delivery
			corrId       string
			respCh       = make(chan []byte, 1)
			closedSignal = make(chan bool, 1)
			timer        *time.Timer
			friendData   = _friend{}
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
		msgs, err = Mq.PublishMessage(corrId, v, "", req.GetHandlerName(), false, false)
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

		// 接收mq的响应,超时时间为3秒,超时则返回错误信息
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
		if result.Code != http.StatusOK {
			err = errors.New(result.GetMsg())
			goto ERR
		}
		friendData.AccountNum = result.GetFriend().GetAccountNum()
		friendData.Name = result.GetFriend().GetName()
		c.JSON(http.StatusOK, gin.H{
			"code": result.GetCode(),
			"msg":  result.GetMsg(),
			"data": friendData,
		})
		return
	ERR:
		c.JSON(500, gin.H{
			"code":  500,
			"msg":   fmt.Sprintf("Add _friend error :%s", result.GetMsg()),
			"error": err.Error(),
		})
	}
}
