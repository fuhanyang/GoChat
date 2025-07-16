package SendMessage

import (
	redis2 "common/redis"
	"github.com/gomodule/redigo/redis"
	"gorm.io/gorm"
	"message/DAO/Redis"
	"message/Models"
	"time"
)

func SendVideo(db *gorm.DB, sender string, receiver string, size int64, format string, content []byte, CreatedAt time.Time) error {
	video := Models.NewVideo(size, format)
	video.Content = content
	video.SenderAccountNum = sender
	video.ReceiverAccountNum = receiver
	Models.WriteMessage(db, video)
	// TODO: Redis
	args := redis.Args{sender + receiver + video.CreatedAt.String()}.AddFlat(video)
	_, err := Redis.RedisDo(redis2.HMSET, args...)
	return err
}
