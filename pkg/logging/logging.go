package logging

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type ILogging interface {
	Info(c *gin.Context, statusCode int, executionTime float64)
	Error(c *gin.Context, statusCode int, executionTime float64, err error)
	InfoLn(message string)
	ErrorLn(err error)
	WarnLn(err error)
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

func (l *Logging) Info(c *gin.Context, statusCode int, executionTime float64) {
	l.logrus.WithFields(logrus.Fields{
		"success":     true,
		"status_code": statusCode,
		"duration":    executionTime,
	}).Info()
}

func (l *Logging) Error(c *gin.Context, statusCode int, executionTime float64, err error) {
	l.logrus.WithFields(logrus.Fields{
		"success":     false,
		"status_code": statusCode,
		"duration":    executionTime,
		"error":       err,
	}).Error()
}

func (l *Logging) ErrorLn(err error) {
	l.logrus.Errorln(err)
}

func (l *Logging) InfoLn(message string) {
	l.logrus.Println(message)
}

func(l *Logging) WarnLn(err error) {
	l.logrus.Warningln(err)
}
