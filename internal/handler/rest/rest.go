package rest

import (
	"fmt"
	"intern-bcc/internal/usecase"
	"intern-bcc/pkg/middleware"
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

type Rest struct {
	router     *gin.Engine
	usecase    *usecase.Usecase
	middleware middleware.IMiddleware
}

func NewRest(c *gin.Engine, usecase *usecase.Usecase, middleware middleware.IMiddleware) *Rest {
	return &Rest{
		router:     gin.Default(),
		usecase:    usecase,
		middleware: middleware,
	}
}

func (r *Rest) MountEndpoint() {
	version := os.Getenv("VERSION")
	path := fmt.Sprintf("api/%v", version)
	routerGroup := r.router.Group(path, r.middleware.LogEvent)

	routerGroup.POST("/register", r.Register)
	routerGroup.POST("/login", r.Login)
	routerGroup.GET("/recoveryaccount", r.PasswordRecovery)
	routerGroup.PATCH("/recoveryaccount/:name/:verPass", r.ChangePassword)

	user := routerGroup.Group("/user")
	user.PATCH("/:userId", r.middleware.Authentication, r.UpdateUser)
	user.PATCH("/:userId/upload-photo", r.middleware.Authentication, r.UploadUserPhoto)

	profile := routerGroup.Group("/profile")
	profile.GET("/:userId", r.middleware.Authentication, r.GetUser)
	profile.GET("/favourite", r.middleware.Authentication, r.GetLikeProduct)
	profile.GET("/merchant", r.middleware.Authentication, r.GetMerchant)
	profile.GET("/product", r.middleware.Authentication, r.GetOwnProducts)
	profile.GET("/product/:productId", r.middleware.Authentication, r.GetOwnProduct)
	profile.GET("/mentor", r.middleware.Authentication, r.GetOwnMentors)

	merchant := routerGroup.Group("/merchant")
	merchant.POST("/", r.middleware.Authentication, r.CreateMerchant)
	merchant.GET("/verify", r.middleware.Authentication, r.SendOtp)
	merchant.PATCH("/verify", r.middleware.Authentication, r.VerifyOtp)
	merchant.PATCH("/:merchantId", r.middleware.Authentication, r.UpdateMerchant)
	merchant.PATCH("/:merchantId/upload-photo", r.middleware.Authentication, r.UploadMerchantPhoto)

	mentor := routerGroup.Group("/mentor")
	mentor.GET("/:mentorId", r.GetMentor)
	mentor.GET("/", r.GetMentors)
	mentor.POST("/", r.middleware.Authentication, r.middleware.OnlyAdmin, r.CreateMentor)
	mentor.PATCH("/:mentorId", r.middleware.Authentication, r.middleware.OnlyAdmin, r.UpdateMentor)
	mentor.PATCH("/:mentorId/upload-photo", r.middleware.Authentication, r.middleware.OnlyAdmin, r.UploadMentorPicture)
	mentor.POST("/:mentorId/transaction", r.middleware.Authentication, r.CreateTransaction)
	mentor.POST("/:mentorId/experience", r.middleware.Authentication, r.middleware.OnlyAdmin, r.AddExperience)
	mentor.POST("/payment-callback", r.VerifyTransaction)

	product := routerGroup.Group("/product")
	product.GET("/:productId", r.GetProduct)
	product.GET("/", r.GetProducts)
	product.POST("/", r.middleware.Authentication, r.CreateProduct)
	product.PATCH("/:productId", r.middleware.Authentication, r.UpdateProduct)
	product.PATCH("/:productId/product-photo", r.middleware.Authentication, r.UploadProductPhoto)
	product.POST("/:productId", r.middleware.Authentication, r.LikeProduct)
	product.DELETE("/:productId", r.middleware.Authentication, r.DeleteLikeProduct)

	information := routerGroup.Group("/information")
	information.GET("/", r.GetInformations)
	information.GET("/:informationId", r.GetArticle)
	information.POST("/", r.middleware.Authentication, r.middleware.OnlyAdmin, r.CreateInformation)
	information.PATCH("/:informationId", r.middleware.Authentication, r.middleware.OnlyAdmin, r.UpdateInformation)
	information.PATCH("/:informationId/upload-photo", r.middleware.Authentication, r.middleware.OnlyAdmin, r.UploadInformationPhoto)

	category := routerGroup.Group("/category")
	category.POST("/", r.middleware.Authentication, r.middleware.OnlyAdmin, r.CreateCategory)

	province := routerGroup.Group("/province")
	province.POST("/", r.middleware.Authentication, r.middleware.OnlyAdmin, r.CreateProvince)

	university := routerGroup.Group("/university")
	university.POST("/", r.middleware.Authentication, r.middleware.OnlyAdmin, r.CreateUniversity)
}

func (r *Rest) Run() {
	address := os.Getenv("APP_ADDRESS")
	port := os.Getenv("APP_PORT")

	err := r.router.Run(fmt.Sprintf("%s:%s", address, port))
	if err != nil {
		log.Fatal("failed to run router")
	}
}
