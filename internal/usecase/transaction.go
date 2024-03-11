package usecase

import (
	"errors"
	"intern-bcc/domain"
	"intern-bcc/internal/repository"
	"intern-bcc/pkg/jwt"
	"intern-bcc/pkg/midtrans"
	"intern-bcc/pkg/response"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/midtrans/midtrans-go/coreapi"
)

type ITransactionUsecase interface {
	CreateTransaction(c *gin.Context, mentorId int, transactionRequest domain.TransactionRequest) (*coreapi.ChargeResponse, any)
	VerifyTransaction(payload map[string]interface{}) any
}

type TransactionUsecase struct {
	transactionRepository repository.ITransactionRepository
	jwt                   jwt.IJwt
	midTrans              midtrans.IMidTrans
}

func NewTransactionRepository(transactionRepository repository.ITransactionRepository, jwt jwt.IJwt, midTrans midtrans.IMidTrans) ITransactionUsecase {
	return &TransactionUsecase{
		transactionRepository: transactionRepository,
		jwt:                   jwt,
		midTrans:              midTrans,
	}
}

func (u *TransactionUsecase) CreateTransaction(c *gin.Context, mentorId int, transactionRequest domain.TransactionRequest) (*coreapi.ChargeResponse, any) {
	user, err := u.jwt.GetLoginUser(c)
	if err != nil {
		return nil, response.ErrorObject{
			Code:    http.StatusNotFound,
			Message: "failed to get user id",
			Err:     err,
		}
	}

	newTransaction := domain.Transactions{
		Id:          uuid.New(),
		UserId:      user.Id,
		MentorId:    mentorId,
		Price:       transactionRequest.Price,
		PaymentType: transactionRequest.PaymentType,
	}

	coreApiRes, err := u.midTrans.ChargeTransaction(newTransaction)
	if err != nil {
		return coreApiRes, response.ErrorObject{
			Code:    http.StatusInternalServerError,
			Message: "an error occured when charge transaction",
			Err:     err,
		}
	}

	err = u.transactionRepository.CreateTransaction(&newTransaction)
	if err != nil {
		return coreApiRes, response.ErrorObject{
			Code:    http.StatusInternalServerError,
			Message: "an error occured when create transaction",
			Err:     err,
		}
	}

	return coreApiRes, err
}

func (u *TransactionUsecase) VerifyTransaction(payload map[string]interface{}) any {
	transactionIdString, exist := payload["order_id"].(string)
	if !exist {
		return response.ErrorObject{
			Code:    http.StatusNotFound,
			Message: "transaction not found",
			Err:     errors.New("can't find the transaction"),
		}
	}

	success, err := u.midTrans.VerifyPayment(transactionIdString)
	if !success {
		return response.ErrorObject{
			Code:    http.StatusBadRequest,
			Message: "transaction failed",
			Err:     err,
		}
	}

	transactionId, err := uuid.Parse(transactionIdString)
	if err != nil {
		return response.ErrorObject{
			Code:    http.StatusInternalServerError,
			Message: "failed to parsing string to uuid",
			Err:     err,
		}
	}

	var transaction domain.Transactions
	transaction.Id = transactionId
	err = u.transactionRepository.GetTransaction(&transaction)
	if err != nil {
		return response.ErrorObject{
			Code:    http.StatusNotFound,
			Message: "transaction not found",
			Err:     err,
		}
	}

	layout := "2006-01-02 15:04:05"
	paymentTime, err := time.Parse(layout, payload["transaction_time"].(string))
	if err != nil {
		return response.ErrorObject{
			Code:    http.StatusInternalServerError,
			Message: "failed to parsing string to time",
			Err:     err,
		}
	}
	transaction.IsPayed = true
	transaction.PayedAt = paymentTime

	err = u.transactionRepository.UpdateTransaction(&transaction)
	if err != nil {
		return response.ErrorObject{
			Code: http.StatusInternalServerError,
			Message: "an error occured when update transaction",
			Err: err,
		}
	}

	return nil
}
