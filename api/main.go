package main

import (
	"api/Mq"
	"api/router"
	"api/rpc/client"
	"api/settings"
	"common/viper"
	"fmt"
	"rpc/Service/Init"
)

func main() {
	Init.InitService()
	var err error
	err = viper.Init(settings.Config)
	if err != nil {
		panic(err)
	}
	fmt.Println("config init success mode:", settings.Config.Mode)
	err = Mq.NewConnCh(settings.Config.RabbitMQConfig)
	if err != nil {
		panic(err)
	}
	fmt.Println("connect rabbitmq success at", settings.Config.RabbitMQConfig.Host, settings.Config.RabbitMQConfig.Port)
	defer Mq.ConnClose()

	client.Init(settings.Config.EtcdConfig.Host, settings.Config.EtcdConfig.Port)

	r := router.Init()
	dsn := fmt.Sprintf("%s:%d", settings.Config.Host, settings.Config.Port)
	fmt.Println("api run at", dsn)
	r.Run(dsn)
}
