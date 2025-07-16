package Test

import (
	"api/Mq"
	"api/middleware"
	"github.com/gin-gonic/gin"
	"testing"
)

func TestWebsocket(t *testing.T) {
	err := Mq.NewConnCh()
	if err != nil {
		t.Error(err)
	}
	defer Mq.ConnClose()
	r := gin.Default()

	r.Use(middleware.WebSocketMiddleware())
	r.GET("/ws", func(c *gin.Context) {

	})
	r.Run("127.0.0.1:8888")
}
