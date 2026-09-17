package httpserver

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func CustomErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err

			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Terjadi Kesalahan pada request",
				"mesasge": err.Error(),
			})
		}
	}
}
