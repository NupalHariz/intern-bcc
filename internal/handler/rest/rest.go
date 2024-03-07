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

func NewRest(usecase *usecase.Usecase, middleware middleware.IMiddleware) *Rest {
	return &Rest{
		router:     gin.Default(),
		usecase:    usecase,
		middleware: middleware,
	}
}

func (r *Rest) UserEndpoint() {
	routerGroup := r.router.Group("api/v1")

	routerGroup.POST("/register", r.Register)
	routerGroup.POST("/login", r.Login)

	merchant := routerGroup.Group("/merchant")
	merchant.POST("/", r.middleware.Authentication, r.CreateMerchant)
}



func (r *Rest) Run() {
	err := r.router.Run()
	if err != nil {
		log.Fatal("failed to run router")
	}
}
