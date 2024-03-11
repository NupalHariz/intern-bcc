package usecase

import (
	"intern-bcc/domain"
	"intern-bcc/internal/repository"
	"intern-bcc/pkg/jwt"
	"intern-bcc/pkg/midtrans"
	"intern-bcc/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/midtrans/midtrans-go/coreapi"
)

type ITransactionUsecase interface{
	CreateTransaction(c *gin.Context, mentorId int, transactionRequest domain.TransactionRequest) (*coreapi.ChargeResponse, any)
}

type TransactionUsecase struct {
	transactionRepository repository.ITransactionRepository
	userRepository        repository.IUserRepository
}

func NewTransactionRepository(transactionRepository repository.ITransactionRepository, userRepository repository.IUserRepository) ITransactionUsecase {
	return &TransactionUsecase{
		transactionRepository: transactionRepository,
		userRepository:        userRepository,
	}
}

func (u *TransactionUsecase) CreateTransaction(c *gin.Context, mentorId int, transactionRequest domain.TransactionRequest) (*coreapi.ChargeResponse, any) {
	userId, err := jwt.GetLoginUserId(c)
	if err != nil {
		return nil, response.ErrorObject{
			Code:    http.StatusNotFound,
			Message: "failed to get user id",
			Err:     err,
		}
	}
	var user domain.Users
	err = u.userRepository.GetUser(&user, domain.UserParam{Id: userId})
	if err != nil {
		return nil, response.ErrorObject{
			Code:    http.StatusNotFound,
			Message: "failed to get user",
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

	coreApiRes, err := midtrans.ChargeTransaction(newTransaction)
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
