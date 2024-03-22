package rest

import (
	"intern-bcc/domain"
	"intern-bcc/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (r *Rest) Register(c *gin.Context) {
	var userRequest domain.UserRequest

	err := c.ShouldBindJSON(&userRequest)
	if err != nil {

		response.Failed(c, response.NewError(http.StatusBadRequest, "failed to bind request", err))
		return
	}

	err = r.usecase.UserUsecase.Register(userRequest)
	if err != nil {

		response.Failed(c, err)
		return
	}

	response.Success(c, "success create account", nil)
}

func (r *Rest) Login(c *gin.Context) {
	var userLogin domain.UserLogin

	err := c.ShouldBindJSON(&userLogin)
	if err != nil {

		response.Failed(c, response.NewError(http.StatusBadRequest, "failed to bind request", err))
		return
	}

	loginRespone, err := r.usecase.UserUsecase.Login(userLogin)
	if err != nil {

		response.Failed(c, err)
		return
	}

	response.Success(c, "login success", loginRespone)
}

func (r *Rest) GetUser(c *gin.Context) {
	userIdString := c.Param("userId")
	userId, err := uuid.Parse(userIdString)
	if err != nil {
		response.Failed(c, response.NewError(http.StatusBadRequest, "failed to parsing product id", err))
		return
	}
	
	userParam := domain.UserParam{
		Id: userId,
	}

	user, err := r.usecase.UserUsecase.GetUser(userParam)
	if err != nil {

		response.Failed(c, err)
		return
	}

	response.Success(c, "success get profile user", user)
}

func (r *Rest) GetOwnProducts(c *gin.Context) {
	products, err := r.usecase.UserUsecase.GetOwnProducts(c)
	if err != nil {

		response.Failed(c, err)
		return
	}

	response.Success(c, "success", products)
}

func (r *Rest) GetLikeProduct(c *gin.Context) {
	products, err := r.usecase.UserUsecase.GetLikeProducts(c)
	if err != nil {

		response.Failed(c, err)
		return
	}

	response.Success(c, "success get like product", products)
}

func (r *Rest) GetOwnMentors(c *gin.Context) {
	mentors, err := r.usecase.UserUsecase.GetOwnMentors(c)
	if err != nil {
		response.Failed(c, err)
		return
	}

	response.Success(c, "success get own mentor", mentors)
}

func (r *Rest) UpdateUser(c *gin.Context) {
	userIdString := c.Param("userId")
	userId, err := uuid.Parse(userIdString)
	if err != nil {
		response.Failed(c, response.NewError(http.StatusBadRequest, "failed to parsing user id", err))
		return
	}
	var userUpdate domain.UserUpdate

	err = c.ShouldBindJSON(&userUpdate)
	if err != nil {
		response.Failed(c, response.NewError(http.StatusBadRequest, "failed to bind request", err))
		return
	}

	updatedUser, err := r.usecase.UserUsecase.UpdateUser(c, userId, userUpdate)
	if err != nil {
		response.Failed(c, err)
		return
	}

	response.Success(c, "success update user", updatedUser)
}

func (r *Rest) UploadUserPhoto(c *gin.Context) {
	userIdString := c.Param("userId")
	userId, err := uuid.Parse(userIdString)
	if err != nil {
		response.Failed(c, response.NewError(http.StatusBadRequest, "failed to parsing user id", err))
		return
	}

	profilePicture, err := c.FormFile("profile_picture")
	if err != nil {
		response.Failed(c, response.NewError(http.StatusBadRequest, "failed to bind request", err))
		return
	}

	user, err := r.usecase.UserUsecase.UploadUserPhoto(c, userId, profilePicture)
	if err != nil {
		response.Failed(c, err)
		return
	}

	response.Success(c, "success updload photo", user)
}

func (r *Rest) PasswordRecovery(c *gin.Context) {
	ctx := c.Request.Context()

	var userParam domain.UserParam
	err := c.ShouldBindJSON(&userParam)
	if err != nil {
		response.Failed(c, response.NewError(http.StatusBadRequest, "failed to bind request", err))
		return
	}

	err = r.usecase.UserUsecase.PasswordRecovery(userParam, ctx)
	if err != nil {
		response.Failed(c, err)
		return
	}

	response.Success(c, "please check your email", nil)
}

func (r *Rest) ChangePassword(c *gin.Context) {
	ctx := c.Request.Context()
	name := c.Param("name")
	verPass := c.Param("verPass")

	var passwordRequest domain.PasswordUpdate
	err := c.ShouldBindJSON(&passwordRequest)
	if err != nil {
		response.Failed(c, response.NewError(http.StatusBadRequest, "failed to bind request", err))
		return
	}

	err = r.usecase.UserUsecase.ChangePassword(ctx, name, verPass, passwordRequest)
	if err != nil {
		response.Failed(c, err)
		return
	}

	response.Success(c, "success change password", nil)
}

func (r *Rest) LikeProduct(c *gin.Context) {
	productIdString := c.Param("productId")
	productId, err := uuid.Parse(productIdString)
	if err != nil {
		response.Failed(c, response.NewError(http.StatusBadRequest, "failed to parsing product id", err))
		return
	}

	err = r.usecase.UserUsecase.LikeProduct(c, productId)
	if err != nil {
		response.Failed(c, err)
		return
	}

	response.Success(c, "success like product", nil)
}

func (r *Rest) DeleteLikeProduct(c *gin.Context) {
	productIdString := c.Param("productId")
	productId, err := uuid.Parse(productIdString)
	if err != nil {
		response.Failed(c, response.NewError(http.StatusBadRequest, "failed to parsing product id", err))
		return
	}

	err = r.usecase.UserUsecase.DeleteLikeProduct(c, productId)
	if err != nil {
		response.Failed(c, err)
		return
	}

	response.Success(c, "success delete liked product", nil)
}
