package rest

import (
	"fmt"
	"intern-bcc/domain"
	"intern-bcc/pkg/response"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (r *Rest) CreateInformation(c *gin.Context) {
	var informationRequest domain.InformationRequest

	err := c.ShouldBindJSON(&informationRequest)
	if err != nil {
		response.Failed(c, response.NewError(http.StatusBadRequest, "failed to bind request", err))
	}

	err = r.usecase.InformationUsecase.CreateInformation(informationRequest)
	if err != nil {
		response.Failed(c, err)
		return
	}

	response.Success(c, "success create information", nil)
}

func (r *Rest) UpdateInformation(c *gin.Context) {
	informationIdString := c.Param("informationId")
	informationId, err := strconv.Atoi(informationIdString)
	if err != nil {
		response.Failed(c, response.NewError(http.StatusBadRequest, "failed to parsing mentor id", err))
		return
	}

	var informationUpdate domain.InformationUpdate
	err = c.ShouldBindJSON(&informationUpdate)
	if err != nil {
		response.Failed(c, response.NewError(http.StatusBadRequest, "failed to bind request", err))
		return
	}

	information, err := r.usecase.InformationUsecase.UpdateInformation(informationId, informationUpdate)
	if err != nil {
		response.Failed(c, err)
		return
	}

	fmt.Println(information)
	response.Success(c, "success update information", information)
}

func (r *Rest) UploadInformationPhoto(c *gin.Context) {
	informationIdString := c.Param("informationId")
	informationId, err := strconv.Atoi(informationIdString)
	if err != nil {
		response.Failed(c, response.NewError(http.StatusBadRequest, "failed to parsing mentor id", err))
		return
	}

	informationPhoto, err := c.FormFile("information_photo")
	if err != nil {
		response.Failed(c, response.NewError(http.StatusBadRequest, "failed to bind request", err))
		return
	}

	information, err := r.usecase.InformationUsecase.UploadInformationPhoto(informationId, informationPhoto)
	if err != nil {
		response.Failed(c, err)
		return
	}

	response.Success(c, "success upload information photo", information)
}
