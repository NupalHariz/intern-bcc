package response

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ErrorObject struct {
	Code    int
	Message string
	Err     error
}

func (eo *ErrorObject) Error() string {
	return fmt.Sprintf("HttpStatus: %v\nMessage: %v\nerror: %v", eo.Code, eo.Message, eo.Err)
}

func NewError(code int, message string, err error) error {
	return &ErrorObject{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

func Failed(c *gin.Context, err error) {
	c.Set("error", err)
	errorObject := err.(*ErrorObject)
	c.JSON(errorObject.Code, gin.H{
		"status":  "error",
		"message": errorObject.Message,
		"error":   errorObject.Err.Error(),
	})
}

func Success(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": message,
		"data":    data,
	})
}
