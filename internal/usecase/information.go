package usecase

import (
	"errors"
	"fmt"
	"intern-bcc/domain"
	"intern-bcc/internal/repository"
	"intern-bcc/pkg/response"
	"intern-bcc/pkg/supabase"
	"mime/multipart"
	"net/http"
	"strings"
	"time"
)

type IInformationUsecase interface {
	CreateInformation(informationRequest domain.InformationRequest) error
	UpdateInformation(informationId int, informationUpdate domain.InformationUpdate) (domain.Information, error)
	UploadInformationPhoto(informationId int, informationPhoto *multipart.FileHeader) (domain.Information, error)
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

func (u *InformationUsecase) CreateInformation(informationRequest domain.InformationRequest) error {
	var category domain.Categories
	err := u.categoryRepository.GetCategory(&category, domain.Categories{Category: informationRequest.Category})
	if err != nil {
		return response.NewError(http.StatusNotFound, "category not found", err)
	}

	if category.Id < 7 || category.Id > 9 {
		return response.NewError(http.StatusBadRequest, "can not use this category for information", errors.New("can not use product category for information"))
	}

	newInformation := domain.Information{
		Title:      informationRequest.Title,
		CategoryId: category.Id,
		Synopsis:   informationRequest.Synopsis,
		Content:    informationRequest.Content,
	}

	err = u.informationRepository.CreateInformation(&newInformation)
	if err != nil {
		return response.NewError(http.StatusInternalServerError, "an error occured when create information", err)
	}

	return nil
}

func (u *InformationUsecase) UpdateInformation(informationId int, informationUpdate domain.InformationUpdate) (domain.Information, error) {
	var information domain.Information
	err := u.informationRepository.GetInformation(&information, domain.Information{Id: informationId})
	if err != nil {
		return domain.Information{}, response.NewError(http.StatusNotFound, "an error occured when get information", err)
	}

	if information.CategoryId != 7 {
		return domain.Information{}, response.NewError(http.StatusBadRequest, "update failed", errors.New("only can update atricle"))
	}

	err = u.informationRepository.UpdateInformation(&informationUpdate, information.Id)
	if err != nil {
		return domain.Information{}, response.NewError(http.StatusInternalServerError, "an error occured when update information", err)
	}

	var updatedInformation domain.Information
	err = u.informationRepository.GetInformation(&updatedInformation, domain.Information{Id: information.Id})
	if err != nil {
		return domain.Information{}, response.NewError(http.StatusInternalServerError, "an error occured when get updated information", err)
	}

	return updatedInformation, nil
}

func (u *InformationUsecase) UploadInformationPhoto(informationId int, informationPhoto *multipart.FileHeader) (domain.Information, error) {
	var information domain.Information
	err := u.informationRepository.GetInformation(&information, domain.Information{Id: informationId})
	if err != nil {
		return domain.Information{}, response.NewError(http.StatusNotFound, "an error occured when get information", err)
	}

	if information.InformationPhoto != "" {
		err = u.supabase.Delete(information.InformationPhoto)
		if err != nil {
			return domain.Information{}, response.NewError(http.StatusInternalServerError, "an error occured when deleting old information photo", err)
		}
	}

	informationPhoto.Filename = fmt.Sprintf("%v-%v", time.Now().String(), informationPhoto.Filename)
		informationPhoto.Filename = strings.Replace(informationPhoto.Filename, " ", "-", -1)
	

	newInformationPhoto, err := u.supabase.Upload(informationPhoto)
	if err != nil {
		return domain.Information{}, response.NewError(http.StatusInternalServerError, "an error occured when upload photo", err)
	}

	err = u.informationRepository.UpdateInformation(&domain.InformationUpdate{InformationPhoto: newInformationPhoto}, information.Id)
	if err != nil {
		return domain.Information{}, response.NewError(http.StatusInternalServerError, "an error occured when update information", err)
	}

	var updatedInformation domain.Information
	err = u.informationRepository.GetInformation(&updatedInformation, domain.Information{Id: information.Id})
	if err != nil {
		return domain.Information{}, response.NewError(http.StatusInternalServerError, "an error occured when get updated information", err)
	}

	return updatedInformation, nil
}
