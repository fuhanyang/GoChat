package Test

import (
	"github.com/gin-gonic/gin"
	"server/Mq"
	"server/WebSocket"
	"testing"
)

func TestWebsocket(t *testing.T) {
	err := Mq.NewConnCh()
	if err != nil {
		t.Error(err)
	}
	defer Mq.ConnClose()
	r := gin.Default()

	r.Use(WebSocket.WebSocketMiddleware())
	r.GET("/ws", func(c *gin.Context) {

	})
	r.Run("127.0.0.1:8888")
}
