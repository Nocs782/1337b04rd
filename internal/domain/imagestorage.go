package domain

import "mime/multipart"

type ImageStorage interface {
	UploadImage(file multipart.File, filename string) error
	DeleteImage(filename string) error
}
