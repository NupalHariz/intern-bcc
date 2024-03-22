package usecase

import (
	"intern-bcc/domain"
	"intern-bcc/internal/repository"
	"intern-bcc/pkg/response"
	"net/http"
)

type IUniversityUsecase interface {
	CreateUniversity(universityRequest domain.Universities) error
}

type UniversityUsecase struct {
	universityRepository repository.IUniversityRepository
}

func NewUniversityUsecase(universityRepository repository.IUniversityRepository) IUniversityUsecase {
	return &UniversityUsecase{universityRepository}
}

func (u *UniversityUsecase) CreateUniversity(universityRequest domain.Universities) error {
	newUniversity := domain.Universities{
		University: universityRequest.University,
	}

	err := u.universityRepository.CreateUniversity(&newUniversity)
	if err != nil {
		return response.NewError(http.StatusInternalServerError, "an error occured when creating university", err)
	}

	return nil
}
