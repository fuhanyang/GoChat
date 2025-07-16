package Mysql

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"user/Models"
)

var Db *gorm.DB

func InitTable(db *gorm.DB) {
	if db == nil {
		panic("db is nil")
	}
	db.AutoMigrate(&Models.User{})
}
