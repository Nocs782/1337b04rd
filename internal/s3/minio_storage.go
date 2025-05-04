package s3

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strings"
)

type MinioStorage struct {
	endpoint string // Example: http://minio:9000
	bucket   string // Example: post-images
}

func NewMinioStorage(endpoint, bucket string) *MinioStorage {
	return &MinioStorage{
		endpoint: strings.TrimSuffix(endpoint, "/"),
		bucket:   bucket,
	}
}

func (m *MinioStorage) UploadImage(file multipart.File, filename string) error {
	defer file.Close()

	var buf bytes.Buffer
	_, err := io.Copy(&buf, file)
	if err != nil {
		return fmt.Errorf("failed to read file into buffer: %v", err)
	}

	url := fmt.Sprintf("%s/%s/%s", m.endpoint, m.bucket, filename)

	req, err := http.NewRequest(http.MethodPut, url, &buf)
	if err != nil {
		return fmt.Errorf("failed to create PUT request: %v", err)
	}

	req.ContentLength = int64(buf.Len())

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("PUT request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("PUT failed: %s", resp.Status)
	}

	return nil
}

func (m *MinioStorage) DownloadImage(filename string) ([]byte, error) {
	url := fmt.Sprintf("%s/%s/%s", m.endpoint, m.bucket, filename)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("GET request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GET failed: %s", resp.Status)
	}

	return io.ReadAll(resp.Body)
}

func (m *MinioStorage) DeleteImage(filename string) error {
	url := fmt.Sprintf("%s/%s/%s", m.endpoint, m.bucket, filename)

	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return fmt.Errorf("failed to create DELETE request: %v", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("DELETE request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("DELETE failed: %s", resp.Status)
	}

	return nil
}
