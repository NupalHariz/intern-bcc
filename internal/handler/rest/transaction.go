package rest

import (
	"intern-bcc/domain"
	"intern-bcc/pkg/response"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (r *Rest) CreateTransaction(c *gin.Context) {
	mentorIdString := c.Param("mentorId")
	mentorId, _ := strconv.Atoi(mentorIdString)

	var transactionRequest domain.TransactionRequest

	err := c.ShouldBindJSON(&transactionRequest)
	if err != nil {
		response.Failed(c, http.StatusBadRequest, "failed to bind json", err)
	}

	coreApiRes, errorObject := r.usecase.TransactionUsecase.CreateTransaction(c, mentorId, transactionRequest)
	if errorObject != nil {
		errorObject := errorObject.(response.ErrorObject)
		response.Failed(c, errorObject.Code, errorObject.Message, errorObject.Err)
		return
	}

	response.Success(c, "waiting for payment", coreApiRes)
}

func (r *Rest) VerifyTransaction(c *gin.Context) {
	var payLoad map[string]interface{}

	err := c.ShouldBindJSON(&payLoad)
	if err != nil {
		response.Failed(c, http.StatusBadRequest, "failed to bind json", err)
	}

	errorObject := r.usecase.TransactionUsecase.VerifyTransaction(payLoad)
	if errorObject != nil {
		errorObject := errorObject.(response.ErrorObject)
		response.Failed(c, errorObject.Code, errorObject.Message, errorObject.Err)
		return
	}

	response.SuccessWithoutData(c, "payment success")
}
