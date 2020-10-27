package handlers

import (
	"github.com/hofmann-works/cloudstatus/db"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Status(c *gin.Context) {
	database := c.MustGet("databaseConn").(db.Database)
	response, _ := database.GetLatestChecks()

	c.JSON(http.StatusOK, response)
}
