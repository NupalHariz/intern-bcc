package handler

import (
	"intern-bcc/domain"
	"intern-bcc/internal/usecase"
	"intern-bcc/pkg/response"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type TransactionHandler struct {
	transactionUsecase usecase.ITransactionUsecase
}

func NewTransactionHandler(transactionUsecase usecase.ITransactionUsecase) *TransactionHandler {
	return &TransactionHandler{transactionUsecase}
}

func (h *TransactionHandler) CreateTransaction(c *gin.Context) {
	mentorIdString := c.Param("mentorId")
	mentorId, _ := strconv.Atoi(mentorIdString)

	var transactionRequest domain.TransactionRequest

	err := c.ShouldBindJSON(&transactionRequest)
	if err != nil {
		response.Failed(c, http.StatusBadRequest, "failed to bind json", err)
	}

	coreApiRes, errorObject := h.transactionUsecase.CreateTransaction(c, mentorId, transactionRequest)
	if errorObject != nil {
		errorObject := errorObject.(response.ErrorObject)
		response.Failed(c, errorObject.Code, errorObject.Message, errorObject.Err)
		return
	}

	response.Success(c, "waiting for payment", coreApiRes)
}

func (h *TransactionHandler) VerifyTransaction(c *gin.Context) {
	var payLoad map[string]interface{}

	err := c.ShouldBindJSON(&payLoad)
	if err != nil {
		response.Failed(c, http.StatusBadRequest, "failed to bind json", err)
	}

	errorObject := h.transactionUsecase.VerifyTransaction(payLoad)
	if errorObject != nil {
		errorObject := errorObject.(response.ErrorObject)
		response.Failed(c, errorObject.Code, errorObject.Message, errorObject.Err)
		return
	}

	response.Success(c, "payment success", nil)
}
