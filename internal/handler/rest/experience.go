package rest

import (
	"intern-bcc/domain"
	"intern-bcc/pkg/response"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (r *Rest) AddExperience(c *gin.Context) {
	mentorIdString := c.Param("mentorId")
	mentorId, _ := strconv.Atoi(mentorIdString)

	var experienceRequest domain.ExperienceRequest

	err := c.ShouldBindJSON(&experienceRequest)
	if err != nil {
		response.Failed(c, response.NewError(http.StatusBadRequest, "failed to bind request", err))
		return
	}

	err = r.usecase.ExperienceUsecase.AddExperience(experienceRequest, mentorId)
	if err != nil {
		response.Failed(c, err)
		return
	}

	response.Success(c, "succes create experience", nil)
}
