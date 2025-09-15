package main

import (
	"common/mysql"
	"common/redis"
	"common/snowflake"
	"common/viper"
	"common/zap"
	"context"
	"flag"
	"fmt"
	"net"
	"os"
	"rpc/Service/Init"
	"strconv"
	"user/DAO/Mysql"
	"user/DAO/Redis"
	"user/Logic/authentication"
	"user/rpc/client"
	"user/rpc/service"
	"user/settings"
)

var (
	lis       net.Listener
	err       error
	bitmapLen int64 = 1008547758
	hashCount int32 = 5
	defers          = make([]func(), 0)
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
	userInit()
	newBloomFilter(bitmapLen, hashCount)
	listen()
	findService()
	serve()
}
func newBloomFilter(bitmapLen int64, hashCount int32) {
	authentication.NewBloomFilter(bitmapLen, hashCount, Mysql.Db)
}
func userInit() {
	Redis.Pool = redis.Init(settings.Config.RedisConfig)
	Redis.Client = redis.InitRedisLock(settings.Config.RedisConfig)
	err = snowflake.Init(settings.Config.SnowflakeConfig)
	if err != nil {
		panic(err)
	}
	Mysql.Db, err = mysql.Init(settings.Config.MysqlConfig)
	if err != nil {
		panic(err)
	}
	defers = append(defers, func() {
		err = mysql.MysqlClose(Mysql.Db)
		if err != nil {
			panic(err)
		}
	})
	Mysql.InitTable(Mysql.Db)

	Init.InitService()
	service.InitServer(settings.Config.ServiceConfig)

	// 初始化zap日志
	path, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	zap.InitLogger(path)
}
func findService() {
	// 其他服务发现
	client.Init(settings.Config.EtcdConfig.Host, settings.Config.EtcdConfig.Port)
}
func listen() {
	lis, err = net.Listen(settings.Config.GrpcConfig.NetWork,
		settings.Config.GrpcConfig.Host+
			":"+
			strconv.Itoa(settings.Config.GrpcConfig.Port))
	if err != nil {
		panic(err)
	}
	fmt.Println("user service start success at", lis.Addr().String())
}
func serve() {
	// 注册服务
	ctx := context.Background()
	server := service.NewServer(ctx, settings.Config.EtcdConfig.Host, settings.Config.EtcdConfig.Port)

	err = server.Serve(lis)
	if err != nil {
		fmt.Println(err)
		return
	}
}
