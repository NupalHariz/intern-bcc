package usecase

import (
	"intern-bcc/domain"
	"intern-bcc/internal/repository"
	"intern-bcc/pkg/response"
	"net/http"
)

type IUniversityUsecase interface {
	CreateUniversity(universityRequest domain.Universities) any
}

type UniversityUsecase struct {
	universityRepository repository.IUniversityRepository
}

func NewUniversityUsecase(universityRepository repository.IUniversityRepository) IUniversityUsecase {
	return &UniversityUsecase{universityRepository}
}

func (u *UniversityUsecase) CreateUniversity(universityRequest domain.Universities) any {
	var university domain.Universities
	err := u.universityRepository.GetUniversity(&university, domain.Universities{University: universityRequest.University})
	if err == nil {
		return response.ErrorObject{
			Code:    http.StatusBadRequest,
			Message: "university already exist",
			Err:     err,
		}
	}

	newUniversity := domain.Universities{
		University: universityRequest.University,
	}

	err = u.universityRepository.CreateUniversity(&newUniversity)
	if err != nil {
		return response.ErrorObject{
			Code:    http.StatusInternalServerError,
			Message: "an error occured when creating university",
			Err:     err,
		}
	}

	return nil
}
