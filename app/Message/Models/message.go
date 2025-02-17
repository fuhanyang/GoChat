package Models

import (
	"Message/DAO/Mysql"
	"fmt"
	"time"
)

type MessageInf interface {
	GetContent() []byte
}
type Message struct {
	SenderAccountNum   string `gorm:"column:sender;index:account_num_idx" redis:"sender"`
	ReceiverAccountNum string `gorm:"column:receiver;index:account_num_idx" redis:"receiver"`
	Content            []byte `gorm:"column:content" redis:"content"`
	HasRead            bool   `gorm:"column:has_read;default:false" redis:"has_read"`
	ID                 uint   `gorm:"primary_key"`
	CreatedAt          time.Time
	UpdatedAt          time.Time
	DeletedAt          *time.Time `sql:"index"`
}

func (M *Message) GetContent() []byte {
	return M.Content
}
func WriteMessage(message MessageInf) {
	Mysql.MysqlDb.Create(message)
}

// GetMessages 查找对应用户特定数量的消息
func GetMessages[M MessageInf](senderAccountNum string, ReceiverAccountNum string, limit int) []M {
	Messages := make([]M, limit)
	Mysql.MysqlDb.Where("sender = ? AND receiver = ?", senderAccountNum, ReceiverAccountNum).
		Or("sender = ? AND receiver = ?", ReceiverAccountNum, senderAccountNum).
		Order("created_at desc").
		Limit(limit).
		Find(&Messages)
	//
	fmt.Println("GetMessages from ", senderAccountNum, "and", ReceiverAccountNum, Messages)
	return Messages
}
