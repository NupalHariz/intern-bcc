package rest

import (
	"intern-bcc/domain"
	"intern-bcc/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (r *Rest) GetMentor(c *gin.Context) {
	var mentorParam domain.MentorParam
	err := c.ShouldBind(&mentorParam)
	if err != nil {
		response.Failed(c, response.NewError(http.StatusBadRequest, "failed to bind request", err))
		return
	}

	mentor, err := r.usecase.MentorUsecase.GetMentor(mentorParam)
	if err != nil {
		response.Failed(c, err)
		return
	}

	response.Success(c, "success get mentor", mentor)
}

func (r *Rest) GetMentors(c *gin.Context) {
	mentors, err := r.usecase.MentorUsecase.GetMentors()
	if err != nil {
		response.Failed(c, err)
		return
	}

	response.Success(c, "succes create mentor", mentors)
}

func (r *Rest) CreateMentor(c *gin.Context) {
	var mentorRequest domain.MentorRequest
	err := c.ShouldBindJSON(&mentorRequest)
	if err != nil {
		response.Failed(c, response.NewError(http.StatusBadRequest, "failed to bind request", err))
		return
	}

	err = r.usecase.MentorUsecase.CreateMentor(mentorRequest)
	if err != nil {
		response.Failed(c, err)
		return
	}

	response.Success(c, "succes create mentor", nil)
}

func (r *Rest) UpdateMentor(c *gin.Context) {
	var mentorParam domain.MentorParam
	err := c.ShouldBind(&mentorParam)
	if err != nil {
		response.Failed(c, response.NewError(http.StatusBadRequest, "failed to bind request", err))
		return
	}

	var mentorUpdate domain.MentorUpdate
	err = c.ShouldBindJSON(&mentorUpdate)
	if err != nil {
		response.Failed(c, response.NewError(http.StatusBadRequest, "failed to bind request", err))
		return
	}

	err = r.usecase.MentorUsecase.UpdateMentor(mentorParam, mentorUpdate)
	if err != nil {
		response.Failed(c, err)
		return
	}

	response.Success(c, "succes create mentor", nil)
}

func (r *Rest) UploadMentorPicture(c *gin.Context) {
	var mentorParam domain.MentorParam
	err := c.ShouldBind(&mentorParam)
	if err != nil {
		response.Failed(c, response.NewError(http.StatusBadRequest, "failed to bind request", err))
		return
	}

	mentorPicture, err := c.FormFile("mentor_picture")
	if err != nil {
		response.Failed(c, response.NewError(http.StatusBadRequest, "failed to bind request", err))
		return
	}

	err = r.usecase.MentorUsecase.UploadMentorPhoto(mentorParam, mentorPicture)
	if err != nil {
		response.Failed(c, err)
		return
	}

	response.Success(c, "succes upload mentor picture", nil)
}
