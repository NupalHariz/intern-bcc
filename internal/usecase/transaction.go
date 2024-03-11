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

type ITransactionUsecase interface {
	CreateTransaction(c *gin.Context, mentorId int, transactionRequest domain.TransactionRequest) (*coreapi.ChargeResponse, any)
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
