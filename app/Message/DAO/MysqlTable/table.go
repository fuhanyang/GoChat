package MysqlTable

import (
	"Message/DAO/Mysql"
	"Message/Models"
)

func InitTable() {
	Mysql.MysqlDb.AutoMigrate(&Models.Text{})
	Mysql.MysqlDb.AutoMigrate(&Models.Audio{})
	Mysql.MysqlDb.AutoMigrate(&Models.Video{})
}
