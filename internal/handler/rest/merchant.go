package rest

import (
	"intern-bcc/domain"
	"intern-bcc/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (r *Rest) GetMerchant(c *gin.Context) {
	merchant, err := r.usecase.MerchantUsecase.GetMerchant(c)
	if err != nil{
		response.Failed(c, err)
		return
	}

	response.Success(c, "success get merchant", merchant)
}

func (r *Rest) CreateMerchant(c *gin.Context) {
	var merchantRequest domain.MerchantRequest

	err := c.ShouldBindJSON(&merchantRequest)
	if err != nil {
		response.Failed(c, response.NewError(http.StatusBadRequest, "failed to bind request", err))
		return
	}

	err = r.usecase.MerchantUsecase.CreateMerchant(c, merchantRequest)
	if err != nil {
		response.Failed(c, err)
		return
	}

	response.Success(c, "success to create merchant, please verify your merchant", nil)
}

func (r *Rest) SendOtp(c *gin.Context) {
	ctx := c.Request.Context()

	err := r.usecase.MerchantUsecase.SendOtp(c, ctx)
	if err != nil {
		response.Failed(c, err)
		return
	}

	response.Success(c, "please check your email for verification", nil)
}

func (r *Rest) VerifyOtp(c *gin.Context) {
	ctx := c.Request.Context()

	var verifyOtp domain.MerchantVerify
	err := c.ShouldBindJSON(&verifyOtp)
	if err != nil {
		response.Failed(c, response.NewError(http.StatusBadRequest, "failed to bind request", err))
		return
	}

	err = r.usecase.MerchantUsecase.VerifyOtp(c, ctx, verifyOtp)
	if err != nil {
		response.Failed(c, err)
		return
	}

	response.Success(c, "verification succeed", nil)
}

func (r *Rest) UpdateMerchant(c *gin.Context) {
	merchantIdString := c.Param("merchantId")
	merchantId, err := uuid.Parse(merchantIdString)
	if err != nil {
		response.Failed(c, response.NewError(http.StatusBadRequest, "failed to parse merchant id", err))
		return
	}

	var updateMerchant domain.UpdateMerchant

	err = c.ShouldBindJSON(&updateMerchant)
	if err != nil {
		response.Failed(c, response.NewError(http.StatusBadRequest, "failed to bind request", err))
		return
	}

	merchant, err := r.usecase.MerchantUsecase.UpdateMerchant(c, merchantId, updateMerchant)
	if err != nil {

		response.Failed(c, err)
		return
	}

	response.Success(c, "success update merchant", merchant)
}

func (r *Rest) UploadMerchantPhoto(c *gin.Context) {
	merchantIdString := c.Param("merchantId")
	merchantId, err := uuid.Parse(merchantIdString)
	if err != nil {
		response.Failed(c, response.NewError(http.StatusBadRequest, "failed to parse merchant id", err))
		return
	}

	merchantPhoto, err := c.FormFile("merchant_photo")
	if err != nil {
		response.Failed(c, response.NewError(http.StatusBadRequest, "failed to bind request", err))
		return
	}

	merchant, err := r.usecase.MerchantUsecase.UploadMerchantPhoto(c, merchantId, merchantPhoto)
	if err != nil {
		response.Failed(c, err)
		return
	}

	response.Success(c, "success upload merchant photo", merchant)
}
