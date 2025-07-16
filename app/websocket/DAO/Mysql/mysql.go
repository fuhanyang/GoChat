package Mysql

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"websocket/Models"
)

var MysqlDb *gorm.DB

func InitTable() {
	MysqlDb.AutoMigrate(&Models.WebSocketData{})
}
