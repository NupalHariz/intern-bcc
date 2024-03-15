package main

import (
	"intern-bcc/internal/handler/rest"
	"intern-bcc/internal/repository"
	"intern-bcc/internal/usecase"
	"intern-bcc/pkg/gomail"
	"intern-bcc/pkg/infrastucture"
	"intern-bcc/pkg/infrastucture/cache"
	"intern-bcc/pkg/infrastucture/database"
	"intern-bcc/pkg/jwt"
	"intern-bcc/pkg/middleware"
	"intern-bcc/pkg/midtrans"
	"intern-bcc/pkg/supabase"

	"github.com/gin-gonic/gin"
)

func main() {
	infrastucture.LoadEnv()
	cache.ConnectToRedis()
	database.ConnectToDB()
	database.Migrate()
	infrastucture.SeedData(database.DB)

	//Repository
	repository := repository.NewRepository(database.DB, cache.RDB)

	//pkg
	jwt := jwt.JwtInit()
	goMail := gomail.GoMailInit()
	midTrans := midtrans.MidTransInit()
	supabase := supabase.SupabaseInit()

	//Usecase
	usecase := usecase.NewUsecase(usecase.UsecaseParam{
		Repository: repository,
		Jwt:        jwt,
		Supabase:   supabase,
		Midtrans:   midTrans,
		GoMail:     goMail,
	})

	//Middleware
	middleware := middleware.MiddlerwareInit(jwt, usecase)

	//Rest
	rest := rest.NewRest(gin.New(), usecase, middleware)

	rest.MerchantEndpoint()
	rest.UserEndpoint()
	rest.MentorEndpoint()
	rest.ProductEndpoint()
	rest.InformationEndpoint()

	rest.Run()
}

//CATATAN URGENT RIL
//Gimana cara dari link email change password bisa ke halaman bikin password?
//Kasus gimana kalo ada orang yang file fotonya sama
//Bikin env untuk host dan port

//To-Do Today
//Bikin tabel khsusu untuk univ dan provinsi
//Kirim pesan ke email yang bagus
//Bikin respon error yang bener(kasih info succes: true atau false)
//Bikin respon yang ga ngebalikin nilai
//Endpoint dijadiin 1 aja

//CATATAN
//Jangan lupa bikin text untuk OTP(text yang sekarang masih nyoba-nyoba)
//ENV nya samain buat yang kayak di deploy
//Jangan lupa benerin respon untuk update
