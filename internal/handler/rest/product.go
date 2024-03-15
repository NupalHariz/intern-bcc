package rest

import (
	"fmt"
	"intern-bcc/domain"
	"intern-bcc/pkg/response"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (r *Rest) CreateProduct(c *gin.Context) {
	var productRequest domain.ProductRequest

	err := c.ShouldBindJSON(&productRequest)
	if err != nil {
		response.Failed(c, http.StatusBadRequest, "failed to bind json", err)
		return
	}

	errorObject := r.usecase.ProductUsecase.CreateProduct(c, productRequest)
	if errorObject != nil {
		errorObject := errorObject.(response.ErrorObject)
		response.Failed(c, errorObject.Code, errorObject.Message, errorObject.Err)
		return
	}

	response.SuccessWithoutData(c, "success create product")
}

func (r *Rest) UpdateProduct(c *gin.Context) {
	productIdString := c.Param("productId")
	productId, err := strconv.Atoi(productIdString)
	if err != nil {
		response.Failed(c, http.StatusBadRequest, "failed to parsing product id", err)
		return
	}

	var updateProduct domain.ProductUpdate
	err = c.ShouldBindJSON(&updateProduct)
	if err != nil {
		response.Failed(c, http.StatusBadRequest, "failed to bind json", err)
		return
	}

	product, errorObject := r.usecase.ProductUsecase.UpdateProduct(c, productId, updateProduct)
	if errorObject != nil {
		errorObject := errorObject.(response.ErrorObject)
		response.Failed(c, errorObject.Code, errorObject.Message, errorObject.Err)
		return
	}

	response.Success(c, "success update product", product)
}

func (r *Rest) UploadProductPhoto(c *gin.Context) {
	productIdString := c.Param("productId")
	productId, err := strconv.Atoi(productIdString)
	if err != nil {
		fmt.Println(1)
		response.Failed(c, http.StatusBadRequest, "failed to parsing product id", err)
		return
	}

	productPhoto, err := c.FormFile("product_photo")
	if err != nil {
		fmt.Println(2)
		response.Failed(c, http.StatusBadRequest, "failed to bind request", err)
		return
	}

	product, errorObject := r.usecase.ProductUsecase.UploadProductPhoto(c, productId, productPhoto)
	if errorObject != nil {
		errorObject := errorObject.(response.ErrorObject)
		response.Failed(c, errorObject.Code, errorObject.Message, errorObject.Err)
		return
	}

	response.Success(c, "success upload product product", product)
}
