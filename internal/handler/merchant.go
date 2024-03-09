package handler

import (
	"intern-bcc/domain"
	"intern-bcc/internal/usecase"
	"intern-bcc/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type MerchantHandler struct {
	merchantUsecase usecase.IMerchantUsecase
}

func NewMerchantHandler(merchantUsecase usecase.IMerchantUsecase) *MerchantHandler {
	return &MerchantHandler{merchantUsecase}
}

func (h *MerchantHandler) CreateMerchant(c *gin.Context) {
	var merchantRequest domain.MerchantRequest

	err := c.ShouldBindJSON(&merchantRequest)
	if err != nil {
		response.Failed(c, http.StatusBadRequest, "failed to bind json", err)
		return
	}

	errorObject := h.merchantUsecase.CreateMerchant(c, merchantRequest)
	if errorObject != nil {
		errorObject := errorObject.(response.ErrorObject)
		response.Failed(c, errorObject.Code, errorObject.Message, errorObject.Err)
		return
	}

	response.Success(c, "success to create merchant, please verify your merchant", nil)
}

func (h *MerchantHandler) SendOtp(c *gin.Context) {
	ctx := c.Request.Context()

	errorObject := h.merchantUsecase.SendOtp(c, ctx)
	if errorObject != nil {
		errorObject := errorObject.(response.ErrorObject)
		response.Failed(c, errorObject.Code, errorObject.Message, errorObject.Err)
		return
	}

	response.Success(c, "please check your email for verification", nil)
}

func (h *MerchantHandler) VerifyOtp(c *gin.Context) {
	ctx := c.Request.Context()

	var verifyOtp domain.MerchantVerify
	err := c.ShouldBindJSON(&verifyOtp)
	if err != nil {
		response.Failed(c, http.StatusBadRequest, "failed to bind json", err)
		return
	}

	errorObject := h.merchantUsecase.VerifyOtp(c, ctx, verifyOtp)
	if errorObject != nil {
		errorObject := errorObject.(response.ErrorObject)
		response.Failed(c, errorObject.Code, errorObject.Message, errorObject.Err)
		return
	}

	response.Success(c, "verification succeed", nil)
}
