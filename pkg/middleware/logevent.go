package middleware

import (
	"encoding/json"
	"time"

	"github.com/gin-gonic/gin"
)

func (m *Middleware) LogEvent(c *gin.Context) {
	timeStart := time.Now()
	accessLogMap := make(map[string]string)
	accessLogMap["endpoint"] = c.Request.URL.Path
	accessLogMap["request_http_method"] = c.Request.Method
	accessLogMap["request_client_ip"] = c.ClientIP()
	logData, err := json.Marshal(accessLogMap)
	if err != nil {
		m.logging.ErrorLn(err)
	}
	accessLogJson := string(logData)
	m.logging.InfoLn(accessLogJson)

	c.Next()

	statusCode := c.Writer.Status()
	duration := time.Since(timeStart).Seconds()
	errorResponse, _ := c.Get("error")
	if errorResponse != nil {
		m.logging.Error(c, statusCode, duration, errorResponse.(error))
	} else {
		m.logging.Info(c, statusCode, duration)
	}
}
