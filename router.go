package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func register(r *gin.Engine) {
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
}
