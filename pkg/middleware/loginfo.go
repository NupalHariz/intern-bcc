package middleware

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
)

func (m *Middleware) LogEvent(c *gin.Context) {
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
}
