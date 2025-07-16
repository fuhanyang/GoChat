package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func LoadStart() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", gin.H{})
	}
}
