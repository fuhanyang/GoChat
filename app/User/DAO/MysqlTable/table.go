package MysqlTable

import (
	"User/DAO/Mysql"
	"User/Models"
)

func InitTable() {
	Mysql.MysqlDb.AutoMigrate(&Models.User{})

}
