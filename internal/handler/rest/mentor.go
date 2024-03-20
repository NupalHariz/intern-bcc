package rest

import (
	"intern-bcc/domain"
	"intern-bcc/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (r *Rest) GetMentor(c *gin.Context) {
	mentorIdString := c.Param("mentorId")
	mentorId, err := uuid.Parse(mentorIdString)
	if err != nil {
		response.Failed(c, response.NewError(http.StatusBadRequest, "failed to parsing string to uuid", err))
		return
	}

	mentorParam := domain.MentorParam{
		Id: mentorId,
	}

	mentor, err := r.usecase.MentorUsecase.GetMentor(mentorParam)
	if err != nil {
		response.Failed(c, err)
		return
	}

	response.Success(c, "success get mentor", mentor)
}

func (r *Rest) GetMentors(c *gin.Context) {
	ctx := c.Request.Context()

	mentors, err := r.usecase.MentorUsecase.GetMentors(ctx)
	if err != nil {
		response.Failed(c, err)
		return
	}

	response.Success(c, "succes get mentor", mentors)
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
	mentorIdString := c.Param("mentorId")
	mentorId, err := uuid.Parse(mentorIdString)
	if err != nil {
		response.Failed(c, response.NewError(http.StatusBadRequest, "failed to parsing string to uuid", err))
		return
	}

	mentorParam := domain.MentorParam{
		Id: mentorId,
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

	response.Success(c, "succes update mentor", nil)
}

func (r *Rest) UploadMentorPicture(c *gin.Context) {
	mentorIdString := c.Param("mentorId")
	mentorId, err := uuid.Parse(mentorIdString)
	if err != nil {
		response.Failed(c, response.NewError(http.StatusBadRequest, "failed to parsing string to uuid", err))
		return
	}

	mentorParam := domain.MentorParam{
		Id: mentorId,
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
