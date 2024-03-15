package usecase

import (
	"intern-bcc/domain"
	"intern-bcc/internal/repository"
	"intern-bcc/pkg/response"
	"net/http"
)

type IExperieceUsecase interface {
	AddExperience(experience domain.ExperienceRequest, mentorId int) any
}

type ExperienceUsecase struct {
	experienceRepository repository.IExperienceRepository
}

func NewExperienceRepository(experienceRepository repository.IExperienceRepository) IExperieceUsecase {
	return &ExperienceUsecase{experienceRepository}
}

func (u *ExperienceUsecase) AddExperience(experienceRequest domain.ExperienceRequest, mentorId int) any {
	experience := domain.Experiences{
		Experience: experienceRequest.Experience,
		MentorId: mentorId,
	}

	err := u.experienceRepository.AddExperience(&experience)
	if err != nil {
		return response.ErrorObject{
			Code:    http.StatusInternalServerError,
			Message: "an error occured when create experience",
			Err:     err,
		}
	}

	return nil
}
