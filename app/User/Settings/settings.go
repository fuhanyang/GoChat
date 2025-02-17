package Settings

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var Config = &AppConfig{}

type AppConfig struct {
	App `mapstructure:"app"`
}

type App struct {
	Name             string `mapstructure:"name"`
	Mode             string `mapstructure:"mode"`
	Version          string `mapstructure:"version"`
	Port             int    `mapstructure:"port"`
	*LogConfig       `mapstructure:"log"`
	*MysqlConfig     `mapstructure:"mysql"`
	*RedisConfig     `mapstructure:"redis"`
	*SnowflakeConfig `mapstructure:"snowflake"`
	*GrpcConfig      `mapstructure:"grpc"`
	*EtcdConfig      `mapstructure:"etcd"`
	*ServiceConfig   `mapstructure:"service"`
}
type ServiceConfig struct {
	Name     string `mapstructure:"name"`
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Protocol string `mapstructure:"protocol"`
}
type EtcdConfig struct {
	Host string `mapstructure:"host"`
	Port string `mapstructure:"port"`
}
type GrpcConfig struct {
	Host    string `mapstructure:"host"`
	Port    int    `mapstructure:"port"`
	NetWork string `mapstructure:"network"`
}
type SnowflakeConfig struct {
	StartTime string `mapstructure:"start_time"`
}
type LogConfig struct {
	Level      string `mapstructure:"level"`
	FileName   string `mapstructure:"filename"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxAge     int    `mapstructure:"max_age"`
	MaxBackups int    `mapstructure:"max_backups"`
}
type MysqlConfig struct {
	Host         string `mapstructure:"host"`
	Port         int    `mapstructure:"port"`
	UserName     string `mapstructure:"user_name"`
	Password     string `mapstructure:"password"`
	DbName       string `mapstructure:"db_name"`
	MaxOpenConns int    `mapstructure:"max_open_conns"`
	MaxIdleConns int    `mapstructure:"max_idle_conns"`
}

type RedisConfig struct {
	Host      string `mapstructure:"host"`
	Port      int    `mapstructure:"port"`
	Password  string `mapstructure:"password"`
	Db        int    `mapstructure:"db"`
	Network   string `mapstructure:"network"`
	MaxIdle   int    `mapstructure:"max_idle"`
	MaxActive int    `mapstructure:"max_active"`
}

func Init() (error error) {
	viper.SetConfigName("config")      // 配置文件名称(无扩展名)
	viper.SetConfigType("yaml")        // 如果配置文件的名称中没有扩展名，则需要配置此项
	viper.AddConfigPath("./Settings/") // 查找配置文件所在的路径
	err := viper.ReadInConfig()        // 查找并读取配置文件
	if err != nil {                    // 处理读取配置文件的错误
		fmt.Println("Fatal error config file: error:", err)
		return err
	}
	//这里我封装了一层结构体
	if err = viper.Unmarshal(&Config); err != nil {
		fmt.Println("Fatal error config file: error:", err)
	}
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Config file changed:", e.Name)
		if err = viper.Unmarshal(&Config); err != nil {
			fmt.Println("Fatal error config file: error:", err)
		}
	})

	fmt.Println("Config file loaded successfully")
	return nil
}
