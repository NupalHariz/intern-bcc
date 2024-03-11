package rest

import (
	"intern-bcc/internal/handler"
	"intern-bcc/pkg/middleware"
	"log"

	"github.com/gin-gonic/gin"
)

type Rest struct {
	router *gin.Engine
}

func NewRest(c *gin.Engine) *Rest {
	return &Rest{
		router: gin.Default(),
	}
}

func (r *Rest) UserEndpoint(userHandler *handler.UserHandler) {
	routerGroup := r.router.Group("api/v1")

	routerGroup.POST("/register", userHandler.Register)
	routerGroup.POST("/login", userHandler.Login)
}

func (r *Rest) MerchantEndpoint(merchantHandler *handler.MerchantHandler) {
	routerGroup := r.router.Group("api/v1")

	merchant := routerGroup.Group("/merchant")
	merchant.POST("/", middleware.Authentication, merchantHandler.CreateMerchant)
	merchant.GET("/verify", middleware.Authentication, merchantHandler.SendOtp)
	merchant.PUT("/verify", middleware.Authentication, merchantHandler.VerifyOtp)
}

func (r *Rest) MentorEndpoint(mentorHandler *handler.MentorHandler, transactionHandler *handler.TransactionHandler) {
	routerGroup := r.router.Group("api/v1")

	mentor := routerGroup.Group("/mentor")
	mentor.POST("/", middleware.Authentication, mentorHandler.CreateMentor)
	mentor.POST("/:mentorId/transaction", middleware.Authentication, transactionHandler.CreateTransaction)

}

func (r *Rest) Run() {
	err := r.router.Run()
	if err != nil {
		log.Fatal("failed to run router")
	}
}
