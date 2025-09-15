package main

import (
	"common/mysql"
	"common/redis"
	"common/viper"
	"common/zap"
	"context"
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"rpc/Service/Init"
	"strconv"
	"websocket/DAO/Mysql"
	"websocket/DAO/Redis"
	"websocket/Logic"
	"websocket/Logic/websocket"
	"websocket/router"
	"websocket/rpc/client"
	"websocket/rpc/service"
	settings "websocket/settings"
)

var (
	lis    net.Listener
	err    error
	server *grpc.Server
	defers = make([]func(), 0)
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
	mode := flag.String("mode", "local", "运行模式，可选值：local(默认)、debug")
	flag.Parse()

	err = viper.Init(settings.Config, *mode)
	if err != nil {
		panic(err)
	}
	fmt.Println("config init success mode:", *mode)

	defer func() {
		for _, f := range defers {
			f()
		}
	}()

	websocketInit()
	listen()
	findService()
	go serve()
	websocket.InitWebsocketMQ()

	r := router.Init()
	r.Run(Logic.MachineURL)
}
func websocketInit() {
	Logic.MachineURL = fmt.Sprintf("%s:%d", settings.Config.App.Host, settings.Config.App.Port)
	fmt.Println("MachineURL:", Logic.MachineURL)

	Redis.Pool = redis.Init(settings.Config.RedisConfig)

	Mysql.MysqlDb, err = mysql.Init(settings.Config.MysqlConfig)
	if err != nil {
		panic(err)
	}
	Mysql.InitTable()
	defers = append(defers, func() {
		err = mysql.MysqlClose(Mysql.MysqlDb)
		if err != nil {
			panic(err)
		}
	})
	err = Logic.NewConnCh(settings.Config.RabbitMQConfig)
	if err != nil {
		panic(err)
	}
	defers = append(defers, func() {
		Logic.ConnClose()
	})

	Init.InitService()
	service.InitServer(settings.Config.ServiceConfig)

	// 初始化zap日志
	path, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	zap.InitLogger(path)
}
func listen() {
	//grpc服务开启监听
	lis, err = net.Listen(settings.Config.GrpcConfig.NetWork,
		settings.Config.GrpcConfig.Host+
			":"+
			strconv.Itoa(settings.Config.GrpcConfig.Port))
	if err != nil {
		panic(err)
	}

	fmt.Println("friend service start success at", lis.Addr().String())
}
func findService() {
	// 其他服务发现
	client.Init(settings.Config.EtcdConfig.Host, settings.Config.EtcdConfig.Port)
}
func serve() { //服务注册
	ctx := context.Background()
	server = service.NewServer(ctx, settings.Config.EtcdConfig.Host, settings.Config.EtcdConfig.Port)

	err = server.Serve(lis)
	if err != nil {
		fmt.Println(err)
		return
	}
}
