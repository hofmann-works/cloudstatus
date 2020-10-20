package main

import (
	"github.com/gin-gonic/gin"
	"github.com/hofmann-works/cloudstatus/config"
	"github.com/hofmann-works/cloudstatus/handlers"
)

func main() {
	var conf = config.New()
	println(conf.PollInterval)

	router := gin.Default()

	v1 := router.Group("v1")
	v1.GET("/status", handlers.Status)

	router.Run(":8080")
}
