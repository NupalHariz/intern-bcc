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
	CreateMentor(mentorRequest domain.MentorRequest) error
	UpdateMentor(mentorId uuid.UUID, mentorUpdate domain.MentorUpdate) (domain.Mentors, error)
	UploadMentorPhoto(mentorId uuid.UUID, mentorPicture *multipart.FileHeader) (domain.Mentors, error)
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

func (u *MentorUsecase) CreateMentor(mentorRequest domain.MentorRequest) error {
	newMentor := domain.Mentors{
		Id:          uuid.New(),
		Name:        mentorRequest.Name,
		CurrentJob:  mentorRequest.CurrentJob,
		Description: mentorRequest.Description,
		Price:       mentorRequest.Price,
	}

	err := u.mentorRepository.CreateMentor(&newMentor)
	if err != nil {
		return response.NewError(http.StatusInternalServerError, "an error occured when create mentor", err)
	}

	return nil
}

func (u *MentorUsecase) UpdateMentor(mentorId uuid.UUID, mentorUpdate domain.MentorUpdate) (domain.Mentors, error) {
	var mentor domain.Mentors
	err := u.mentorRepository.GetMentor(&mentor, domain.Mentors{Id: mentorId})
	if err != nil {
		return domain.Mentors{}, response.NewError(http.StatusNotFound, "an error occured when get mentor", err)
	}

	err = u.mentorRepository.UpdateMentor(&mentorUpdate, mentor.Id)
	if err != nil {
		return domain.Mentors{}, response.NewError(http.StatusInternalServerError, "an error occured when update mentor", err)
	}

	var updatedMentor domain.Mentors
	err = u.mentorRepository.GetMentor(&updatedMentor, domain.Mentors{Id: mentor.Id})
	if err != nil {
		return domain.Mentors{}, response.NewError(http.StatusInternalServerError, "an error occured when get updated mentor", err)
	}

	return updatedMentor, nil
}

func (u *MentorUsecase) UploadMentorPhoto(mentorId uuid.UUID, mentorPicture *multipart.FileHeader) (domain.Mentors, error) {
	var mentor domain.Mentors
	err := u.mentorRepository.GetMentor(&mentor, domain.Mentors{Id: mentorId})
	if err != nil {
		return domain.Mentors{}, response.NewError(http.StatusNotFound, "an error occured when get mentor", err)
	}

	if mentor.MentorPicture != "" {
		err = u.supabase.Delete(mentor.MentorPicture)
		if err != nil {
			return domain.Mentors{}, response.NewError(http.StatusInternalServerError, "an error occured when delete old mentor picture", err)
		}
	}

	mentorPicture.Filename = fmt.Sprintf("%v-%v", time.Now().String(), mentorPicture.Filename)
	mentorPicture.Filename = strings.Replace(mentorPicture.Filename, " ", "-", -1)

	newMentorPicture, err := u.supabase.Upload(mentorPicture)
	if err != nil {
		return domain.Mentors{}, response.NewError(http.StatusInternalServerError, "an error occured when upload photo", err)
	}

	err = u.mentorRepository.UpdateMentor(&domain.MentorUpdate{MentorPicture: newMentorPicture}, mentor.Id)
	if err != nil {
		return domain.Mentors{}, response.NewError(http.StatusInternalServerError, "an error occured when update mentor", err)
	}

	var updatedMentor domain.Mentors
	err = u.mentorRepository.GetMentor(&updatedMentor, domain.Mentors{Id: mentor.Id})
	if err != nil {
		return domain.Mentors{}, response.NewError(http.StatusInternalServerError, "an error occured when get updated mentor", err)
	}

	return updatedMentor, nil
}
