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
	jwtAuth     jwt.IJwt
	userUsecase usecase.IUserUsecase
}

func MiddlerwareInit(jwtAuth jwt.IJwt, userUsecase usecase.IUserUsecase) IMiddleware {
	return &Middleware{
		jwtAuth:     jwtAuth,
		userUsecase: userUsecase,
	}
}
