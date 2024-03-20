package usecase

import (
	"context"
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
	GetInformations(ctx context.Context) (domain.InformationResponses, error)
	GetArticle(informationParam domain.InformationParam) (domain.Article, error)
	CreateInformation(informationRequest domain.InformationRequest) error
	UpdateInformation(informationParam domain.InformationParam, informationUpdate domain.InformationUpdate) error
	UploadInformationPhoto(informationParam domain.InformationParam, informationPhoto *multipart.FileHeader) error
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

func (u *InformationUsecase) GetInformations(ctx context.Context) (domain.InformationResponses, error) {
	var articles []domain.Articles
	err := u.informationRepository.GetArticles(ctx, &articles)
	if err != nil {
		return domain.InformationResponses{}, response.NewError(http.StatusInternalServerError, "an error occured when get information", err)
	}

	var otherInformation []domain.Information
	err = u.informationRepository.GetWebinarNCompetition(ctx, &otherInformation)
	if err != nil {
		return domain.InformationResponses{}, response.NewError(http.StatusInternalServerError, "an error occured when get information", err)
	}

	var webinarNCompetitions []domain.WebinarNCompetition

	for _, i := range otherInformation {
		webinarnCompetition := domain.WebinarNCompetition{
			Id:               i.Id,
			Title:            i.Title,
			Category:         i.Category.Category,
			InformationPhoto: i.InformationPhoto,
		}

		webinarNCompetitions = append(webinarNCompetitions, webinarnCompetition)
	}

	informationResponses := domain.InformationResponses{
		Articles:            articles,
		WebinarNCompetition: webinarNCompetitions,
	}

	return informationResponses, nil
}

func (u *InformationUsecase) GetArticle(informationParam domain.InformationParam) (domain.Article, error) {
	var information domain.Information
	err := u.informationRepository.GetInformation(&information, informationParam)
	if err != nil {
		return domain.Article{}, response.NewError(http.StatusInternalServerError, "an error occured when get artivle", err)
	}

	if information.CategoryId != 7 {
		return domain.Article{}, response.NewError(http.StatusBadRequest, "an error occured when get article", errors.New("category id is not article"))
	}

	article := domain.Article{
		Id:               information.Id,
		Title:            information.Title,
		Content:          information.Content,
		InformationPhoto: information.InformationPhoto,
		CreatedAt:        information.CreatedAt.Format(time.RFC822Z),
		UpdatedAt:        information.UpdatedAt.Format(time.RFC822Z),
	}

	return article, err
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

func (u *InformationUsecase) UpdateInformation(informationParam domain.InformationParam, informationUpdate domain.InformationUpdate) error {
	var information domain.Information
	err := u.informationRepository.GetInformation(&information, informationParam)
	if err != nil {
		return response.NewError(http.StatusNotFound, "an error occured when get information", err)
	}

	if information.CategoryId != 7 {
		return response.NewError(http.StatusBadRequest, "update failed", errors.New("only can update atricle"))
	}

	err = u.informationRepository.UpdateInformation(&informationUpdate, information.Id)
	if err != nil {
		return response.NewError(http.StatusInternalServerError, "an error occured when update information", err)
	}

	var updatedInformation domain.Information
	err = u.informationRepository.GetInformation(&updatedInformation, informationParam)
	if err != nil {
		return response.NewError(http.StatusInternalServerError, "an error occured when get updated information", err)
	}

	return nil
}

func (u *InformationUsecase) UploadInformationPhoto(informationParam domain.InformationParam, informationPhoto *multipart.FileHeader) error {
	var information domain.Information
	err := u.informationRepository.GetInformation(&information, informationParam)
	if err != nil {
		return response.NewError(http.StatusNotFound, "an error occured when get information", err)
	}

	if information.InformationPhoto != "" {
		err = u.supabase.Delete(information.InformationPhoto)
		if err != nil {
			return response.NewError(http.StatusInternalServerError, "an error occured when deleting old information photo", err)
		}
	}

	informationPhoto.Filename = fmt.Sprintf("%v-%v", time.Now().String(), informationPhoto.Filename)
	informationPhoto.Filename = strings.Replace(informationPhoto.Filename, " ", "-", -1)

	newInformationPhoto, err := u.supabase.Upload(informationPhoto)
	if err != nil {
		return response.NewError(http.StatusInternalServerError, "an error occured when upload photo", err)
	}

	err = u.informationRepository.UpdateInformation(&domain.InformationUpdate{InformationPhoto: newInformationPhoto}, information.Id)
	if err != nil {
		return response.NewError(http.StatusInternalServerError, "an error occured when update information", err)
	}

	var updatedInformation domain.Information
	err = u.informationRepository.GetInformation(&updatedInformation, informationParam)
	if err != nil {
		return response.NewError(http.StatusInternalServerError, "an error occured when get updated information", err)
	}

	return nil
}
