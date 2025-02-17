package ParseMessage

import (
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
	"rpc/handler"
)

// ParseDelivery 解析消息
func ParseDelivery(v handler.HandlerRequest, delivery amqp.Delivery) (handler.HandlerRequest, error) {
	if v == nil {
		return nil, fmt.Errorf("nil handler")
	}
	body := delivery.Body
	if len(body) == 0 {
		return nil, fmt.Errorf("empty body")
	}
	fmt.Println("request body:", string(body))
	err := json.Unmarshal(body, &v)
	return v, err
}
