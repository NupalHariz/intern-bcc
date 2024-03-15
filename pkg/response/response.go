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
		"status": "error",
		"message": message,
		"error":   err.Error(),
	})
}

func Success(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"message": message,
		"data":    data,
	})
}


func SuccessWithoutData(c *gin.Context, message string) {
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"message": message,
	})
}