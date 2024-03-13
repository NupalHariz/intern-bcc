package rest

import (
	"intern-bcc/domain"
	"intern-bcc/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (r *Rest) CreateInformation(c *gin.Context) {
	var informationRequest domain.InformationRequest

	err := c.ShouldBindJSON(&informationRequest)
	if err != nil {
		response.Failed(c, http.StatusBadRequest, "failed to bind json", err)
	}

	errorObject := r.usecase.InformationUsecase.CreateInformation(informationRequest)
	if errorObject != nil {
		errorObject := errorObject.(response.ErrorObject)
		response.Failed(c, errorObject.Code, errorObject.Message, errorObject.Err)
		return
	}

	response.Success(c, "success create information", nil)
}
