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
		response.Failed(c, response.NewError(http.StatusBadRequest, "failed to bind request", err))
	}

	coreApiRes, err := r.usecase.TransactionUsecase.CreateTransaction(c, mentorId, transactionRequest)
	if err != nil {
		response.Failed(c, err)
		return
	}

	response.Success(c, "waiting for payment", coreApiRes)
}

func (r *Rest) VerifyTransaction(c *gin.Context) {
	var payLoad map[string]interface{}

	err := c.ShouldBindJSON(&payLoad)
	if err != nil {
		response.Failed(c, response.NewError(http.StatusBadRequest, "failed to bind request", err))
	}

	err = r.usecase.TransactionUsecase.VerifyTransaction(payLoad)
	if err != nil {
		response.Failed(c, err)
		return
	}

	response.Success(c, "payment success", nil)
}
