package main

import (
	"fmt"
	"server/Mq"
	"server/Settings"
	"server/router"
)

func main() {
	var err error
	err = Settings.Init()
	if err != nil {
		panic(err)
	}
	err = Mq.NewConnCh(Settings.Config.RabbitMQConfig)
	if err != nil {
		panic(err)
	}
	fmt.Println("connect rabbitmq success at", Settings.Config.RabbitMQConfig.Host, Settings.Config.RabbitMQConfig.Port)
	defer Mq.ConnClose()
	r := router.Init()

	r.Run(":8080")
}
