package WebSocket

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/streadway/amqp"
	"rpc/User"
	"server/Handler"
	"server/Mq"
	"server/random"
	"time"
)

// SendMessage 发送给消息队列并调用微服务
func SendMessage(message WebSocketData) error {
	var (
		req = &User.FindUserRequest{
			AccountNum:  message.Receiver,
			HandlerName: "FindUser",
		}
		v      []byte
		err    error
		result = User.FindUserResponse{}
		resp   []byte
		msgs   <-chan amqp.Delivery
		corrId string
	)
	data, err := ParseData([]byte(message.Data), message.Type, message.CreatedAt)
	if err != nil {
		fmt.Println("Error while parsing data")
		return err
	}
	// 调用微服务查找接收者是否存在
	v, err = json.Marshal(req)
	if err != nil {
		return err
	}
	corrId = random.RandomString(32)
	msgs, err = Mq.PublishMessage(corrId, v, "", req.GetHandlerName(), false, false)
	if err != nil {
		return err
	}
	for {
		select {
		case d := <-msgs:
			if corrId == d.CorrelationId {
				resp = d.Body
				err = d.Ack(false)
				if err != nil {
					return err
				}
				goto Correction
			}
		case <-time.NewTimer(time.Second * 3).C:
			return errors.New("Timeout while waiting for response")
		}
	}
Correction:
	err = json.Unmarshal(resp, &result)
	if err != nil {
		return err
	}
	if result.Code != 200 {
		return errors.New("Receiver not found")
	}
	corrId = random.RandomString(32)
	msgs, err = Mq.PublishMessage(corrId, data, "", message.Type, false, false)
	if err != nil {
		return err
	}
	for {
		select {
		case d := <-msgs:
			if corrId == d.CorrelationId {
				return d.Ack(false)
			}
		case <-time.NewTimer(time.Second * 3).C:
			return errors.New("Timeout while waiting for response")
		}
	}
}

// ParseData 解析数据到HandlerWithTime实例并转化成json字节
func ParseData(data []byte, messageType string, CreatedAt string) ([]byte, error) {
	handler := Handler.GetHandlersType(messageType)
	if handler == nil {
		return nil, fmt.Errorf("Client not found for message type %s", messageType)
	}
	err := json.Unmarshal(data, &handler)
	if err != nil {
		return nil, err
	}
	return json.Marshal(Handler.HandlerWithTime{Handler: handler, CreatedAt: CreatedAt, Type: messageType})
}
