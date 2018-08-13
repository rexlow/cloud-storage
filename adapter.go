package storage

import (
	"io"
	"mime/multipart"
	"time"
)

// Adapter :
type Adapter interface {
	UploadFile(file *multipart.FileHeader, bucket, name string) (string, error)
	DeleteFileUsingURL(bucket, fileURL string) error
	UploadReader(string, string, io.Reader, string) (string, error)
	TemporaryServingFile(bucket string, fileURL string, expiredTime time.Time, client interface{}) (string, error)
	UploadBuffer(string, string, string) (*Buffer, error)
	ReadFile(string, string) ([]byte, error)
}

var _ Adapter = &GCSAdapter{}
