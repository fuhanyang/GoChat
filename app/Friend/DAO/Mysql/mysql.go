package Mysql

import (
	"common/mysql"
	Models2 "friend/Models"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

var MysqlDb *gorm.DB

func InitTable(config *mysql.MysqlConfig) {
	MysqlDb.AutoMigrate(&Models2.Friend{})
	err := Models2.AddCompositeUniqueIndex(MysqlDb, config.DbName)
	if err != nil {
		panic(err)
	}
}
