package main

import (
	rest "intern-bcc/internal/handler/rest"
	"intern-bcc/internal/repository"
	"intern-bcc/internal/usecase"
	"intern-bcc/pkg/infrastucture"
	"intern-bcc/pkg/infrastucture/database"
	"intern-bcc/pkg/jwt"
)

func main() {
	infrastucture.LoadEnv()
	database.ConnectToDB()
	database.Migrate()

	jwt := jwt.JwtInit()

	repository := repository.NewRepository(database.DB)
	usecase := usecase.NewUsecase(usecase.InitParam{Repository: repository, JWT: jwt})
	rest := rest.NewRest(usecase)

	rest.UserEndpoint()

	rest.Run()
}
