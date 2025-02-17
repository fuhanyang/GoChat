package SendMessage

import (
	"Message/DAO/Redis"
	"Message/Models"
	"github.com/gomodule/redigo/redis"
	"time"
)

func SendAudio(sender string, receiver string, size int64, format string, content []byte, CreatedAt time.Time) error {
	audio := Models.NewAudio(size, format)
	audio.Content = content
	audio.SenderAccountNum = sender
	audio.ReceiverAccountNum = receiver
	Models.WriteMessage(audio)
	// TODO: Redis
	args := redis.Args{sender + receiver + audio.CreatedAt.String()}.AddFlat(audio)
	_, err := Redis.RedisDo(Redis.HMSET, args...)
	return err
}
