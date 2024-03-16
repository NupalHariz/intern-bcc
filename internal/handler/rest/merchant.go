package rest

import (
	"intern-bcc/domain"
	"intern-bcc/pkg/response"
	"net/http"
	"strconv"

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

	response.SuccessWithoutData(c, "success to create merchant, please verify your merchant")
}

func (r *Rest) SendOtp(c *gin.Context) {
	ctx := c.Request.Context()

	errorObject := r.usecase.MerchantUsecase.SendOtp(c, ctx)
	if errorObject != nil {
		errorObject := errorObject.(response.ErrorObject)
		response.Failed(c, errorObject.Code, errorObject.Message, errorObject.Err)
		return
	}

	response.SuccessWithoutData(c, "please check your email for verification")
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

	response.SuccessWithoutData(c, "verification succeed")
}

func (r *Rest) UpdateMerchant(c *gin.Context) {
	merchantIdString := c.Param("merchantId")
	merchantId, err := strconv.Atoi(merchantIdString)
	if err != nil {
		response.Failed(c, http.StatusBadRequest, "failed to parse merchant id", err)
		return
	}

	var updateMerchant domain.UpdateMerchant

	err = c.ShouldBindJSON(&updateMerchant)
	if err != nil {
		response.Failed(c, http.StatusBadRequest, "failed to bind json", err)
		return
	}

	merchant, errorObject := r.usecase.MerchantUsecase.UpdateMerchant(c, merchantId, updateMerchant)
	if errorObject != nil {
		errorObject := errorObject.(response.ErrorObject)
		response.Failed(c, errorObject.Code, errorObject.Message, errorObject.Err)
		return
	}

	response.Success(c, "success update merchant", merchant)
}

func (r *Rest) UploadMerchantPhoto(c *gin.Context) {
	merchantIdString := c.Param("merchantId")
	merchantId, err := strconv.Atoi(merchantIdString)
	if err != nil {
		response.Failed(c, http.StatusBadRequest, "failed to parse merchant id", err)
		return
	}

	merchantPhoto, err := c.FormFile("merchant_photo")
	if err != nil {
		response.Failed(c, http.StatusBadRequest, "failed to bind request", err)
		return
	}

	merchant, errorObject := r.usecase.MerchantUsecase.UploadMerchantPhoto(c, merchantId, merchantPhoto)
	if errorObject != nil {
		errorObject := errorObject.(response.ErrorObject)
		response.Failed(c, errorObject.Code, errorObject.Message, errorObject.Err)
		return
	}

	response.Success(c, "success upload merchant photo", merchant)
}
