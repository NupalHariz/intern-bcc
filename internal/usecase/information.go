package usecase

import (
	"intern-bcc/domain"
	"intern-bcc/internal/repository"
	"intern-bcc/pkg/response"
	"net/http"
)

type IInformationUsecase interface{
	CreateInformation(informationRequest domain.InformationRequest) any
}

type InformationUsecase struct {
	informationRepository repository.IInformationRepository
	categoryRepository    repository.ICategoryRepository
}

func NewInformatinUsecase(informationRepository repository.IInformationRepository, categoryRepository repository.ICategoryRepository) IInformationUsecase {
	return &InformationUsecase{
		informationRepository: informationRepository,
		categoryRepository:    categoryRepository,
	}
}

func (u *InformationUsecase) CreateInformation(informationRequest domain.InformationRequest) any {
	var category domain.Categories

	err := u.categoryRepository.GetCategory(&category, domain.Categories{Category: informationRequest.Category})
	if err != nil {
		return response.ErrorObject{
			Code:    http.StatusNotFound,
			Message: "category not found",
			Err:     err,
		}
	}

	newInformation := domain.Information{
		Title:      informationRequest.Title,
		CategoryId: category.Id,
		Synopsis:   informationRequest.Synopsis,
		Content:    informationRequest.Content,
	}

	err = u.informationRepository.CreateInformation(&newInformation)
	if err != nil {
		return response.ErrorObject{
			Code:    http.StatusInternalServerError,
			Message: "an error occured when create information",
			Err: err,
		}
	}

	return nil
}
