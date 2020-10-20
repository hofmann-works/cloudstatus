package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Status(c *gin.Context) {
	//check cloud status here
	c.String(http.StatusOK, "Cloud status OK")
}
