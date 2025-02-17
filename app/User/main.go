package main

import (
	"User/DAO/Mysql"
	"User/DAO/MysqlTable"
	"User/DAO/Redis"
	"User/Logic/Snowflake"
	"User/Settings"
	"User/rpc/Service"
	"context"
	"fmt"
	"net"
	"rpc/handler/Init"
	"strconv"
)

func main() {
	Init.InitService()
	err := Settings.Init()
	if err != nil {
		panic(err)
	}
	Service.InitServer(Settings.Config.ServiceConfig)
	Redis.Init(Settings.Config.RedisConfig)
	err = Snowflake.Init(Settings.Config.SnowflakeConfig)
	if err != nil {
		panic(err)
	}
	err = Mysql.Init(Settings.Config.MysqlConfig)
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

	lis, err := net.Listen(Settings.Config.GrpcConfig.NetWork,
		Settings.Config.GrpcConfig.Host+
			":"+
			strconv.Itoa(Settings.Config.GrpcConfig.Port))
	if err != nil {
		panic(err)
	}
	fmt.Println("User service start success at", lis.Addr().String())
	ctx := context.Background()
	// 注册服务
	server := Service.NewServer(ctx, Settings.Config.EtcdConfig.Host, Settings.Config.EtcdConfig.Port)

	err = server.Serve(lis)
	if err != nil {
		fmt.Println(err)
		return
	}
}
