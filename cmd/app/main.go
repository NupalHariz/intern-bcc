package main

import (
	rest "intern-bcc/internal/handler/rest"
	"intern-bcc/internal/repository"
	"intern-bcc/internal/usecase"
	"intern-bcc/pkg/infrastucture"
	"intern-bcc/pkg/infrastucture/cache"
	"intern-bcc/pkg/infrastucture/database"
	"intern-bcc/pkg/jwt"
	"intern-bcc/pkg/middleware"
)

func main() {
	infrastucture.LoadEnv()
	cache.ConnectToRedis()
	database.ConnectToDB()
	database.Migrate()

	jwt := jwt.JwtInit()

	repository := repository.NewRepository(database.DB)
	usecase := usecase.NewUsecase(usecase.InitParam{Repository: repository, JWT: jwt})
	middleware := middleware.MiddlerwareInit(jwt, usecase)
	rest := rest.NewRest(usecase, middleware)

	rest.UserEndpoint()

	rest.Run()
}
