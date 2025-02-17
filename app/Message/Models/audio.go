package Models

import (
	"sync"
)

type Audio struct {
	Message
	Size   int64  `gorm:"column: size" json:"size" redis:"size'"`
	Format string `gorm:" column: format" json:"format" redis:"format'"`
}

var audioPool = sync.Pool{
	New: func() interface{} {
		return &Audio{}
	},
}

func NewAudio(size int64, format string) *Audio {
	var audio = audioPool.Get().(*Audio)
	audio.ID = 0
	audio.Size = size
	audio.Format = format
	audio.Content = make([]byte, size)
	return audio
}
func (audio *Audio) Release() {
	audio.ID = 0
	audioPool.Put(audio)
}
