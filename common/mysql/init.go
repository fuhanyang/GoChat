package mysql

import (
	"fmt"
	"github.com/jinzhu/gorm"
)

func Init(config *MysqlConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=true", config.UserName, config.Password, config.Host, config.Port, config.DbName)
	db, err := gorm.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(100)
	fmt.Println("mysql connect success at:", dsn)
	return db, nil

}
func MysqlClose(db *gorm.DB) error {
	if db == nil {
		return nil
	}
	return db.Close()
}
