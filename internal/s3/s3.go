package s3

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type S3Client struct {
	bucketPath string
}

func NewS3Client(bucket string) (*S3Client, error) {
	// If the bucket directory doesn't exist, create it
	if _, err := os.Stat(bucket); os.IsNotExist(err) {
		err = os.Mkdir(bucket, 0755)
		if err != nil {
			return nil, fmt.Errorf("failed to create bucket directory: %v", err)
		}
	}

	return &S3Client{
		bucketPath: bucket,
	}, nil
}

// UploadImage copies a local file into the bucket folder
func (s *S3Client) UploadImage(sourcePath, objectName string) error {
	srcFile, err := os.Open(sourcePath)
	if err != nil {
		return fmt.Errorf("failed to open source file: %v", err)
	}
	defer srcFile.Close()

	dstPath := filepath.Join(s.bucketPath, objectName)
	dstFile, err := os.Create(dstPath)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %v", err)
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	return err
}

func (s *S3Client) DownloadImage(objectName, destPath string) error {
	srcPath := filepath.Join(s.bucketPath, objectName)
	srcFile, err := os.Open(srcPath)
	if err != nil {
		return fmt.Errorf("failed to open source file in bucket: %v", err)
	}
	defer srcFile.Close()

	dstFile, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %v", err)
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	return err
}

func (s *S3Client) DeleteImage(objectName string) error {
	return os.Remove(filepath.Join(s.bucketPath, objectName))
}
