package SendMessage

import (
	"Message/DAO/Redis"
	"Message/Models"
	"github.com/gomodule/redigo/redis"
	"time"
)

func SendVideo(sender string, receiver string, size int64, format string, content []byte, CreatedAt time.Time) error {
	video := Models.NewVideo(size, format)
	video.Content = content
	video.SenderAccountNum = sender
	video.ReceiverAccountNum = receiver
	Models.WriteMessage(video)
	// TODO: Redis
	args := redis.Args{sender + receiver + video.CreatedAt.String()}.AddFlat(video)
	_, err := Redis.RedisDo(Redis.HMSET, args...)
	return err
}
