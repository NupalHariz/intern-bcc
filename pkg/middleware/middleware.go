package middleware

import (
	"intern-bcc/internal/usecase"
	"intern-bcc/pkg/jwt"
	"intern-bcc/pkg/logging"

	"github.com/gin-gonic/gin"
)

type IMiddleware interface {
	Authentication(c *gin.Context)
	OnlyAdmin(c *gin.Context)
	LogEvent(c *gin.Context)
}

type Middleware struct {
	jwtAuth jwt.IJwt
	usecase *usecase.Usecase
	logging logging.ILogging
}

func MiddlerwareInit(jwtAuth jwt.IJwt, usecase *usecase.Usecase, logging logging.ILogging) IMiddleware {
	return &Middleware{
		jwtAuth: jwtAuth,
		usecase: usecase,
		logging: logging,
	}
}
