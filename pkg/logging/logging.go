package logging

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type ILogging interface {
	Info(c *gin.Context, message string, time time.Duration, responseBody interface{})
	Error(c *gin.Context, message string, err error)
	InfoLn(message string)
	ErrorLn(err error)
	WarnLn(message string)
}

type Logging struct {
	logrus *logrus.Logger
}

func LoggingInit() ILogging {
	logger := logrus.New()
	customFormatter := new(logrus.JSONFormatter)

	logger.SetFormatter(customFormatter)
	logger.SetReportCaller(true)
	logger.SetLevel(logrus.DebugLevel)

	return &Logging{logrus: logger}
}

func (l *Logging) Info(c *gin.Context, message string, time time.Duration, responseBody interface{}) {
	l.logrus.WithFields(logrus.Fields{
		"path":          c.Request.RequestURI,
		"response_body": responseBody,
		"duration":      time.Seconds(),
	}).Info(message)
}

func (l *Logging) Error(c *gin.Context, message string, err error) {
	l.logrus.WithFields(logrus.Fields{
		"path":  c.Request.RequestURI,
		"error": err.Error(),
	}).Error(message)
}

func (l *Logging) WarnLn(message string) {
	l.logrus.Warnln(message)
}

func (l *Logging) ErrorLn(err error) {
	l.logrus.Errorln(err)
}

func (l *Logging) InfoLn(message string) {
	l.logrus.Println(message)
}
