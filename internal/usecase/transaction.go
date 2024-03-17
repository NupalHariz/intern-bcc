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
	CreateTransaction(c *gin.Context, mentorId uuid.UUID, transactionRequest domain.TransactionRequest) (domain.TransactionResponse, error)
	VerifyTransaction(payload map[string]interface{}) error
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

func (u *TransactionUsecase) CreateTransaction(c *gin.Context, mentorId uuid.UUID, transactionRequest domain.TransactionRequest) (domain.TransactionResponse, error) {
	user, err := u.jwt.GetLoginUser(c)
	if err != nil {
		return domain.TransactionResponse{}, response.NewError(http.StatusNotFound, "an error occured when get login user", err)
	}

	//Null value for Payed_At
	layoutFormat := "2006-01-02 15:04:05"
	nullTime, err := time.Parse(layoutFormat, "1970-01-01 00:00:01")
	if err != nil {
		return domain.TransactionResponse{}, response.NewError(http.StatusInternalServerError, "failed to parsing null time", err)
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
		return domain.TransactionResponse{}, response.NewError(http.StatusInternalServerError, "an error occured when charge transaction", err)
	}

	err = u.transactionRepository.CreateTransaction(&newTransaction)
	if err != nil {
		return domain.TransactionResponse{}, response.NewError(http.StatusInternalServerError, "an error occured when create transaction", err)
	}

	transactionResponse := generateTransactionResponse(transactionRequest, coreApiRes)

	return transactionResponse, err
}

func generateTransactionResponse(transactionRequest domain.TransactionRequest, coreApiRes *coreapi.ChargeResponse) domain.TransactionResponse {
	transactionResponse := domain.TransactionResponse{
		TransactionId: coreApiRes.TransactionID,
		PaymentType:   transactionRequest.PaymentType,
	}

	if transactionRequest.PaymentType == "gopay" {
		transactionResponse.URL = coreApiRes.Actions[0].URL
	} else if transactionRequest.PaymentType == "mandiri" {
		transactionResponse.BillKey = coreApiRes.BillKey
		transactionResponse.BillerCode = coreApiRes.BillerCode
	} else if transactionRequest.PaymentType == "bca" || transactionRequest.PaymentType == "bri" || transactionRequest.PaymentType == "bni" {
		transactionResponse.VaNumber = coreApiRes.VaNumbers[0].VANumber
	}

	return transactionResponse
}

func (u *TransactionUsecase) VerifyTransaction(payload map[string]interface{}) error {
	transactionIdString, exist := payload["order_id"].(string)
	if !exist {
		return response.NewError(http.StatusNotFound, "transaction not found", errors.New("can't find the transaction"))
	}

	success, err := u.midTrans.VerifyPayment(transactionIdString)
	if !success {
		return response.NewError(http.StatusBadRequest, "transaction failed", err)
	}

	transactionId, err := uuid.Parse(transactionIdString)
	if err != nil {
		return response.NewError(http.StatusInternalServerError, "failed to parsing string to uuid", err)
	}

	var transaction domain.Transactions
	transaction.Id = transactionId
	err = u.transactionRepository.GetTransaction(&transaction)
	if err != nil {
		return response.NewError(http.StatusNotFound, "transaction not found", err)
	}

	layout := "2006-01-02 15:04:05"
	paymentTime, err := time.Parse(layout, payload["transaction_time"].(string))
	if err != nil {
		return response.NewError(http.StatusInternalServerError, "failed to parsing string to time", err)
	}
	transaction.IsPayed = true
	transaction.PayedAt = paymentTime

	err = u.transactionRepository.UpdateTransaction(&transaction)
	if err != nil {
		return response.NewError(http.StatusInternalServerError, "an error occured when update transaction", err)
	}

	mentor := domain.HasMentor{
		UserId:   transaction.UserId,
		MentorId: transaction.MentorId,
	}

	err = u.userRepository.CreateHasMentor(&mentor)
	if err != nil {
		return response.NewError(http.StatusInternalServerError, "an error occured when create mentor and student relation", err)
	}

	return nil
}
