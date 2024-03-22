package rest

import (
	"intern-bcc/domain"
	"intern-bcc/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (r *Rest) CreateProvince(c *gin.Context) {
	var provinceRequest domain.Province
	err := c.ShouldBindJSON(&provinceRequest)
	if err != nil {
		response.Failed(c, response.NewError(http.StatusBadRequest, "failed to bind request", err))
		return
	}

	err = r.usecase.ProvinceUsecase.CreateProvince(provinceRequest)
	if err != nil {
		response.Failed(c, err)
		return
	}

	response.Success(c, "success create province", nil)
}
