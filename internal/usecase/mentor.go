package usecase

import (
	"intern-bcc/domain"
	"intern-bcc/internal/repository"
	"intern-bcc/pkg/jwt"
	"intern-bcc/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type IMentorUsecase interface {
	CreateMentor(c *gin.Context, mentorRequest domain.MentorRequest) any
}

type MentorUsecase struct {
	mentorRepository repository.IMentorRepository
	jwt              jwt.IJwt
}

func NewMentorUsecase(mentorRepository repository.IMentorRepository, jwt jwt.IJwt) IMentorUsecase {
	return &MentorUsecase{
		mentorRepository: mentorRepository,
		jwt:              jwt,
	}
}

func (u *MentorUsecase) CreateMentor(c *gin.Context, mentorRequest domain.MentorRequest) any {
	newMentor := domain.Mentors{
		Name:          mentorRequest.Name,
		CurrentJob:    mentorRequest.CurrentJob,
		Description:   mentorRequest.Description,
		MentorPicture: mentorRequest.MentorPicture,
	}

	err := u.mentorRepository.CreateMentor(&newMentor)
	if err != nil {
		return response.ErrorObject{
			Code:    http.StatusInternalServerError,
			Message: "an error occured when creating mentor",
			Err:     err,
		}
	}

	return nil
}
