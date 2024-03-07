package rest

import (
	"intern-bcc/domain"
	"intern-bcc/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (r *Rest) CreateMerchant(c *gin.Context){
	var merchantRequest domain.MerchantRequest

	err := c.ShouldBindJSON(&merchantRequest)
	if err != nil {
		response.Failed(c, http.StatusBadRequest, "failed to bind json", err)
		return
	}

	errorObject := r.usecase.MerchantUsecase.CreateMerchant(c, merchantRequest)
	if errorObject != nil {
		errorObject := errorObject.(response.ErrorObject)
		response.Failed(c, errorObject.Code, errorObject.Message, errorObject.Err)
		return
	}

	response.Succes(c, "success to create merchant, please verify your merchant", nil)
}