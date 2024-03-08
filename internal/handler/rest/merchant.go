package rest

import (
	"intern-bcc/domain"
	"intern-bcc/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (r *Rest) CreateMerchant(c *gin.Context) {
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

	response.Success(c, "success to create merchant, please verify your merchant", nil)
}

func (r *Rest) SendOtp(c *gin.Context) {
	ctx := c.Request.Context()

	errorObject := r.usecase.MerchantUsecase.SendOtp(c, ctx)
	if errorObject != nil {
		errorObject := errorObject.(response.ErrorObject)
		response.Failed(c, errorObject.Code, errorObject.Message, errorObject.Err)
		return
	}

	response.Success(c, "please check your email for verification", nil)
}

func (r *Rest) VerifyOtp(c *gin.Context) {
	ctx := c.Request.Context()

	var verifyOtp domain.MerchantVerify
	err := c.ShouldBindJSON(&verifyOtp)
	if err != nil {
		response.Failed(c, http.StatusBadRequest, "failed to bind json", err)
		return
	}

	errorObject := r.usecase.MerchantUsecase.VerifyOtp(c, ctx, verifyOtp)
	if errorObject != nil {
		errorObject := errorObject.(response.ErrorObject)
		response.Failed(c, errorObject.Code, errorObject.Message, errorObject.Err)
		return
	}

	response.Success(c, "verification succeed", nil)
}
