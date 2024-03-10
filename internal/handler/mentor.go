package handler

import (
	"intern-bcc/domain"
	"intern-bcc/internal/usecase"
	"intern-bcc/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type MentorHandler struct {
	mentorUsecase usecase.IMentorUsecase
}

func NewMentorHandler(mentorUsecase usecase.IMentorUsecase) *MentorHandler {
	return &MentorHandler{mentorUsecase}
}

func (h *MentorHandler) CreateMentor(c *gin.Context) {
	var mentorRequest domain.MentorRequest
	err := c.ShouldBindJSON(&mentorRequest)
	if err != nil {
		response.Failed(c, http.StatusBadRequest, "failed to bind json", err)
	}

	errorObject := h.mentorUsecase.CreateMentor(c, mentorRequest)
	if errorObject != nil {
		errorObject := errorObject.(response.ErrorObject)
		response.Failed(c, errorObject.Code, errorObject.Message, errorObject.Err)
		return
	}

	response.Success(c, "succes create mentor", nil)
}
