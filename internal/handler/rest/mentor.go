package rest

import (
	"intern-bcc/domain"
	"intern-bcc/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

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
		response.Failed(c, response.NewError(http.StatusBadRequest, "failed to parsing mentor id", err))
		return
	}

	var mentorUpdate domain.MentorUpdate
	err = c.ShouldBindJSON(&mentorUpdate)
	if err != nil {
		response.Failed(c, response.NewError(http.StatusBadRequest, "failed to bind request", err))
		return
	}

	mentor, err := r.usecase.MentorUsecase.UpdateMentor(mentorId, mentorUpdate)
	if err != nil {
		response.Failed(c, err)
		return
	}

	response.Success(c, "succes create mentor", mentor)
}

func (r *Rest) UploadMentorPicture(c *gin.Context) {
	mentorIdString := c.Param("mentorId")
	mentorId, err := uuid.Parse(mentorIdString)
	if err != nil {
		response.Failed(c, response.NewError(http.StatusBadRequest, "failed to parsing mentor id", err))
		return
	}

	mentorPicture, err := c.FormFile("mentor_picture")
	if err != nil {
		response.Failed(c, response.NewError(http.StatusBadRequest, "failed to bind request", err))
		return
	}

	mentor, err := r.usecase.MentorUsecase.UploadMentorPhoto(mentorId, mentorPicture)
	if err != nil {
		response.Failed(c, err)
		return
	}

	response.Success(c, "succes upload mentor picture", mentor)
}
