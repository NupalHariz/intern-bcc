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

	user.PATCH("/:userId", r.middleware.Authentication, r.UpdateUser)
	user.PATCH("/:userId/upload-photo", r.middleware.Authentication, r.UploadUserPhoto)

}

func (r *Rest) MerchantEndpoint() {
	routerGroup := r.router.Group("api/v1")

	merchant := routerGroup.Group("/merchant")
	merchant.POST("/", r.middleware.Authentication, r.CreateMerchant)
	merchant.GET("/verify", r.middleware.Authentication, r.SendOtp)
	merchant.PUT("/verify", r.middleware.Authentication, r.VerifyOtp)
	merchant.PATCH("/:merchantId", r.middleware.Authentication, r.UpdateMerchant)
	merchant.PATCH("/:merchantId/upload-photo", r.middleware.Authentication, r.UploadMerchantPhoto)
}

func (r *Rest) MentorEndpoint() {
	routerGroup := r.router.Group("api/v1")

	mentor := routerGroup.Group("/mentor")
	mentor.POST("/", r.middleware.Authentication, r.middleware.OnlyAdmin, r.CreateMentor)
	mentor.POST("/:mentorId/transaction", r.middleware.Authentication, r.CreateTransaction)
	mentor.POST("/:mentorId/experience", r.middleware.Authentication, r.middleware.OnlyAdmin, r.AddExperience)
	mentor.POST("/payment-callback", r.VerifyTransaction)
}

func (r *Rest) ProductEndpoint() {
	routerGroup := r.router.Group("api/v1")

	product := routerGroup.Group("/product")
	product.POST("/", r.middleware.Authentication, r.CreateProduct)
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
