package queue

import (
	"fmt"
	"github.com/streadway/amqp"
	"rabbitmq/settings"
)

var conn *amqp.Connection
var ch *amqp.Channel

func NewConnCh(config *settings.RabbitMQConfig) error {
	var err error
	dsn := fmt.Sprintf("amqp://%s:%s@%s:%d/", config.Username, config.Password, config.Host, config.Port)
	conn, err = amqp.Dial(dsn)
	if err != nil {
		return err
	}
	ch, err = conn.Channel()
	return err
}
func ConnClose() {
	ch.Close()
	conn.Close()
}

// QueueDeclare 声明队列
func QueueDeclare(queueName string, durable bool, autoDelete bool, exclusive bool, noWait bool, args amqp.Table) (amqp.Queue, error) {
	q, err := ch.QueueDeclare(
		queueName,
		durable,
		autoDelete,
		exclusive,
		noWait,
		args,
	)
	return q, err
}

// GivesResponseTo 给指定队列发送响应消息
func GivesResponseTo(key string, correlationId string, content []byte) error {
	fmt.Println("gives response to content:", string(content), " correlationId:", correlationId)
	err := ch.Publish(
		"",
		key,
		false,
		false,
		amqp.Publishing{
			ContentType:   "text/plain",
			CorrelationId: correlationId,
			Body:          content,
		},
	)
	return err
}

// GetMessages 获取消息队列的管道
func GetMessages(queueName string, consumer string, autoAck bool, exclusive bool, noLocal bool, noWait bool, args amqp.Table) (<-chan amqp.Delivery, error) {
	msgs, err := ch.Consume(
		queueName,
		consumer,
		autoAck,
		exclusive,
		noLocal,
		noWait,
		args,
	)
	return msgs, err
}
