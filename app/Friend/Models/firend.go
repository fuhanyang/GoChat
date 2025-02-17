package Models

import (
	"Friend/DAO/Mysql"
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	"log"
	"strings"
)

var ErrDuplicateFriend = errors.New("好友关系已存在")

type Friend struct {
	gorm.Model
	User1     string `gorm:"column:user1;index:idx_user1" redis:"user1"`
	User2     string `gorm:"column:user2;index:idx_user2" redis:"user2"`
	User1Name string `gorm:"column:user1_name" redis:"user1_name"`
	User2Name string `gorm:"column:user2_name" redis:"user2_name"`
}

func AddCompositeUniqueIndex(db *gorm.DB, dbName string) error {
	// 检查是否存在 unique_friendship 约束
	var indexExists int
	row := db.Raw("SELECT COUNT(*) FROM information_schema.STATISTICS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? AND INDEX_NAME = ?", dbName, "friends", "unique_friendship").Row()
	if err := row.Scan(&indexExists); err != nil {
		return fmt.Errorf("检查索引失败: %v", err)
	}

	if indexExists > 0 {
		log.Println("唯一约束 unique_friendship 已存在，跳过添加")
		return nil
	}

	// 添加唯一约束
	return db.Exec("ALTER TABLE friends ADD CONSTRAINT unique_friendship UNIQUE (user1, user2);").Error
}
func GetFriend(userAccountNum string, limit int) []Friend {
	Friends := make([]Friend, limit)
	Mysql.MysqlDb.Where("user1 = ?", userAccountNum).
		Or("user2 = ?", userAccountNum).
		Order("created_at desc").
		Limit(limit).
		Find(&Friends)
	return Friends
}
func CreateFriend(user1, user2, user1Name, user2Name string) error {
	err := Mysql.MysqlDb.Create(&Friend{
		User1:     user1,
		User2:     user2,
		User1Name: user1Name,
		User2Name: user2Name,
	}).Error

	if err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") {
			return ErrDuplicateFriend
		}
		return fmt.Errorf("创建好友关系失败: %v", err)
	}

	return nil
}
