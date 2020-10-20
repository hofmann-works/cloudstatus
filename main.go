package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hofmann-works/cloudstatus/config"
)

func main() {
	var conf = config.New()
	println(conf.PollInterval)

	router := gin.Default()

	v1 := router.Group("v1")
	v1.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	router.Run(":8080")
}
