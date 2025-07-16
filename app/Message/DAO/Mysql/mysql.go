package Mysql

import (
	mysql2 "common/mysql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/sharding"
	"hash/crc32"
	"message/Models"
)

var MysqlDb *gorm.DB

func Init(config *mysql2.MysqlConfig) error {
	var (
		dsn string
		err error
	)

	dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=true", config.UserName, config.Password, config.Host, config.Port, config.DbName)

	MysqlDb, err = gorm.Open(mysql.New(mysql.Config{
		DSN: dsn,
	}))
	if err != nil {
		return err
	}

	fmt.Println("mysql connect success")
	return nil

}

func InitTable(db *gorm.DB) {
	var cfg sharding.Config
	// 定义分片配置
	cfg = sharding.Config{
		ShardingKey:    "sender", // 分片字段名
		NumberOfShards: 4,
		ShardingAlgorithm: func(val interface{}) (suffix string, err error) {
			if _, ok := val.(string); !ok {
				// 初始化的时候会按NumberOfShards遍历创建一次，用hash算法会不均匀，所以特殊处理
				suffix = fmt.Sprintf("_%02d", val.(int))
				return suffix, nil
			}
			sender := val.(string) // 确保类型与模型中 sender 字段类型 一致
			hash := crc32.ChecksumIEEE([]byte(sender))
			return fmt.Sprintf("_%02d", hash%4), nil // 表后缀规则（如 _01）
		},
		PrimaryKeyGenerator: sharding.PKSnowflake, // 主键生成器（推荐使用 snowflake）
	}
	db.Use(sharding.Register(cfg, &Models.Text{}, &Models.Audio{}, &Models.Video{}))

	db.AutoMigrate(&Models.Text{})
	db.AutoMigrate(&Models.Audio{})
	db.AutoMigrate(&Models.Video{})
}
