package rest

import (
	"intern-bcc/internal/usecase"
	"intern-bcc/pkg/middleware"
	"log"

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

func (r *Rest) UserEndpoint() {
	routerGroup := r.router.Group("api/v1")

	user := routerGroup.Group("/user")

	routerGroup.POST("/register", r.Register)
	routerGroup.POST("/login", r.Login)
	routerGroup.POST("/recoveryaccount", r.PasswordRecovery)
	routerGroup.PATCH("/recoveryaccount/:email/:verPass", r.ChangePassword)

	user.PATCH("/:userId", r.middleware.Authentication, r.UpdateUser)
	user.PATCH("/:userId/upload-photo", r.middleware.Authentication, r.UploadUserPhoto)

}

func (r *Rest) MerchantEndpoint() {
	routerGroup := r.router.Group("api/v1")

	merchant := routerGroup.Group("/merchant")
	merchant.POST("/", r.middleware.Authentication, r.CreateMerchant)
	merchant.GET("/verify", r.middleware.Authentication, r.SendOtp)
	merchant.PATCH("/verify", r.middleware.Authentication, r.VerifyOtp)
	merchant.PATCH("/:merchantId", r.middleware.Authentication, r.UpdateMerchant)
	merchant.PATCH("/:merchantId/upload-photo", r.middleware.Authentication, r.UploadMerchantPhoto)
}

func (r *Rest) MentorEndpoint() {
	routerGroup := r.router.Group("api/v1")

	mentor := routerGroup.Group("/mentor")
	mentor.POST("/", r.middleware.Authentication, r.middleware.OnlyAdmin, r.CreateMentor)
	mentor.PATCH("/:mentorId", r.middleware.Authentication, r.middleware.OnlyAdmin, r.UpdateMentor)
	mentor.PATCH("/:mentorId/upload-photo", r.middleware.Authentication, r.middleware.OnlyAdmin, r.UploadMentorPicture)
	mentor.POST("/:mentorId/transaction", r.middleware.Authentication, r.CreateTransaction)
	mentor.POST("/:mentorId/experience", r.middleware.Authentication, r.middleware.OnlyAdmin, r.AddExperience)
	mentor.POST("/payment-callback", r.VerifyTransaction)
}

func (r *Rest) ProductEndpoint() {
	routerGroup := r.router.Group("api/v1")

	product := routerGroup.Group("/product")
	product.POST("/", r.middleware.Authentication, r.CreateProduct)
	product.PATCH("/:productId", r.middleware.Authentication, r.UpdateProduct)
	product.PATCH("/:productId/product-photo", r.middleware.Authentication, r.UploadProductPhoto)
	product.POST("/:productId", r.middleware.Authentication, r.LikeProduct)
	product.DELETE("/:productId", r.middleware.Authentication, r.DeleteLikeProduct)
}

func (r *Rest) InformationEndpoint() {
	routerGroup := r.router.Group("api/v1")

	information := routerGroup.Group("/information")
	information.POST("/", r.middleware.Authentication, r.middleware.OnlyAdmin, r.CreateInformation)
	information.PATCH("/:informationId", r.middleware.Authentication, r.middleware.OnlyAdmin, r.UpdateInformation)
	information.PATCH("/:informationId/upload-photo", r.middleware.Authentication, r.middleware.OnlyAdmin, r.UploadInformationPhoto)
}

func (r *Rest) Run() {
	err := r.router.Run()
	if err != nil {
		log.Fatal("failed to run router")
	}
}

// // domain driven design
// -- app
// -- -- user
// -- -- -- repository
// -- -- -- usecase
// -- -- -- handler

// // event driven design
// -- repository
// -- -- user repo
// -- usecase
// -- -- user uc
// -- handler
// -- -- user handler
