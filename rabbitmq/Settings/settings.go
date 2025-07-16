package settings

import (
	"common/etcd"
)

var Config = &AppConfig{}

type AppConfig struct {
	App `mapstructure:"app"`
}

type App struct {
	Name             string `mapstructure:"name"`
	Mode             string `mapstructure:"mode"`
	Version          string `mapstructure:"version"`
	*RabbitMQConfig  `mapstructure:"rabbitmq"`
	*etcd.EtcdConfig `mapstructure:"etcd"`
}

type RabbitMQConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}
