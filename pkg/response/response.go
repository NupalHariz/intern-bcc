package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ErrorObject struct {
	Code    int
	Message string
	Err     error
}

func Failed(c *gin.Context, code int, message string, err error) {
	c.JSON(code, gin.H{
		"message": message,
		"error":   err,
	})
}

func Succes(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"message": message,
		"data": data,
	})
}
