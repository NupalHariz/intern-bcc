package usecase

import (
	"errors"
	"intern-bcc/domain"
	"intern-bcc/internal/repository"
	"intern-bcc/pkg/response"
	"intern-bcc/pkg/supabase"
	"mime/multipart"
	"net/http"
)

type IInformationUsecase interface {
	CreateInformation(informationRequest domain.InformationRequest) any
	UpdateInformation(informationId int, informationUpdate domain.InformationUpdate) (domain.Information, any)
	UploadInformationPhoto(informationId int, informationPhoto *multipart.FileHeader) (domain.Information, any)
}

type InformationUsecase struct {
	informationRepository repository.IInformationRepository
	categoryRepository    repository.ICategoryRepository
	supabase              supabase.ISupabase
}

func NewInformatinUsecase(informationRepository repository.IInformationRepository, categoryRepository repository.ICategoryRepository,
	supabase supabase.ISupabase) IInformationUsecase {
	return &InformationUsecase{
		informationRepository: informationRepository,
		categoryRepository:    categoryRepository,
		supabase:              supabase,
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

	if category.Id < 7 {
		return response.ErrorObject{
			Code:    http.StatusBadRequest,
			Message: "wrong category",
			Err:     errors.New("not category for information"),
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
			Err:     err,
		}
	}

	return nil
}

func (u *InformationUsecase) UpdateInformation(informationId int, informationUpdate domain.InformationUpdate) (domain.Information, any) {
	var information domain.Information
	err := u.informationRepository.GetInformation(&information, domain.Information{Id: informationId})
	if err != nil {
		return domain.Information{}, response.ErrorObject{
			Code:    http.StatusNotFound,
			Message: "failed to get information",
			Err:     err,
		}
	}

	if information.CategoryId != 7 {
		return domain.Information{}, response.ErrorObject{
			Code:    http.StatusBadRequest,
			Message: "update failed",
			Err:     errors.New("only can update article"),
		}
	}

	if informationUpdate.Content != "" {
		information.Content = informationUpdate.Content
	}
	if informationUpdate.Sysnopsis != "" {
		information.Synopsis = informationUpdate.Sysnopsis
	}

	err = u.informationRepository.UpdateInformation(&information)
	if err != nil {
		return domain.Information{}, response.ErrorObject{
			Code:    http.StatusInternalServerError,
			Message: "an error occured when update information",
			Err:     err,
		}
	}

	return information, nil
}

func (u *InformationUsecase) UploadInformationPhoto(informationId int, informationPhoto *multipart.FileHeader) (domain.Information, any) {
	var information domain.Information
	err := u.informationRepository.GetInformation(&information, domain.Information{Id: informationId})
	if err != nil {
		return domain.Information{}, response.ErrorObject{
			Code:    http.StatusNotFound,
			Message: "failed to get information",
			Err:     err,
		}
	}

	if information.InformationPhoto != "" {
		err = u.supabase.Delete(information.InformationPhoto)
		if err != nil {
			return domain.Information{}, response.ErrorObject{
				Code:    http.StatusInternalServerError,
				Message: "error occured when deleting old information photo",
				Err:     err,
			}
		}
	}

	newInformationPhoto, err := u.supabase.Upload(informationPhoto)
	if err != nil {
		return domain.Information{}, response.ErrorObject{
			Code:    http.StatusInternalServerError,
			Message: "failed to upload photo",
			Err:     err,
		}
	}

	information.InformationPhoto = newInformationPhoto
	err = u.informationRepository.UpdateInformation(&information)
	if err != nil {
		return domain.Information{}, response.ErrorObject{
			Code:    http.StatusInternalServerError,
			Message: "an error occured when update information",
			Err:     err,
		}
	}

	return information, nil
}
