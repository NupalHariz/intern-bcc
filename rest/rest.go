package rest

import (
	"intern-bcc/internal/handler"
	"intern-bcc/pkg/middleware"
	"log"

	"github.com/gin-gonic/gin"
)

type Rest struct {
	router             *gin.Engine
	userHandler        *handler.UserHandler
	merchantHandler    *handler.MerchantHandler
	mentorHandler      *handler.MentorHandler
	transactionHandler *handler.TransactionHandler
	middleware         middleware.IMiddleware
}

func NewRest(c *gin.Engine, userHandler *handler.UserHandler,
	merchantHandler *handler.MerchantHandler, mentorHandler *handler.MentorHandler,
	transactionHandler *handler.TransactionHandler, middleware middleware.IMiddleware) *Rest {
	return &Rest{
		router:             gin.Default(),
		userHandler:        userHandler,
		merchantHandler:    merchantHandler,
		mentorHandler:      mentorHandler,
		transactionHandler: transactionHandler,
		middleware:         middleware,
	}
}

func (r *Rest) UserEndpoint() {
	routerGroup := r.router.Group("api/v1")

	user := routerGroup.Group("/user")

	routerGroup.POST("/register", r.userHandler.Register)
	routerGroup.POST("/login", r.userHandler.Login)

	user.PATCH("/:userId", r.middleware.Authentication, r.userHandler.UpdateUser)
	user.PATCH("/:userId/upload-photo", r.middleware.Authentication, r.userHandler.UploadPhoto)

}

func (r *Rest) MerchantEndpoint() {
	routerGroup := r.router.Group("api/v1")

	merchant := routerGroup.Group("/merchant")
	merchant.POST("/", r.middleware.Authentication, r.merchantHandler.CreateMerchant)
	merchant.GET("/verify", r.middleware.Authentication, r.merchantHandler.SendOtp)
	merchant.PUT("/verify", r.middleware.Authentication, r.merchantHandler.VerifyOtp)
}

func (r *Rest) MentorEndpoint() {
	routerGroup := r.router.Group("api/v1")

	mentor := routerGroup.Group("/mentor")
	mentor.POST("/", r.middleware.Authentication, r.middleware.OnlyAdmin, r.mentorHandler.CreateMentor)
	mentor.POST("/:mentorId/transaction", r.middleware.Authentication, r.transactionHandler.CreateTransaction)
	mentor.POST("/payment-callback", r.transactionHandler.VerifyTransaction)

}

func (r *Rest) Run() {
	err := r.router.Run()
	if err != nil {
		log.Fatal("failed to run router")
	}
}
