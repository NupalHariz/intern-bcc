package usecase

import (
	"errors"
	"intern-bcc/domain"
	"intern-bcc/internal/repository"
	"intern-bcc/pkg/response"
	"net/http"
)

type IUniversityUsecase interface {
	CreateUniversity(universityRequest domain.UniversityRequest) error
}

type UniversityUsecase struct {
	universityRepository repository.IUniversityRepository
}

func NewUniversityUsecase(universityRepository repository.IUniversityRepository) IUniversityUsecase {
	return &UniversityUsecase{universityRepository}
}

func (u *UniversityUsecase) CreateUniversity(universityRequest domain.UniversityRequest) error {
	var university domain.Universities
	err := u.universityRepository.GetUniversity(&university, domain.Universities{University: universityRequest.University})
	if err == nil {
		return response.NewError(http.StatusBadRequest, "university already exist", errors.New("can not make same unviersity"))
	}

	newUniversity := domain.Universities{
		University: universityRequest.University,
	}

	err = u.universityRepository.CreateUniversity(&newUniversity)
	if err != nil {
		return response.NewError(http.StatusInternalServerError, "an error occured when creating university", err)
	}

	return nil
}
