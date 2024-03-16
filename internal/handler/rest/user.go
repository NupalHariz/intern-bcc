package rest

import (
	"intern-bcc/domain"
	"intern-bcc/pkg/response"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

	response.SuccessWithoutData(c, "success create account")
}

func (r *Rest) Login(c *gin.Context) {
	var userLogin domain.UserLogin

	err := c.ShouldBindJSON(&userLogin)
	if err != nil {
		response.Failed(c, http.StatusBadRequest, "failed to bind request", err)
		return
	}

	loginRespone, errorObject := r.usecase.UserUsecase.Login(userLogin)
	if errorObject != nil {
		errorObject := errorObject.(response.ErrorObject)
		response.Failed(c, errorObject.Code, errorObject.Message, errorObject.Err)
		return
	}

	response.Success(c, "login success", loginRespone)
}

func (r *Rest) UpdateUser(c *gin.Context) {
	userIdString := c.Param("userId")
	userId, err := uuid.Parse(userIdString)
	if err != nil {
		response.Failed(c, http.StatusBadRequest, "failed to parsing user id", err)
	}
	var userUpdate domain.UserUpdate

	err = c.ShouldBindJSON(&userUpdate)
	if err != nil {
		response.Failed(c, http.StatusBadRequest, "failed to bind request", err)
		return
	}

	updatedUser, errorObject := r.usecase.UserUsecase.UpdateUser(c, userId, userUpdate)
	if errorObject != nil {
		errorObject := errorObject.(response.ErrorObject)
		response.Failed(c, errorObject.Code, errorObject.Message, errorObject.Err)
		return
	}

	response.Success(c, "success update user", updatedUser)
}

func (r *Rest) UploadUserPhoto(c *gin.Context) {
	userIdString := c.Param("userId")
	userId, err := uuid.Parse(userIdString)
	if err != nil {
		response.Failed(c, http.StatusBadRequest, "failed to parsing user id", err)
	}

	profilePicture, err := c.FormFile("profile_picture")
	if err != nil {
		response.Failed(c, http.StatusBadRequest, "failed to bind request", err)
		return
	}

	user, errorObject := r.usecase.UserUsecase.UploadUserPhoto(c, userId, profilePicture)
	if errorObject != nil {
		errorObject := errorObject.(response.ErrorObject)
		response.Failed(c, errorObject.Code, errorObject.Message, errorObject.Err)
		return
	}

	response.Success(c, "success updload photo", user)
}

func (r *Rest) PasswordRecovery(c *gin.Context) {
	ctx := c.Request.Context()

	var userParam domain.UserParam
	err := c.ShouldBindJSON(&userParam)
	if err != nil {
		response.Failed(c, http.StatusBadRequest, "failed to bind request", err)
		return
	}

	errorObject := r.usecase.UserUsecase.PasswordRecovery(userParam, ctx)
	if errorObject != nil {
		errorObject := errorObject.(response.ErrorObject)
		response.Failed(c, errorObject.Code, errorObject.Message, errorObject.Err)
		return
	}

	response.SuccessWithoutData(c, "please check your email")
}

func (r *Rest) ChangePassword(c *gin.Context) {
	ctx := c.Request.Context()
	name := c.Param("name")
	verPass := c.Param("verPass")

	var passwordRequest domain.PasswordUpdate
	err := c.ShouldBindJSON(&passwordRequest)
	if err != nil {
		response.Failed(c, http.StatusBadRequest, "failed to bind request", err)
		return
	}

	errorObject := r.usecase.UserUsecase.ChangePassword(ctx, name, verPass, passwordRequest)
	if errorObject != nil {
		errorObject := errorObject.(response.ErrorObject)
		response.Failed(c, errorObject.Code, errorObject.Message, errorObject.Err)
		return
	}

	response.SuccessWithoutData(c, "success change password")
}

func (r *Rest) LikeProduct(c *gin.Context) {
	productIdString := c.Param("productId")
	productId, err := strconv.Atoi(productIdString)
	if err != nil {
		response.Failed(c, http.StatusBadRequest, "failed to parsing product id", err)
	}

	errorObject := r.usecase.UserUsecase.LikeProduct(c, productId)
	if errorObject != nil {
		errorObject := errorObject.(response.ErrorObject)
		response.Failed(c, errorObject.Code, errorObject.Message, errorObject.Err)
		return
	}

	response.SuccessWithoutData(c, "success like product")
}

func (r *Rest) DeleteLikeProduct(c *gin.Context) {
	productIdString := c.Param("productId")
	productId, err := uuid.Parse(productIdString)
	if err != nil {
		response.Failed(c, http.StatusBadRequest, "failed to parsing product id", err)
	}

	errorObject := r.usecase.UserUsecase.DeleteLikeProduct(c, productId)
	if errorObject != nil {
		errorObject := errorObject.(response.ErrorObject)
		response.Failed(c, errorObject.Code, errorObject.Message, errorObject.Err)
		return
	}

	response.SuccessWithoutData(c, "success delete liked product")
}
