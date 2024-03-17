package usecase

import (
	"errors"
	"fmt"
	"intern-bcc/domain"
	"intern-bcc/internal/repository"
	"intern-bcc/pkg/jwt"
	"intern-bcc/pkg/response"
	"intern-bcc/pkg/supabase"
	"mime/multipart"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type IProductUsecase interface {
	CreateProduct(c *gin.Context, productRequest domain.ProductRequest) error
	UpdateProduct(c *gin.Context, productId uuid.UUID, updateProduct domain.ProductUpdate) (domain.Products, error)
	UploadProductPhoto(c *gin.Context, productId uuid.UUID, productPhoto *multipart.FileHeader) (domain.Products, error)
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

func (u *ProductUsecase) CreateProduct(c *gin.Context, productRequest domain.ProductRequest) error {
	user, err := u.jwt.GetLoginUser(c)
	if err != nil {
		return response.NewError(http.StatusNotFound, "an error occured when get login user", err)
	}

	var merchant domain.Merchants
	err = u.merchantRepository.GetMerchant(&merchant, domain.MerchantParam{UserId: user.Id})
	if err != nil {
		return response.NewError(http.StatusNotFound, "an er", err)
	}

	if !merchant.IsActive {
		return response.NewError(http.StatusNotFound, "failed to create product", errors.New("please verify your merchant"))
	}

	var category domain.Categories
	err = u.categoryRepository.GetCategory(&category, domain.Categories{Category: productRequest.Category})
	if err != nil {
		return response.NewError(http.StatusNotFound, "category not found", err)
	}

	if category.Id > 6 {
		return response.NewError(http.StatusBadRequest, "can no use this category for product", errors.New("can not use information category"))

	}

	newProduct := domain.Products{
		Id:          uuid.New(),
		Name:        productRequest.Name,
		MerchantId:  merchant.Id,
		Description: productRequest.Description,
		Price:       productRequest.Price,
		CategoryId:  category.Id,
	}

	err = u.productRepository.CreateProduct(&newProduct)
	if err != nil {
		return response.NewError(http.StatusInternalServerError, "an error occured when creating product", err)
	}

	return nil
}

func (u *ProductUsecase) UpdateProduct(c *gin.Context, productId uuid.UUID, updateProduct domain.ProductUpdate) (domain.Products, error) {
	user, err := u.jwt.GetLoginUser(c)
	if err != nil {
		return domain.Products{}, response.NewError(http.StatusNotFound, "an error occured when get login user", err)
	}

	var product domain.Products
	err = u.productRepository.GetProduct(&product, domain.ProductParam{Id: productId})
	if err != nil {
		return domain.Products{}, response.NewError(http.StatusNotFound, "an error occured when get product", err)
	}

	var merchant domain.Merchants
	err = u.merchantRepository.GetMerchant(&merchant, domain.MerchantParam{UserId: user.Id})
	if err != nil {
		return domain.Products{}, response.NewError(http.StatusNotFound, "merchant not found", err)
	}

	if merchant.Id != product.MerchantId {
		return domain.Products{}, response.NewError(http.StatusUnauthorized, "access denied", errors.New("can not edit other people merchant"))
	}

	err = u.productRepository.UpdateProduct(&updateProduct, product.Id)
	if err != nil {
		return domain.Products{}, response.NewError(http.StatusInternalServerError, "an error occured when update product", err)
	}

	var updatedProduct domain.Products
	err = u.productRepository.GetProduct(&updatedProduct, domain.ProductParam{Id: product.Id})
	if err != nil {
		return domain.Products{}, response.NewError(http.StatusInternalServerError, "an error occured when get updated product", err)
	}

	return updatedProduct, nil
}

func (u *ProductUsecase) UploadProductPhoto(c *gin.Context, productId uuid.UUID, productPhoto *multipart.FileHeader) (domain.Products, error) {
	user, err := u.jwt.GetLoginUser(c)
	if err != nil {
		return domain.Products{}, response.NewError(http.StatusNotFound, "an error occured when get login user", err)
	}

	var product domain.Products
	err = u.productRepository.GetProduct(&product, domain.ProductParam{Id: productId})
	if err != nil {
		return domain.Products{}, response.NewError(http.StatusNotFound, "an error occured when get product", err)
	}

	var merchant domain.Merchants
	err = u.merchantRepository.GetMerchant(&merchant, domain.MerchantParam{UserId: user.Id})
	if err != nil {
		return domain.Products{}, response.NewError(http.StatusNotFound, "merchant not found", err)
	}

	if merchant.Id != product.MerchantId {
		return domain.Products{}, response.NewError(http.StatusUnauthorized, "access denied", errors.New("can not edit other people merchant"))
	}

	if product.ProductPhoto != "" {
		err = u.supabase.Delete(product.ProductPhoto)
		if err != nil {
			return domain.Products{}, response.NewError(http.StatusInternalServerError, "error occured when deleting old product photo", err)
		}
	}

	productPhoto.Filename = fmt.Sprintf("%v-%v", time.Now().String(), productPhoto.Filename)
	productPhoto.Filename = strings.Replace(productPhoto.Filename, " ", "-", -1)

	newProductPhoto, err := u.supabase.Upload(productPhoto)
	if err != nil {
		return domain.Products{}, response.NewError(http.StatusInternalServerError, "failed to upload photo", err)
	}

	err = u.productRepository.UpdateProduct(&domain.ProductUpdate{ProductPhoto: newProductPhoto}, product.Id)
	if err != nil {
		return domain.Products{}, response.NewError(http.StatusInternalServerError, "an error occured when update product", err)
	}

	var updatedProduct domain.Products
	err = u.productRepository.GetProduct(&updatedProduct, domain.ProductParam{Id: product.Id})
	if err != nil {
		return domain.Products{}, response.NewError(http.StatusInternalServerError, "an error occured when get updated product", err)
	}

	return updatedProduct, nil
}
