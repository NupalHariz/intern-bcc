package usecase

import (
	"intern-bcc/domain"
	"intern-bcc/internal/repository"
	"intern-bcc/pkg/response"
	"net/http"
)

type ICategoryUsecase interface {
	CreateCategory(categoryRequest domain.Categories) error
}

type CategoryUsecase struct {
	categoryRepository repository.ICategoryRepository
}

func NewCategoryUsecase(categoryRepository repository.ICategoryRepository) ICategoryUsecase {
	return &CategoryUsecase{categoryRepository}
}

func (u *CategoryUsecase) CreateCategory(categoryRequest domain.Categories) error {
	newCategory := domain.Categories{
		Category: categoryRequest.Category,
	}

	err := u.categoryRepository.CreateCategory(&newCategory)
	if err != nil {
		return response.NewError(http.StatusInternalServerError, "an error occured when creating category", err)
	}

	return nil
}
