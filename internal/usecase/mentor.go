package usecase

import (
	"intern-bcc/domain"
	"intern-bcc/internal/repository"
	"intern-bcc/pkg/jwt"
	"intern-bcc/pkg/middleware"
	"intern-bcc/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type IMentorUsecase interface {
	CreateMentor(c *gin.Context, mentorRequest domain.MentorRequest) any
}

type MentorUsecase struct {
	mentorRepository repository.IMentorRepository
	userRepository   repository.IUserRepository
}

func NewMentorUsecase(mentorRepository repository.IMentorRepository, userRepository repository.IUserRepository) IMentorUsecase {
	return &MentorUsecase{
		mentorRepository: mentorRepository,
		userRepository:   userRepository,
	}
}

func (u *MentorUsecase) CreateMentor(c *gin.Context, mentorRequest domain.MentorRequest) any {
	userId, err := jwt.GetLoginUserId(c)
	if err != nil {
		return response.ErrorObject{
			Code:    http.StatusNotFound,
			Message: "failed to get user id",
			Err:     err,
		}
	}

	var user domain.Users
	err = u.userRepository.GetUser(&user, domain.UserParam{Id: userId})
	if err != nil {
		return response.ErrorObject{
			Code:    http.StatusNotFound,
			Message: "failed to get user",
			Err:     err,
		}
	}

	err = middleware.CheckAdmin(user.IsAdmin)
	if err != nil {
		return response.ErrorObject {
			Code: http.StatusUnauthorized,
			Message: "admin only",
			Err: err,
		}
	}

	newMentor := domain.Mentors{
		Name:          mentorRequest.Name,
		CurrentJob:    mentorRequest.CurrentJob,
		Description:   mentorRequest.Description,
		MentorPicture: mentorRequest.MentorPicture,
	}

	err = u.mentorRepository.CreateMentor(&newMentor)
	if err != nil {
		return response.ErrorObject{
			Code:    http.StatusInternalServerError,
			Message: "an error occured when creating mentor",
			Err:     err,
		}
	}

	return nil
}
