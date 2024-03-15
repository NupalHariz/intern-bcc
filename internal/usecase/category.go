package usecase

import (
	"intern-bcc/domain"
	"intern-bcc/internal/repository"
	"intern-bcc/pkg/response"
	"net/http"
)

type ICategoryUsecase interface {
	CreateCategory(categoryRequest domain.CategoryRequest) any
}

type CategoryUsecase struct {
	categoryRepository repository.ICategoryRepository
}

func NewCategoryUsecase(categoryRepository repository.ICategoryRepository) ICategoryUsecase {
	return &CategoryUsecase{categoryRepository}
}

func (u *CategoryUsecase) CreateCategory(categoryRequest domain.CategoryRequest) any {
	var category domain.Categories
	err := u.categoryRepository.GetCategory(&category, domain.Categories{Category: categoryRequest.Category})
	if err == nil {
		return response.ErrorObject{
			Code:    http.StatusBadRequest,
			Message: "Category already exist",
			Err:     err,
		}
	}

	newCategory := domain.Categories{
		Category: categoryRequest.Category,
	}

	err = u.categoryRepository.CreateCategory(&newCategory)
	if err != nil {
		return response.ErrorObject{
			Code:    http.StatusInternalServerError,
			Message: "an error occured when creating category",
			Err:     err,
		}
	}

	return nil
}
