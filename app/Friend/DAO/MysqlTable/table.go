package MysqlTable

import (
	"Friend/DAO/Mysql"
	"Friend/Models"
	settings "Friend/Settings"
)

func InitTable(config *settings.MysqlConfig) {
	Mysql.MysqlDb.AutoMigrate(&Models.Friend{})
	err := Models.AddCompositeUniqueIndex(Mysql.MysqlDb, config.DbName)
	if err != nil {
		panic(err)
	}
}
