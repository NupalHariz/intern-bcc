package supabase

import (
	"mime/multipart"
	"os"

	supabasestorageuploader "github.com/adityarizkyramadhan/supabase-storage-uploader"
)

type ISupabase interface {
	Upload(file *multipart.FileHeader) (string, error)
	Delete(link string) error
}

type Supabase struct {
	client *supabasestorageuploader.Client
}

func SupabaseInit() ISupabase {
	client := supabasestorageuploader.New(
		os.Getenv("SUPABASE_URL"),
		os.Getenv("SUPABASE_TOKEN"),
		os.Getenv("SUPABASE_BUCKET"),
	)

	return &Supabase{client}
}

func (s *Supabase) Upload(file *multipart.FileHeader) (string, error) {
	link, err := s.client.Upload(file)
	if err != nil {
		return link, err
	}

	return link, nil
}

func (s *Supabase) Delete(link string) error {
	err := s.client.Delete(link)
	if err != nil {
		return err
	}

	return nil
}
