package main

import (
	"Friend/DAO/Mysql"
	"Friend/DAO/MysqlTable"
	"Friend/DAO/Redis"
	settings "Friend/Settings"
	"Friend/rpc/Service"
	"Friend/rpc/client"
	"context"
	"fmt"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"rpc/handler/Init"
	"strconv"
)

// execPath returns the executable path.
func execPath() (string, error) {
	file, err := exec.LookPath(os.Args[0])
	if err != nil {
		return "", err
	}
	return filepath.Abs(file)
}
func main() {
	s, _ := execPath()
	fmt.Println("exec path:", s)

	Init.InitService()
	err := settings.Init()
	if err != nil {
		panic(err)
	}
	Redis.Init(settings.Config.RedisConfig)

	Service.InitServer(settings.Config.ServiceConfig)
	err = Mysql.Init(settings.Config.MysqlConfig)
	MysqlTable.InitTable(settings.Config.MysqlConfig)
	if err != nil {
		panic(err)
	}
	defer func() {
		err = Mysql.MysqlClose()
		if err != nil {
			panic(err)
		}
	}()
	//grpc服务开启监听
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

	// 其他服务发现
	client.Init(settings.Config.EtcdConfig.Host, settings.Config.EtcdConfig.Port)

	err = server.Serve(lis)
	if err != nil {
		fmt.Println(err)
		return
	}
}
