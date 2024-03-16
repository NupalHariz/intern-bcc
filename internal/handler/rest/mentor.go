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
		response.Failed(c, http.StatusBadRequest, "failed to bind json", err)
		return
	}

	errorObject := r.usecase.MentorUsecase.CreateMentor(mentorRequest)
	if errorObject != nil {
		errorObject := errorObject.(response.ErrorObject)
		response.Failed(c, errorObject.Code, errorObject.Message, errorObject.Err)
		return
	}

	response.SuccessWithoutData(c, "succes create mentor")
}

func (r *Rest) UpdateMentor(c *gin.Context) {
	mentorIdString := c.Param("mentorId")
	mentorId, err := uuid.Parse(mentorIdString)
	if err != nil {
		response.Failed(c, http.StatusBadRequest, "failed to parsing mentor id", err)
		return
	}

	var mentorUpdate domain.MentorUpdate
	err = c.ShouldBindJSON(&mentorUpdate)
	if err != nil {
		response.Failed(c, http.StatusBadRequest, "failed to bind json", err)
		return
	}

	mentor, errorObject := r.usecase.MentorUsecase.UpdateMentor(mentorId, mentorUpdate)
	if errorObject != nil {
		errorObject := errorObject.(response.ErrorObject)
		response.Failed(c, errorObject.Code, errorObject.Message, errorObject.Err)
		return
	}

	response.Success(c, "succes create mentor", mentor)
}

func (r *Rest) UploadMentorPicture(c *gin.Context) {
	mentorIdString := c.Param("mentorId")
	mentorId, err := uuid.Parse(mentorIdString)
	if err != nil {
		response.Failed(c, http.StatusBadRequest, "failed to parsing mentor id", err)
		return
	}

	mentorPicture, err := c.FormFile("mentor_picture")
	if err != nil {
		response.Failed(c, http.StatusBadRequest, "failed to bind request", err)
		return
	}

	mentor, errorObject := r.usecase.MentorUsecase.UploadMentorPhoto(mentorId, mentorPicture)
	if errorObject != nil {
		errorObject := errorObject.(response.ErrorObject)
		response.Failed(c, errorObject.Code, errorObject.Message, errorObject.Err)
		return
	}

	response.Success(c, "succes upload mentor picture", mentor)
}
