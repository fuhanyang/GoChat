package main

import (
	"common/bloomFilter"
	"common/redis"
	"common/viper"
	"common/zap"
	"context"
	"flag"
	"message/DAO/Mysql"
	"message/DAO/Redis"
	"message/rpc/service"
	settings "message/settings"
	"os"
	"rpc/Service/Init"
	"strconv"

	"fmt"
	"net"
)

var (
	lis          net.Listener
	err          error
	_bloomFilter bloomFilter.BloomFilter
	bitmapLen    int64 = 10000000
	hashCount    int32 = 10
	defers             = make([]func(), 0)
)

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
	messageInit()
	listen()
	serve()
}
func messageInit() {

	Redis.Pool = redis.Init(settings.Config.RedisConfig)

	err = Mysql.Init(settings.Config.MysqlConfig)
	if err != nil {
		panic(err)
	}
	Mysql.InitTable(Mysql.MysqlDb)

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
	lis, err = net.Listen(settings.Config.GrpcConfig.NetWork,
		settings.Config.GrpcConfig.Host+
			":"+
			strconv.Itoa(settings.Config.GrpcConfig.Port))
	if err != nil {
		panic(err)
	}
	fmt.Println("message service start success at", lis.Addr().String())
}
func serve() {
	//服务注册
	ctx := context.Background()
	server := service.NewServer(ctx, settings.Config.EtcdConfig.Host, settings.Config.EtcdConfig.Port)

	err = server.Serve(lis)
	if err != nil {
		fmt.Println(err)
		return
	}
}
