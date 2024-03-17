package rest

import (
	"intern-bcc/domain"
	"intern-bcc/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (r *Rest) CreateProduct(c *gin.Context) {
	var productRequest domain.ProductRequest

	err := c.ShouldBindJSON(&productRequest)
	if err != nil {
		response.Failed(c, response.NewError(http.StatusBadRequest, "failed to bind request", err))
		return
	}

	err = r.usecase.ProductUsecase.CreateProduct(c, productRequest)
	if err != nil {
		response.Failed(c, err)
		return
	}

	response.Success(c, "success create product", nil)
}

func (r *Rest) UpdateProduct(c *gin.Context) {
	productIdString := c.Param("productId")
	productId, err := uuid.Parse(productIdString)
	if err != nil {
		response.Failed(c, response.NewError(http.StatusBadRequest, "failed to parsing product id", err))
		return
	}

	var updateProduct domain.ProductUpdate
	err = c.ShouldBindJSON(&updateProduct)
	if err != nil {
		response.Failed(c, response.NewError(http.StatusBadRequest, "failed to bind request", err))
		return
	}

	product, err := r.usecase.ProductUsecase.UpdateProduct(c, productId, updateProduct)
	if err != nil {
		response.Failed(c, err)
		return
	}

	response.Success(c, "success update product", product)
}

func (r *Rest) UploadProductPhoto(c *gin.Context) {
	productIdString := c.Param("productId")
	productId, err := uuid.Parse(productIdString)
	if err != nil {
		response.Failed(c, response.NewError(http.StatusBadRequest, "failed to parsing product id", err))
		return
	}

	productPhoto, err := c.FormFile("product_photo")
	if err != nil {
		response.Failed(c, response.NewError(http.StatusBadRequest, "failed to bind request", err))
		return
	}

	product, err := r.usecase.ProductUsecase.UploadProductPhoto(c, productId, productPhoto)
	if err != nil {
		response.Failed(c, err)
		return
	}

	response.Success(c, "success upload product product", product)
}
