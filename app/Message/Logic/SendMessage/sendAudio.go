package SendMessage

import (
	redis2 "common/redis"
	"github.com/gomodule/redigo/redis"
	"gorm.io/gorm"
	"message/DAO/Redis"
	"message/Models"
	"time"
)

func SendAudio(db *gorm.DB, sender string, receiver string, size int64, format string, content []byte, CreatedAt time.Time) error {
	audio := Models.NewAudio(size, format)
	audio.Content = content
	audio.SenderAccountNum = sender
	audio.ReceiverAccountNum = receiver
	Models.WriteMessage(db, audio)
	// TODO: Redis
	args := redis.Args{sender + receiver + audio.CreatedAt.String()}.AddFlat(audio)
	_, err := Redis.RedisDo(redis2.HMSET, args...)
	return err
}
