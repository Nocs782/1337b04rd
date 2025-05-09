package domain

import "mime/multipart"

type ImageStorage interface {
	UploadImage(file multipart.File, filename string) error
	DownloadImage(filename string) ([]byte, error)
	DeleteImage(filename string) error
}
