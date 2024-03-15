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
		response.Failed(c, http.StatusBadRequest, "failed to bind json", err)
		return
	}

	errorObject := r.usecase.ExperienceUsecase.AddExperience(experienceRequest, mentorId)
	if errorObject != nil {
		errorObject := errorObject.(response.ErrorObject)
		response.Failed(c, errorObject.Code, errorObject.Message, errorObject.Err)
		return
	}

	response.SuccessWithoutData(c, "succes create experience")
}
