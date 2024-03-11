package main

import (
	"intern-bcc/internal/handler"
	"intern-bcc/internal/repository"
	"intern-bcc/internal/usecase"
	"intern-bcc/pkg/gomail"
	"intern-bcc/pkg/infrastucture"
	"intern-bcc/pkg/infrastucture/cache"
	"intern-bcc/pkg/infrastucture/database"
	"intern-bcc/pkg/jwt"
	"intern-bcc/pkg/middleware"
	"intern-bcc/pkg/midtrans"
	"intern-bcc/rest"

	"github.com/gin-gonic/gin"
)

func main() {
	infrastucture.LoadEnv()
	cache.ConnectToRedis()
	database.ConnectToDB()
	database.Migrate()

	//Repository
	merchantRepository := repository.NewMerchantRepository(database.DB)
	merchantRedis := repository.NewMerchantRedis(cache.RDB)
	userRepository := repository.NewUserRepository(database.DB)
	mentorRepository := repository.NewMentorRepository(database.DB)
	transactionRepository := repository.NewTransactionRepository(database.DB)

	//pkg
	jwt := jwt.JwtInit()
	goMail := gomail.GoMailInit()
	midTrans := midtrans.MidTransInit()

	//Usecase
	merchantUsecase := usecase.NewMerchantUsecase(merchantRepository, merchantRedis, jwt, goMail)
	userUsecase := usecase.NewUserUsecase(userRepository, jwt)
	mentorUsecase := usecase.NewMentorUsecase(mentorRepository, jwt)
	transactionUsecase := usecase.NewTransactionRepository(transactionRepository, jwt, midTrans)

	//Middleware
	middleware := middleware.MiddlerwareInit(jwt, userUsecase)

	//Handler
	merchantHandler := handler.NewMerchantHandler(merchantUsecase)
	userHandler := handler.NewUserHandler(userUsecase)
	mentorHandler := handler.NewMentorHandler(mentorUsecase)
	transactionHandler := handler.NewTransactionHandler(transactionUsecase)

	rest := rest.NewRest(gin.New(), userHandler, merchantHandler, mentorHandler, transactionHandler, middleware)

	rest.MerchantEndpoint()
	rest.UserEndpoint()
	rest.MentorEndpoint()

	rest.Run()
}

//CATATAN
//Jangan lupa bikin text untuk OTP(text yang sekarang masih nyoba-nyoba)
