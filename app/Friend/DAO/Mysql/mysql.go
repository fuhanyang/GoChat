package Mysql

import (
	settings "Friend/Settings"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

var MysqlDb *gorm.DB

func Init(config *settings.MysqlConfig) error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=true", config.UserName, config.Password, config.Host, config.Port, config.DbName)
	var err error
	MysqlDb, err = gorm.Open("mysql", dsn)
	if err != nil {
		return err
	}
	MysqlDb.DB().SetMaxIdleConns(10)
	MysqlDb.DB().SetMaxOpenConns(100)
	fmt.Println("mysql connect success")
	return nil

}
func MysqlClose() error {
	return MysqlDb.Close()

}
