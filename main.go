package main

import (
	"github.com/hofmann-works/cloudstatus/db"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hofmann-works/cloudstatus/checks"
	"github.com/hofmann-works/cloudstatus/config"
	"github.com/hofmann-works/cloudstatus/handlers"
)

func main() {
	var conf = config.New()

	database, err := db.Init(conf.PGHost, conf.PGPort, conf.PGDatabase, conf.PGUser, conf.PGPassword)
	if err != nil {
		log.Fatal("Could not set up database: ", err)
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
		log.Println("Polling Cloud Status...")
		checks.AzureStatus(database)
		checks.GitHubStatus(database)
	}
}
