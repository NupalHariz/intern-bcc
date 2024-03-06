package rest

import (
	"intern-bcc/internal/usecase"
	"log"

	"github.com/gin-gonic/gin"
)

type Rest struct {
	router  *gin.Engine
	usecase *usecase.Usecase
}

func NewRest(usecase *usecase.Usecase) *Rest {
	return &Rest{
		router:  gin.Default(),
		usecase: usecase,
	}
}

func (r *Rest) Run() {
	err := r.router.Run()
	if err != nil {
		log.Fatal("failed to run router")
	}
}
