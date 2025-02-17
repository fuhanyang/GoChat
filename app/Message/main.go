package main

import (
	"Message/DAO/Mysql"
	"Message/DAO/MysqlTable"
	"Message/DAO/Redis"
	settings "Message/Settings"
	"Message/rpc/Service"
	"context"
	"rpc/handler/Init"
	"strconv"

	"fmt"
	"net"
)

func main() {
	Init.InitService()
	err := settings.Init()
	if err != nil {
		panic(err)
	}
	Redis.Init(settings.Config.RedisConfig)
	Service.InitServer(settings.Config.ServiceConfig)
	err = Mysql.Init(settings.Config.MysqlConfig)
	MysqlTable.InitTable()
	if err != nil {
		panic(err)
	}
	defer func() {
		err = Mysql.MysqlClose()
		if err != nil {
			panic(err)
		}
	}()

	lis, err := net.Listen(settings.Config.GrpcConfig.NetWork,
		settings.Config.GrpcConfig.Host+
			":"+
			strconv.Itoa(settings.Config.GrpcConfig.Port))
	if err != nil {
		panic(err)
	}
	ctx := context.Background()
	//服务注册
	server := Service.NewServer(ctx, settings.Config.EtcdConfig.Host, settings.Config.EtcdConfig.Port)

	err = server.Serve(lis)
	if err != nil {
		fmt.Println(err)
		return
	}
}
