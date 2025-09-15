package viper

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"os"
)

func Init(Config any, mode string) (error error) {
	switch mode {
	case "debug":
		return InitDebugger(Config)
	case "local":
		return InitLocal(Config)
	case "production":
		return InitProduction(Config)
	default:
		return fmt.Errorf("Invalid mode")
	}
}

func InitLocal(Config any) (error error) {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	path := fmt.Sprintf("%s%s", wd, "/settings/")
	viper.SetConfigName("config") // 配置文件名称(无扩展名)
	viper.SetConfigType("yaml")   // 如果配置文件的名称中没有扩展名，则需要配置此项
	viper.AddConfigPath(path)     //"./viper/") // 查找配置文件所在的路径
	err = viper.ReadInConfig()    // 查找并读取配置文件
	if err != nil {               // 处理读取配置文件的错误
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

func InitProduction(Config any) (error error) {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	path := fmt.Sprintf("%s%s", wd, "/settings/")
	viper.SetConfigName("config") // 配置文件名称(无扩展名)
	viper.SetConfigType("yaml")   // 如果配置文件的名称中没有扩展名，则需要配置此项
	viper.AddConfigPath(path)     //"./viper/") // 查找配置文件所在的路径
	err = viper.ReadInConfig()    // 查找并读取配置文件
	if err != nil {               // 处理读取配置文件的错误
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

func InitDebugger(Config any) (error error) {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	path := fmt.Sprintf("%s%s", wd, "/settings/debug/")
	viper.SetConfigName("config") // 配置文件名称(无扩展名)
	viper.SetConfigType("yaml")   // 如果配置文件的名称中没有扩展名，则需要配置此项
	viper.AddConfigPath(path)     //"./viper/") // 查找配置文件所在的路径
	err = viper.ReadInConfig()    // 查找并读取配置文件
	if err != nil {               // 处理读取配置文件的错误
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
