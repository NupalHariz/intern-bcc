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
	GetMentor(mentorParam domain.MentorParam) (domain.MentorResponse, error)
	GetMentors() ([]domain.MentorResponses, error)
	CreateMentor(mentorRequest domain.MentorRequest) error
	UpdateMentor(mentorParam domain.MentorParam, mentorUpdate domain.MentorUpdate) error
	UploadMentorPhoto(mentorParam domain.MentorParam, mentorPicture *multipart.FileHeader) error
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

func (u *MentorUsecase) GetMentor(mentorParam domain.MentorParam) (domain.MentorResponse, error) {
	var mentor domain.Mentors
	err := u.mentorRepository.GetMentor(&mentor, mentorParam)
	if err != nil {
		return domain.MentorResponse{}, response.NewError(http.StatusNotFound, "an error occured when get mentors", err)
	}

	fmt.Println(mentor.Experiences)

	mentorResponse := domain.MentorResponse{
		Id:          mentor.Id,
		Name:        mentor.Name,
		CurrentJob:  mentor.CurrentJob,
		Description: mentor.Description,
		Price:       mentor.Price,
		Experiences: []string{
			mentor.Experiences[0].Experience,
			mentor.Experiences[1].Experience,
			mentor.Experiences[2].Experience,
		},
	}

	return mentorResponse, nil
}

func (u *MentorUsecase) GetMentors() ([]domain.MentorResponses, error) {
	var mentors []domain.Mentors
	err := u.mentorRepository.GetMentors(&mentors)
	if err != nil {
		return []domain.MentorResponses{}, response.NewError(http.StatusInternalServerError, "an error occured when get mentors", err)
	}

	var mentorResponses []domain.MentorResponses
	for _, m := range mentors {
		mentorResponse := domain.MentorResponses{
			Id:            m.Id,
			Name:          m.Name,
			CurrentJob:    m.CurrentJob,
			Price:         m.Price,
			MentorPicture: m.MentorPicture,
		}

		mentorResponses = append(mentorResponses, mentorResponse)
	}

	return mentorResponses, nil
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

func (u *MentorUsecase) UpdateMentor(mentorParam domain.MentorParam, mentorUpdate domain.MentorUpdate) error {
	var mentor domain.Mentors
	err := u.mentorRepository.GetMentor(&mentor, mentorParam)
	if err != nil {
		return response.NewError(http.StatusNotFound, "an error occured when get mentor", err)
	}

	err = u.mentorRepository.UpdateMentor(&mentorUpdate, mentor.Id)
	if err != nil {
		return response.NewError(http.StatusInternalServerError, "an error occured when update mentor", err)
	}

	var updatedMentor domain.Mentors
	err = u.mentorRepository.GetMentor(&updatedMentor, mentorParam)
	if err != nil {
		return response.NewError(http.StatusInternalServerError, "an error occured when get updated mentor", err)
	}

	return nil
}

func (u *MentorUsecase) UploadMentorPhoto(mentorParam domain.MentorParam, mentorPicture *multipart.FileHeader) error {
	var mentor domain.Mentors
	err := u.mentorRepository.GetMentor(&mentor, mentorParam)
	if err != nil {
		return response.NewError(http.StatusNotFound, "an error occured when get mentor", err)
	}

	if mentor.MentorPicture != "" {
		err = u.supabase.Delete(mentor.MentorPicture)
		if err != nil {
			return response.NewError(http.StatusInternalServerError, "an error occured when delete old mentor picture", err)
		}
	}

	mentorPicture.Filename = fmt.Sprintf("%v-%v", time.Now().String(), mentorPicture.Filename)
	mentorPicture.Filename = strings.Replace(mentorPicture.Filename, " ", "-", -1)

	newMentorPicture, err := u.supabase.Upload(mentorPicture)
	if err != nil {
		return response.NewError(http.StatusInternalServerError, "an error occured when upload photo", err)
	}

	err = u.mentorRepository.UpdateMentor(&domain.MentorUpdate{MentorPicture: newMentorPicture}, mentor.Id)
	if err != nil {
		return response.NewError(http.StatusInternalServerError, "an error occured when update mentor", err)
	}

	var updatedMentor domain.Mentors
	err = u.mentorRepository.GetMentor(&updatedMentor, mentorParam)
	if err != nil {
		return response.NewError(http.StatusInternalServerError, "an error occured when get updated mentor", err)
	}

	return nil
}
