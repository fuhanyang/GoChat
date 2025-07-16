package main

import (
	"common/bloomFilter"
	"common/mysql"
	"common/redis"
	"common/viper"
	"common/zap"
	"context"
	"fmt"
	"friend/DAO/Mysql"
	"friend/DAO/Redis"
	"friend/rpc/client"
	"friend/rpc/service"
	settings "friend/settings"
	"google.golang.org/grpc"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"rpc/Service/Init"
	"strconv"
)

var (
	lis          net.Listener
	err          error
	_bloomFilter bloomFilter.BloomFilter
	bitmapLen    int64 = 100000000
	hashCount    int32 = 10
	server       *grpc.Server
	defers       = make([]func(), 0)
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
	defer func() {
		for _, f := range defers {
			f()
		}
	}()
	friendInit()
	listen()
	findService()
	serve()
}
func friendInit() {
	err := viper.Init(settings.Config)
	if err != nil {
		panic(err)
	}
	fmt.Println("config init success mode:", settings.Config.Mode)

	Mysql.MysqlDb, err = mysql.Init(settings.Config.MysqlConfig)
	if err != nil {
		panic(err)
	}
	defers = append(defers, func() {
		err = mysql.MysqlClose(Mysql.MysqlDb)
		if err != nil {
			panic(err)
		}
	})
	Mysql.InitTable(settings.Config.MysqlConfig)
	Redis.Pool = redis.Init(settings.Config.RedisConfig)

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
