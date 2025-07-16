package Logic

import (
	"fmt"
	"github.com/streadway/amqp"
	"sync"
	"websocket/settings"
)

var MachineURL string

var (
	conn *amqp.Connection
	ch   *amqp.Channel
	once sync.Once
)

const (
	BExchange = "b_exchange"
	DExchange = "d_exchange"
)

func ExchangeInit() {
	once.Do(func() {
		err := ch.ExchangeDeclare(
			BExchange,
			"fanout",
			true,
			false,
			false,
			false,
			nil,
		)
		if err != nil {
			panic(err)
		}
		err = ch.ExchangeDeclare(
			DExchange,
			"direct",
			true,
			false,
			false,
			false,
			nil,
		)
		if err != nil {
			panic(err)
		}

	})
}
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

func BroadcastMessage(data []byte, mandatory bool, immediate bool) error {
	ExchangeInit()
	err := ch.Publish(
		BExchange,
		"",
		mandatory,
		immediate,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        data,
		},
	)
	return err
}
func StartConsume(queueName string, autoAck bool) <-chan amqp.Delivery {
	msg, err := ch.Consume(
		queueName,
		"",
		autoAck,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		panic(err)
		return nil
	}
	return msg
}

func DirectionMessage(data []byte, direction string, mandatory bool, immediate bool) error {
	err := ch.Publish(
		DExchange,
		direction,
		mandatory,
		immediate,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        data,
		})
	return err
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
func QueueBind(queueName string, exchangeName string, routingKey string, noWait bool, args amqp.Table) error {
	err := ch.QueueBind(
		queueName,
		routingKey,
		exchangeName,
		noWait,
		args,
	)
	return err
}
