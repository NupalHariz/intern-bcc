package rest

import (
	"intern-bcc/domain"
	"intern-bcc/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (r *Rest) CreateUniversity(c *gin.Context) {
	var universityRequest domain.UniversityRequest
	err := c.ShouldBindJSON(&universityRequest)
	if err != nil {
		response.Failed(c, response.NewError(http.StatusBadRequest, "failed to bind request", err))
	}

	err = r.usecase.UniversityUsecase.CreateUniversity(universityRequest)
	if err != nil {
		response.Failed(c, err)
		return
	}

	response.Success(c, "success create university", nil)
}
