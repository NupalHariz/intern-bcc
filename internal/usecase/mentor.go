package usecase

import (
	"intern-bcc/domain"
	"intern-bcc/internal/repository"
	"intern-bcc/pkg/jwt"
	"intern-bcc/pkg/response"
	"intern-bcc/pkg/supabase"
	"mime/multipart"
	"net/http"
)

type IMentorUsecase interface {
	CreateMentor(mentorRequest domain.MentorRequest) any
	UpdateMentor(mentorId int, mentorUpdate domain.MentorUpdate) (domain.Mentors, any)
	UploadMentorPhoto(mentorId int, mentorPicture *multipart.FileHeader) (domain.Mentors, any)
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

func (u *MentorUsecase) UpdateMentor(mentorId int, mentorUpdate domain.MentorUpdate) (domain.Mentors, any) {
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

func (u *MentorUsecase) UploadMentorPhoto(mentorId int, mentorPicture *multipart.FileHeader) (domain.Mentors, any) {
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
		return domain.Mentors{}, response.ErrorObject {
			Code: http.StatusInternalServerError,
			Message: "an error occured when update mentor",
			Err: err,
		}
	}

	return mentor, nil
}
