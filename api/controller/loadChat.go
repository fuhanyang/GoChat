package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func LoadChat() gin.HandlerFunc {
	return func(c *gin.Context) {
		accountNum := c.Query("account_num")
		if accountNum == "" {
			c.Redirect(http.StatusFound, "/")
			return
		}
		c.HTML(http.StatusOK, "chat.html", gin.H{
			"account_num": accountNum,
		})
	}
}
