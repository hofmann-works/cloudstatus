package main

import (
	"fmt"
	"github.com/hofmann-works/cloudstatus/db"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hofmann-works/cloudstatus/checks"
	"github.com/hofmann-works/cloudstatus/config"
	"github.com/hofmann-works/cloudstatus/handlers"
)

func main() {
	var conf = config.New()

	database, err := db.Init("localhost", 5432, "cloudstatus", "cloudstatus", "mypw")
	if err != nil {
		fmt.Println("Could not set up database: %v", err)
	}
	defer database.Conn.Close()

	go pollStatus(conf.PollInterval, database)

	router := gin.Default()
	router.Use(func(c *gin.Context) {
		c.Set("databaseConn", database)
		c.Next()
	})

	v1 := router.Group("v1")
	v1.GET("/status", handlers.Status)

	router.Run(":8080")
}

func pollStatus(pollInterval int, database db.Database) {
	for range time.Tick(time.Second * time.Duration(pollInterval)) {
		checks.AzureStatus(database)
		checks.GitHubStatus(database)
	}
}
