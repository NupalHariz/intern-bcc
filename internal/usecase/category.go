package usecase

import (
	"errors"
	"intern-bcc/domain"
	"intern-bcc/internal/repository"
	"intern-bcc/pkg/response"
	"net/http"
)

type ICategoryUsecase interface {
	CreateCategory(categoryRequest domain.CategoryRequest) error
}

type CategoryUsecase struct {
	categoryRepository repository.ICategoryRepository
}

func NewCategoryUsecase(categoryRepository repository.ICategoryRepository) ICategoryUsecase {
	return &CategoryUsecase{categoryRepository}
}

func (u *CategoryUsecase) CreateCategory(categoryRequest domain.CategoryRequest) error {
	var category domain.Categories
	err := u.categoryRepository.GetCategory(&category, domain.Categories{Category: categoryRequest.Category})
	if err == nil {
		return response.NewError(http.StatusBadRequest, "category already exist", errors.New("can not make same category"))
	}

	newCategory := domain.Categories{
		Category: categoryRequest.Category,
	}

	err = u.categoryRepository.CreateCategory(&newCategory)
	if err != nil {
		return response.NewError(http.StatusInternalServerError, "an error occured when creating category", err)
	}

	return nil
}
