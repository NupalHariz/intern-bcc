package middleware

import (
	"errors"
	"intern-bcc/pkg/jwt"
	"intern-bcc/pkg/response"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func Authentication(c *gin.Context) {
	bearer := c.GetHeader("Authorization")
	if bearer == "" {
		response.Failed(c, http.StatusUnauthorized, "token is empty", errors.New("need token"))
		c.Abort()
		return
	}

	tokenString := strings.Split(bearer, " ")[1]
	userId, err := jwt.ValidateToken(tokenString)
	if err != nil {
		response.Failed(c, http.StatusUnauthorized, "failed to validate token", err)
		c.Abort()
		return
	}

	// user, errorObject := m.usecase.UserUsecase.GetUser(domain.UserParam{Id: userId})
	// if errorObject != nil {
	// 	errorObject := errorObject.(response.ErrorObject)
	// 	response.Failed(c, errorObject.Code, errorObject.Message, errorObject.Err)
	// }

	c.Set("userId", userId)
	c.Next()
}
