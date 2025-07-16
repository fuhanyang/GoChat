package settings

import (
	"common/etcd"
	"common/mysql"
	"common/redis"
)

var Config = &AppConfig{}

type AppConfig struct {
	App `mapstructure:"app"`
}

type App struct {
	Name               string `mapstructure:"name"`
	Mode               string `mapstructure:"mode"`
	Host               string `mapstructure:"host"`
	Port               int    `mapstructure:"port"`
	*LogConfig         `mapstructure:"log"`
	*mysql.MysqlConfig `mapstructure:"mysql"`
	*redis.RedisConfig `mapstructure:"redis"`
	*GrpcConfig        `mapstructure:"grpc"`
	*etcd.EtcdConfig   `mapstructure:"etcd"`
	*ServiceConfig     `mapstructure:"service"`
	*RabbitMQConfig    `mapstructure:"rabbitmq"`
}
type RabbitMQConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}
type ServiceConfig struct {
	Name     string `mapstructure:"name"`
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Protocol string `mapstructure:"protocol"`
}

type GrpcConfig struct {
	Host    string `mapstructure:"host"`
	Port    int    `mapstructure:"port"`
	NetWork string `mapstructure:"network"`
}

type LogConfig struct {
	Level      string `mapstructure:"level"`
	FileName   string `mapstructure:"filename"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxAge     int    `mapstructure:"max_age"`
	MaxBackups int    `mapstructure:"max_backups"`
}
