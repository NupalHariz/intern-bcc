package middleware

import (
	"errors"
	"intern-bcc/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (m *Middleware) OnlyAdmin(c *gin.Context) {
	user, err := m.jwtAuth.GetLoginUser(c)
	if err != nil {
		response.Failed(c, response.NewError(http.StatusNotFound, "failed to get account", err))
		c.Abort()
	}
	if !user.IsAdmin {
		response.Failed(c, response.NewError(http.StatusUnauthorized, "access denied", errors.New("only admin")))
	}

	c.Next()
}
