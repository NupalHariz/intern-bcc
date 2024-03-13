package usecase

import (
	"errors"
	"intern-bcc/domain"
	"intern-bcc/internal/repository"
	"intern-bcc/pkg/jwt"
	"intern-bcc/pkg/response"
	"intern-bcc/pkg/supabase"
	"mime/multipart"
	"net/http"

	"github.com/gin-gonic/gin"
)

type IProductUsecase interface {
	CreateProduct(c *gin.Context, productRequest domain.ProductRequest) any
	UpdateProduct(c *gin.Context, productId int, updateProduct domain.ProductUpdate) (domain.Products, any)
	UploadProductPhoto(c *gin.Context, productId int, productPhoto *multipart.FileHeader) (domain.Products, any)
}

type ProductUsecase struct {
	productRepository  repository.IProductRepository
	merchantRepository repository.IMerchantRepository
	categoryRepository repository.ICategoryRepository
	jwt                jwt.IJwt
	supabase           supabase.ISupabase
}

func NewProductUsecase(productRepository repository.IProductRepository, jwt jwt.IJwt,
	merchantRepository repository.IMerchantRepository, categoryRepository repository.ICategoryRepository,
	supabase supabase.ISupabase) IProductUsecase {
	return &ProductUsecase{
		productRepository:  productRepository,
		jwt:                jwt,
		merchantRepository: merchantRepository,
		categoryRepository: categoryRepository,
		supabase:           supabase,
	}
}

func (u *ProductUsecase) CreateProduct(c *gin.Context, productRequest domain.ProductRequest) any {
	user, err := u.jwt.GetLoginUser(c)
	if err != nil {
		return response.ErrorObject{
			Code:    http.StatusNotFound,
			Message: "failed to get login user",
			Err:     err,
		}
	}

	var merchant domain.Merchants
	err = u.merchantRepository.GetMerchant(&merchant, domain.MerchantParam{UserId: user.Id})
	if err != nil {
		return response.ErrorObject{
			Code:    http.StatusNotFound,
			Message: "failed to get user merchant",
			Err:     err,
		}
	}

	if !merchant.IsActive {
		return response.ErrorObject{
			Code:    http.StatusNotFound,
			Message: "can't make product",
			Err:     errors.New("please verify your merchant"),
		}
	}

	var category domain.Categories
	err = u.categoryRepository.GetCategory(&category, domain.Categories{Category: productRequest.Category})
	if err != nil {
		return response.ErrorObject{
			Code:    http.StatusNotFound,
			Message: "category not found",
			Err:     err,
		}
	}

	newProduct := domain.Products{
		Name:        productRequest.Name,
		MerchantId:  merchant.Id,
		Description: productRequest.Description,
		Price:       productRequest.Price,
		CategoryId:  category.Id,
	}

	err = u.productRepository.CreateProduct(&newProduct)
	if err != nil {
		return response.ErrorObject{
			Code:    http.StatusInternalServerError,
			Message: "an error occured when creating product",
			Err:     err,
		}
	}

	return nil
}

func (u *ProductUsecase) UpdateProduct(c *gin.Context, productId int, updateProduct domain.ProductUpdate) (domain.Products, any) {
	user, err := u.jwt.GetLoginUser(c)
	if err != nil {
		return domain.Products{}, response.ErrorObject{
			Code:    http.StatusNotFound,
			Message: "failed to get login user",
			Err:     err,
		}
	}

	var product domain.Products
	err = u.productRepository.GetProduct(&product, &domain.ProductParam{Id: productId})
	if err != nil {
		return domain.Products{}, response.ErrorObject{
			Code:    http.StatusNotFound,
			Message: "failed to get product",
			Err:     err,
		}
	}

	var merchant domain.Merchants
	err = u.merchantRepository.GetMerchant(&merchant, domain.MerchantParam{UserId: user.Id})
	if err != nil {
		return domain.Products{}, response.ErrorObject{
			Code:    http.StatusNotFound,
			Message: "merchant not found",
			Err:     err,
		}
	}

	if merchant.Id != product.MerchantId {
		return domain.Products{}, response.ErrorObject{
			Code:    http.StatusUnauthorized,
			Message: "access denied",
			Err:     errors.New("can not edit other people merchant"),
		}
	}

	if updateProduct.Name != "" {
		product.Name = updateProduct.Name
	}
	if updateProduct.Price != 0 {
		product.Price = updateProduct.Price
	}
	if updateProduct.Description != "" {
		product.Description = updateProduct.Description
	}

	err = u.productRepository.UpdateProduct(&product)
	if err != nil {
		return domain.Products{}, response.ErrorObject{
			Code:    http.StatusInternalServerError,
			Message: "an error occured when update product",
			Err:     err,
		}
	}

	return product, nil
}

func (u *ProductUsecase) UploadProductPhoto(c *gin.Context, productId int, productPhoto *multipart.FileHeader) (domain.Products, any) {
	user, err := u.jwt.GetLoginUser(c)
	if err != nil {
		return domain.Products{}, response.ErrorObject{
			Code:    http.StatusNotFound,
			Message: "failed to get login user",
			Err:     err,
		}
	}

	var product domain.Products
	err = u.productRepository.GetProduct(&product, &domain.ProductParam{Id: productId})
	if err != nil {
		return domain.Products{}, response.ErrorObject{
			Code:    http.StatusNotFound,
			Message: "failed to get product",
			Err:     err,
		}
	}

	var merchant domain.Merchants
	err = u.merchantRepository.GetMerchant(&merchant, domain.MerchantParam{UserId: user.Id})
	if err != nil {
		return domain.Products{}, response.ErrorObject{
			Code:    http.StatusNotFound,
			Message: "merchant not found",
			Err:     err,
		}
	}

	if merchant.Id != product.MerchantId {
		return domain.Products{}, response.ErrorObject{
			Code:    http.StatusUnauthorized,
			Message: "access denied",
			Err:     errors.New("can not edit other people merchant"),
		}
	}


	if product.ProductPhoto != "" {
		err = u.supabase.Delete(product.ProductPhoto)
		if err != nil {
			return domain.Products{}, response.ErrorObject{
				Code:    http.StatusInternalServerError,
				Message: "error occured when deleting old product photo",
				Err:     err,
			}
		}
	}

	newProductPhoto, err := u.supabase.Upload(productPhoto)
	if err != nil {
		return domain.Products{}, response.ErrorObject{
			Code:    http.StatusInternalServerError,
			Message: "failed to upload photo",
			Err:     err,
		}
	}

	product.ProductPhoto = newProductPhoto
	err = u.productRepository.UpdateProduct(&product)
	if err != nil {
		return domain.Products{}, response.ErrorObject{
			Code:    http.StatusInternalServerError,
			Message: "an error occured when update product",
			Err:     err,
		}
	}

	return product, nil
}
