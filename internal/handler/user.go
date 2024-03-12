package handler

import (
	"intern-bcc/domain"
	"intern-bcc/internal/usecase"
	"intern-bcc/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userUsecase usecase.IUserUsecase
}

func NewUserHandler(userUsecase usecase.IUserUsecase) *UserHandler {
	return &UserHandler{userUsecase}
}

func (h *UserHandler) Register(c *gin.Context) {
	var userRequest domain.UserRequest

	err := c.ShouldBindJSON(&userRequest)
	if err != nil {
		response.Failed(c, http.StatusBadRequest, "failed to bind request", err)
		return
	}

	errorObject := h.userUsecase.Register(userRequest)
	if errorObject != nil {
		errorObject := errorObject.(response.ErrorObject)
		response.Failed(c, errorObject.Code, errorObject.Message, errorObject.Err)
		return
	}

	response.Success(c, "success create account", nil)
}

func (h *UserHandler) Login(c *gin.Context) {
	var userLogin domain.UserLogin

	err := c.ShouldBindJSON(&userLogin)
	if err != nil {
		response.Failed(c, http.StatusBadRequest, "failed to bind request", err)
		return
	}

	loginRespone, errorObject := h.userUsecase.Login(userLogin)
	if errorObject != nil {
		errorObject := errorObject.(response.ErrorObject)
		response.Failed(c, errorObject.Code, errorObject.Message, errorObject.Err)
		return
	}

	response.Success(c, "login success", loginRespone)
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
	var userUpdate domain.UserUpdate

	err := c.ShouldBindJSON(&userUpdate)
	if err != nil {
		response.Failed(c, http.StatusBadRequest, "failed to bind request", err)
		return
	}

	updatedUser, errorObject := h.userUsecase.UpdateUser(c, userUpdate)
	if errorObject != nil {
		errorObject := errorObject.(response.ErrorObject)
		response.Failed(c, errorObject.Code, errorObject.Message, errorObject.Err)
		return
	}

	response.Success(c, "success update user", updatedUser)
}

func (h *UserHandler) UploadPhoto(c *gin.Context) {
	profilePicture, err := c.FormFile("profile_picture")
	if err != nil {
		response.Failed(c, http.StatusBadRequest, "failed to bind request", err)
		return
	}

	var userPhoto domain.UploadUserPhoto
	userPhoto.ProfilePicture = profilePicture

	errorObject := h.userUsecase.UploadPhoto(c, profilePicture)
	if errorObject != nil {
		errorObject := errorObject.(response.ErrorObject)
		response.Failed(c, errorObject.Code, errorObject.Message, errorObject.Err)
		return
	}

	response.Success(c, "success updload photo", nil)
}
