package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func parseIDParam(c *gin.Context) (uint, bool) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id invalido"})
		return 0, false
	}
	return uint(id), true
}
