package middleware

import (
	"intern-bcc/internal/usecase"
	"intern-bcc/pkg/jwt"

	"github.com/gin-gonic/gin"
)

type IMiddleware interface {
	Authentication(c *gin.Context)
	OnlyAdmin(c *gin.Context)
}

type Middleware struct {
	jwtAuth jwt.IJwt
	usecase *usecase.Usecase
}

func MiddlerwareInit(jwtAuth jwt.IJwt, usecase *usecase.Usecase) IMiddleware {
	return &Middleware{
		jwtAuth:     jwtAuth,
		usecase: usecase,
	}
}
