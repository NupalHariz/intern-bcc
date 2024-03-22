package usecase

import (
	"intern-bcc/domain"
	"intern-bcc/internal/repository"
	"intern-bcc/pkg/response"
	"net/http"
)

type IExperieceUsecase interface {
	AddExperience(experience domain.Experiences, mentorParam domain.MentorParam) error
}

type ExperienceUsecase struct {
	experienceRepository repository.IExperienceRepository
}

func NewExperienceRepository(experienceRepository repository.IExperienceRepository) IExperieceUsecase {
	return &ExperienceUsecase{experienceRepository}
}

func (u *ExperienceUsecase) AddExperience(experienceRequest domain.Experiences, mentorParam domain.MentorParam) error {
	experience := domain.Experiences{
		Experience: experienceRequest.Experience,
		MentorId:   mentorParam.Id,
	}

	err := u.experienceRepository.AddExperience(&experience)
	if err != nil {
		return response.NewError(http.StatusInternalServerError, "an error occured when create experience", err)
	}

	return nil
}
