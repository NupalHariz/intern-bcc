package rest

import (
	"intern-bcc/domain"
	"intern-bcc/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (r *Rest) Register(c *gin.Context) {
	var userRequest domain.UserRequest

	err := c.ShouldBindJSON(&userRequest)
	if err != nil {
		response.Failed(c, http.StatusBadRequest, "failed to bind request", err)
		return
	}

	errorObject := r.usecase.UserUsecase.Register(userRequest)
	if errorObject != nil {
		errorObject := errorObject.(response.ErrorObject)
		response.Failed(c, errorObject.Code, errorObject.Message, errorObject.Err)
		return
	}

	response.Success(c, "success create account", nil)
}

func (r *Rest) Login(c *gin.Context) {
	var userLogin domain.UserLogin

	err := c.ShouldBindJSON(&userLogin)
	if err != nil {
		response.Failed(c, http.StatusBadRequest, "failed to bind request", err)
		return
	}

	loginRespone, errorObject := r.usecase.UserUsecase.Login(userLogin)
	if errorObject != nil{
		errorObject := errorObject.(response.ErrorObject)
		response.Failed(c, errorObject.Code, errorObject.Message, errorObject.Err)
		return
	}

	response.Success(c, "login success", loginRespone)
}
