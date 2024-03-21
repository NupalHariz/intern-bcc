package main

import (
	"intern-bcc/internal/handler/rest"
	"intern-bcc/internal/repository"
	"intern-bcc/internal/usecase"
	"intern-bcc/pkg/gomail"
	"intern-bcc/pkg/infrastucture"
	"intern-bcc/pkg/infrastucture/database"
	"intern-bcc/pkg/jwt"
	"intern-bcc/pkg/logging"
	"intern-bcc/pkg/middleware"
	"intern-bcc/pkg/midtrans"
	"intern-bcc/pkg/redis"
	"intern-bcc/pkg/supabase"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	env := os.Getenv("ENV")
	if err != nil && env == "" {
		log.Fatal("error loading .env file")
	}
	redis.ConnectToRedis()
	database.ConnectToDB()
	database.Migrate()

	//pkg
	jwt := jwt.JwtInit()
	goMail := gomail.GoMailInit()
	midTrans := midtrans.MidTransInit()
	supabase := supabase.SupabaseInit()
	redis := redis.RedisInit(redis.RDB)
	logging := 	logging.LoggingInit()

	//Repository
	repository := repository.NewRepository(database.DB, repository.RepositoryParam{Redis: redis})

	//Usecase
	usecase := usecase.NewUsecase(usecase.UsecaseParam{
		Repository: repository,
		Jwt:        jwt,
		Supabase:   supabase,
		Midtrans:   midTrans,
		GoMail:     goMail,
		Redis:      redis,
	})

	//Middleware
	middleware := middleware.MiddlerwareInit(jwt, usecase, logging)
	infrastucture.SeedData(database.DB)

	//Rest
	rest := rest.NewRest(gin.New(), usecase, middleware)

	rest.MountEndpoint()

	rest.Run()
}
