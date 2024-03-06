package main

import (
	rest "intern-bcc/internal/handler/rest"
	"intern-bcc/internal/repository"
	"intern-bcc/internal/usecase"
	"intern-bcc/pkg/infrastucture"
	"intern-bcc/pkg/infrastucture/database"
)

func main() {
	infrastucture.LoadEnv()
	database.ConnectToDB()
	database.Migrate()

	repository := repository.NewRepository(database.DB)
	usecase := usecase.NewUsecase(usecase.InitParam{Repository: repository})
	rest := rest.NewRest(usecase)

	rest.Run()
}
