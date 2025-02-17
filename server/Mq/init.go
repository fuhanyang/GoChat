package Mq

import (
	"fmt"
	"github.com/streadway/amqp"
	"server/Settings"
	"sync"
)

var conn *amqp.Connection
var ch *amqp.Channel
var once sync.Once
var respQueue amqp.Queue
var err error
var msgs <-chan amqp.Delivery

func respQueueInit() {
	respQueue, err = ch.QueueDeclare(
		"",    // 不指定队列名，默认使用随机生成的队列名
		false, // durable
		false, // delete when unused
		false, // exclusive
		false, // noWait
		nil,   // arguments
	)
	if err != nil {
		panic(err)
	}
	msgs, err = ch.Consume(
		respQueue.Name, // queue
		"",             // consumer
		false,          // auto-ack
		false,          // exclusive
		false,          // no-local
		false,          // no-wait
		nil,            // args
	)
	if err != nil {
		panic(err)
	}
}
func NewConnCh(config *Settings.RabbitMQConfig) error {
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

func PublishMessage(corrId string, data []byte, exchangeName string, key string, mandatory bool, immediate bool) (<-chan amqp.Delivery, error) {
	respQueueInit()
	err = ch.Publish(
		exchangeName,
		key,
		mandatory,
		immediate,
		amqp.Publishing{
			ContentType:   "text/plain",
			ReplyTo:       respQueue.Name,
			CorrelationId: corrId,
			Body:          data,
		},
	)
	if err != nil {
		return nil, err
	}
	return msgs, nil
}
