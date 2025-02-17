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
	Name            string `mapstructure:"name"`
	Mode            string `mapstructure:"mode"`
	Version         string `mapstructure:"version"`
	Port            int    `mapstructure:"port"`
	*RabbitMQConfig `mapstructure:"rabbitmq"`
	*EtcdConfig     `mapstructure:"etcd"`
}

type EtcdConfig struct {
	Host string `mapstructure:"host"`
	Port string `mapstructure:"port"`
}
type RabbitMQConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
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
