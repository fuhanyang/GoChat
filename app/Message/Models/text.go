package Models

import (
	"sync"
)

type Text struct {
	Message
	Size int64 `gorm:"column:size" json:"size" redis:"size"`
}

var textPool = sync.Pool{
	New: func() interface{} {
		return &Text{}
	},
}

func NewText(size int64) *Text {
	var text = textPool.Get().(*Text)
	text.ID = 0
	text.Size = size
	text.Content = make([]byte, 0, size)
	return text
}

func (text *Text) Release() {
	text.ID = 0
	textPool.Put(text)
}
