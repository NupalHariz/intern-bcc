package usecase

import (
	"errors"
	"intern-bcc/domain"
	"intern-bcc/internal/repository"
	"intern-bcc/pkg/jwt"
	"intern-bcc/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type IProductUsecase interface {
	CreateProduct(c *gin.Context, productRequest domain.ProductRequest) any
}

type ProductUsecase struct {
	productRepository  repository.IProductRepository
	merchantRepository repository.IMerchantRepository
	categoryRepository repository.ICategoryRepository
	jwt                jwt.IJwt
}

func NewProductUsecase(productRepository repository.IProductRepository, jwt jwt.IJwt,
	merchantRepository repository.IMerchantRepository, categoryRepository repository.ICategoryRepository) IProductUsecase {
	return &ProductUsecase{
		productRepository:  productRepository,
		jwt:                jwt,
		merchantRepository: merchantRepository,
		categoryRepository: categoryRepository,
	}
}

func (u *ProductUsecase) CreateProduct(c *gin.Context, productRequest domain.ProductRequest,) any {
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

	if !merchant.IsActive{
		return response.ErrorObject{
			Code:    http.StatusNotFound,
			Message: "can't make product",
			Err:     errors.New("please verify your merchant"),
		}
	}

	var category domain.Categories
	err = u.categoryRepository.GetCategory(&category)
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
		CategoryId: category.Id,
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
