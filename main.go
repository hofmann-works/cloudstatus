package main

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hofmann-works/cloudstatus/checks"
	"github.com/hofmann-works/cloudstatus/config"
	"github.com/hofmann-works/cloudstatus/handlers"
)

func main() {
	var conf = config.New()

	fmt.Println("Polling cloud status every ", conf.PollInterval, " seconds.")
	go pollStatus(conf.PollInterval)

	router := gin.Default()

	v1 := router.Group("v1")
	v1.GET("/status", handlers.Status)

	router.Run(":8080")
}

func pollStatus(pollInterval int) {
	for range time.Tick(time.Second * time.Duration(pollInterval)) {
		checks.AzureStatus()
		checks.GitHubStatus()
	}
}
