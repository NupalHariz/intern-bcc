package main

import (
	"intern-bcc/internal/handler"
	"intern-bcc/internal/repository"
	"intern-bcc/internal/usecase"
	"intern-bcc/pkg/infrastucture"
	"intern-bcc/pkg/infrastucture/cache"
	"intern-bcc/pkg/infrastucture/database"
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

	//Usecase
	merchantUsecase := usecase.NewMerchantUsecase(merchantRepository, merchantRedis, userRepository)
	userUsecase := usecase.NewUserUsecase(userRepository)

	//Handler
	merchantHandler := handler.NewMerchantHandler(merchantUsecase)
	userHandler := handler.NewUserHandler(userUsecase)

	rest := rest.NewRest(gin.New())

	rest.MerchantEndpoint(merchantHandler)
	rest.UserEndpoint(userHandler)

	rest.Run()
}

//CATATAN
//Jangan lupa bikin text untuk OTP(text yang sekarang masih nyoba-nyoba)
