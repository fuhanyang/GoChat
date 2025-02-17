package Models

import (
	"sync"
)

type Video struct {
	Message
	Size   int64  `gorm:"column:size" json:"size" redis:"size"`
	Format string `gorm:"column:format" json:"format" redis:"format"`
}

var videoPool = sync.Pool{
	New: func() interface{} {
		return &Video{}
	},
}

func NewVideo(size int64, format string) *Video {
	var video = videoPool.Get().(*Video)
	video.ID = 0
	video.Size = size
	video.Format = format
	video.Content = make([]byte, size)
	return video
}

func (video *Video) Release() {
	video.ID = 0
	videoPool.Put(video)
}
