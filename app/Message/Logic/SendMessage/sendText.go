package SendMessage

import (
	"Message/DAO/Redis"
	"Message/Models"
	"github.com/gomodule/redigo/redis"
	"time"
)

func SendText(sender string, receiver string, size int64, content []byte, CreatedAt time.Time) error {
	text := Models.NewText(size)
	text.Content = append(text.Content, content...)
	text.SenderAccountNum = sender
	text.ReceiverAccountNum = receiver
	text.CreatedAt = CreatedAt
	Models.WriteMessage(text)
	// TODO: Redis
	args := redis.Args{sender + receiver + text.CreatedAt.String()}.AddFlat(text)
	_, err := Redis.RedisDo(Redis.HMSET, args...)
	return err
}
