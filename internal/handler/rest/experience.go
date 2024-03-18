package rest

import (
	"intern-bcc/domain"
	"intern-bcc/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (r *Rest) AddExperience(c *gin.Context) {
	mentorIdString := c.Param("mentorId")
	mentorId, _ := uuid.Parse(mentorIdString)
	mentorParam := domain.MentorParam{
		MentorId: mentorId,
	}

	var experienceRequest domain.ExperienceRequest

	err := c.ShouldBindJSON(&experienceRequest)
	if err != nil {
		response.Failed(c, response.NewError(http.StatusBadRequest, "failed to bind request", err))
		return
	}

	err = r.usecase.ExperienceUsecase.AddExperience(experienceRequest, mentorParam)
	if err != nil {
		response.Failed(c, err)
		return
	}

	response.Success(c, "succes create experience", nil)
}
