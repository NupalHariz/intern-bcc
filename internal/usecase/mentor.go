package usecase

import (
	"fmt"
	"intern-bcc/domain"
	"intern-bcc/internal/repository"
	"intern-bcc/pkg/jwt"
	"intern-bcc/pkg/response"
	"intern-bcc/pkg/supabase"
	"mime/multipart"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
)

type IMentorUsecase interface {
	CreateMentor(mentorRequest domain.MentorRequest) any
	UpdateMentor(mentorId uuid.UUID, mentorUpdate domain.MentorUpdate) (domain.Mentors, any)
	UploadMentorPhoto(mentorId uuid.UUID, mentorPicture *multipart.FileHeader) (domain.Mentors, any)
}

type MentorUsecase struct {
	mentorRepository repository.IMentorRepository
	jwt              jwt.IJwt
	supabase         supabase.ISupabase
}

func NewMentorUsecase(mentorRepository repository.IMentorRepository, jwt jwt.IJwt, supabase supabase.ISupabase) IMentorUsecase {
	return &MentorUsecase{
		mentorRepository: mentorRepository,
		jwt:              jwt,
		supabase:         supabase,
	}
}

func (u *MentorUsecase) CreateMentor(mentorRequest domain.MentorRequest) any {
	newMentor := domain.Mentors{
		Id:          uuid.New(),
		Name:        mentorRequest.Name,
		CurrentJob:  mentorRequest.CurrentJob,
		Description: mentorRequest.Description,
		Price:       mentorRequest.Price,
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

func (u *MentorUsecase) UpdateMentor(mentorId uuid.UUID, mentorUpdate domain.MentorUpdate) (domain.Mentors, any) {
	var mentor domain.Mentors
	err := u.mentorRepository.GetMentor(&mentor, domain.Mentors{Id: mentorId})
	if err != nil {
		return domain.Mentors{}, response.ErrorObject{
			Code:    http.StatusNotFound,
			Message: "failed to get mentor",
			Err:     err,
		}
	}

	if mentorUpdate.CurrentJob != "" {
		mentor.CurrentJob = mentorUpdate.CurrentJob
	}
	if mentorUpdate.Description != "" {
		mentor.Description = mentorUpdate.Description
	}
	if mentorUpdate.Price != 0 {
		mentor.Price = mentorUpdate.Price
	}

	err = u.mentorRepository.UpdateMentor(&mentor)
	if err != nil {
		return domain.Mentors{}, response.ErrorObject{
			Code:    http.StatusInternalServerError,
			Message: "an error occured when update mentor",
			Err:     err,
		}
	}

	return mentor, nil
}

func (u *MentorUsecase) UploadMentorPhoto(mentorId uuid.UUID, mentorPicture *multipart.FileHeader) (domain.Mentors, any) {
	var mentor domain.Mentors
	err := u.mentorRepository.GetMentor(&mentor, domain.Mentors{Id: mentorId})
	if err != nil {
		return domain.Mentors{}, response.ErrorObject{
			Code:    http.StatusNotFound,
			Message: "failed to get mentor",
			Err:     err,
		}
	}

	if mentor.MentorPicture != "" {
		err = u.supabase.Delete(mentor.MentorPicture)
		if err != nil {
			return domain.Mentors{}, response.ErrorObject{
				Code:    http.StatusInternalServerError,
				Message: "error occured when deleting old mentor picture",
				Err:     err,
			}
		}
	}

	mentorPicture.Filename = fmt.Sprintf("%v-%v", time.Now().String(), mentorPicture.Filename)
	if strings.Contains(mentorPicture.Filename, " ") {
		mentorPicture.Filename = strings.Replace(mentorPicture.Filename, " ", "-", -1)
	}

	newMentorPicture, err := u.supabase.Upload(mentorPicture)
	if err != nil {
		return domain.Mentors{}, response.ErrorObject{
			Code:    http.StatusInternalServerError,
			Message: "failed to upload photo",
			Err:     err,
		}
	}

	mentor.MentorPicture = newMentorPicture
	err = u.mentorRepository.UpdateMentor(&mentor)
	if err != nil {
		return domain.Mentors{}, response.ErrorObject{
			Code:    http.StatusInternalServerError,
			Message: "an error occured when update mentor",
			Err:     err,
		}
	}

	return mentor, nil
}
