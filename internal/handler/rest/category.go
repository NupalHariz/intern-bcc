package rest

import (
	"intern-bcc/domain"
	"intern-bcc/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (r *Rest) CreateCategory(c *gin.Context) {
	var categoryRequest domain.Categories

	err := c.ShouldBindJSON(&categoryRequest)
	if err != nil {
		response.Failed(c, response.NewError(http.StatusBadRequest, "failed to bind request", err))
		return
	}

	err = r.usecase.CategoryUsecase.CreateCategory(categoryRequest)
	if err != nil {
		response.Failed(c, err)
		return
	}

	response.Success(c, "success create category", nil)
}
