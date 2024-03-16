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
	CreateTransaction(c *gin.Context, mentorId int, transactionRequest domain.TransactionRequest) (domain.TransactionResponse, any)
	VerifyTransaction(payload map[string]interface{}) any
}

type TransactionUsecase struct {
	transactionRepository repository.ITransactionRepository
	userRepository        repository.IUserRepository
	jwt                   jwt.IJwt
	midTrans              midtrans.IMidTrans
}

func NewTransactionUsecase(transactionRepository repository.ITransactionRepository, userRepository repository.IUserRepository, jwt jwt.IJwt, midTrans midtrans.IMidTrans) ITransactionUsecase {
	return &TransactionUsecase{
		transactionRepository: transactionRepository,
		userRepository:        userRepository,
		jwt:                   jwt,
		midTrans:              midTrans,
	}
}

func (u *TransactionUsecase) CreateTransaction(c *gin.Context, mentorId int, transactionRequest domain.TransactionRequest) (domain.TransactionResponse, any) {
	user, err := u.jwt.GetLoginUser(c)
	if err != nil {
		return domain.TransactionResponse{}, response.ErrorObject{
			Code:    http.StatusNotFound,
			Message: "failed to get user id",
			Err:     err,
		}
	}

	layoutFormat := "2006-01-02 15:04:05"
	nullTime, err := time.Parse(layoutFormat, "1970-01-01 00:00:01")
	if err != nil {
		return domain.TransactionResponse{}, response.ErrorObject{
			Code:    http.StatusInternalServerError,
			Message: "failed parsing null time",
			Err:     err,
		}
	}
	newTransaction := domain.Transactions{
		Id:          uuid.New(),
		UserId:      user.Id,
		MentorId:    mentorId,
		Price:       transactionRequest.Price,
		PaymentType: transactionRequest.PaymentType,
		PayedAt:     nullTime,
	}

	coreApiRes, err := u.midTrans.ChargeTransaction(newTransaction)
	if err != nil {
		return domain.TransactionResponse{}, response.ErrorObject{
			Code:    http.StatusInternalServerError,
			Message: "an error occured when charge transaction",
			Err:     err,
		}
	}

	err = u.transactionRepository.CreateTransaction(&newTransaction)
	if err != nil {
		return domain.TransactionResponse{}, response.ErrorObject{
			Code:    http.StatusInternalServerError,
			Message: "an error occured when create transaction",
			Err:     err,
		}
	}

	transactionResponse := generateTransactionResponse(transactionRequest, coreApiRes)
	
	return transactionResponse, err
}

func generateTransactionResponse(transactionRequest domain.TransactionRequest, coreApiRes *coreapi.ChargeResponse) domain.TransactionResponse{
	transactionResponse := domain.TransactionResponse {
		TransactionId: coreApiRes.TransactionID,
		PaymentType: transactionRequest.PaymentType,
	}

	switch transactionRequest.PaymentType {
	case "gopay" :
		transactionResponse.URL = coreApiRes.Actions[0].URL
	case "mandiri" :
		transactionResponse.BillKey = coreApiRes.BillKey
		transactionResponse.BillerCode = coreApiRes.BillerCode
	default :
	transactionResponse.VaNumber = coreApiRes.VaNumbers[0].VANumber
	}

	return transactionResponse
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
			Code:    http.StatusInternalServerError,
			Message: "an error occured when update transaction",
			Err:     err,
		}
	}

	mentor := domain.HasMentor{
		UserId:   transaction.UserId,
		MentorId: transaction.MentorId,
	}

	err = u.userRepository.CreateHasMentor(&mentor)
	if err != nil {
		return response.ErrorObject{
			Code:    http.StatusInternalServerError,
			Message: "an error occured when create mentor and student relation",
			Err:     err,
		}
	}

	return nil
}
