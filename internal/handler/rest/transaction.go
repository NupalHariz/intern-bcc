package rest

import (
	"intern-bcc/domain"
	"intern-bcc/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (r *Rest) CreateTransaction(c *gin.Context) {
	mentorIdString := c.Param("mentorId")
	mentorId, _ := uuid.Parse(mentorIdString)

	var transactionRequest domain.TransactionRequest

	err := c.ShouldBindJSON(&transactionRequest)
	if err != nil {

		response.Failed(c, response.NewError(http.StatusBadRequest, "failed to bind request", err))
		return
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
		return
	}

	err = r.usecase.TransactionUsecase.VerifyTransaction(payLoad)
	if err != nil {
		response.Failed(c, err)
		return
	}

	response.Success(c, "payment success", nil)
}
