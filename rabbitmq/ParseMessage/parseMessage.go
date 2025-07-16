package ParseMessage

import (
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
	"rpc/Service/method"
)

// ParseDelivery 解析消息
func ParseDelivery(v method.MethodRequest, delivery amqp.Delivery) (method.MethodRequest, error) {
	if v == nil {
		return nil, fmt.Errorf("nil method")
	}
	body := delivery.Body
	if len(body) == 0 {
		return nil, fmt.Errorf("empty body")
	}
	fmt.Println("request body:", string(body))
	err := json.Unmarshal(body, &v)
	return v, err
}
