package storage

import (
	"errors"
	"io"
	"mime/multipart"
	"strings"
	"time"
)

// Builder :
type Builder struct {
	adapter Adapter
	err     error
}

// AliyunClient :
type AliyunClient struct {
	Endpoint        string
	AccessKeyID     string
	AccessKeySecret string
}

// GoogleCloudStorageClient :
type GoogleCloudStorageClient struct {
}

var client = map[string]bool{
	GCS: true,
}

// NewClient :
func NewClient(name string) *Builder {
	builder := new(Builder)
	if !client[strings.ToUpper(name)] {
		builder.err = errors.New("The client not supported")
	}

	switch strings.ToUpper(name) {
	case GCS:
		adapter := new(GCSAdapter)
		builder.adapter = adapter
		break
	}

	return builder
}

// New :
func New(client interface{}) *Builder {
	builder := new(Builder)

	switch v := client.(type) {
	case AliyunClient:
		adapter := new(AliyunAdapter)
		adapter.Endpoint = v.Endpoint
		adapter.AccessKeyID = v.AccessKeyID
		adapter.AccessKeySecret = v.AccessKeySecret
		builder.adapter = adapter

	case GoogleCloudStorageClient:
		adapter := new(GCSAdapter)
		builder.adapter = adapter

	default:
		builder.err = errors.New("invalid client interface")
		return builder
	}

	return builder
}

// UploadFile :
func (b *Builder) UploadFile(file *multipart.FileHeader, bucket, name string) (string, error) {
	if b.err != nil {
		return "", b.err
	}
	return b.adapter.UploadFile(file, bucket, name)
}

// ReadFile :
func (b *Builder) ReadFile(bucket, path string) ([]byte, error) {
	if b.err != nil {
		return nil, b.err
	}
	return b.adapter.ReadFile(bucket, path)
}

// DeleteFileUsingURL :
func (b *Builder) DeleteFileUsingURL(bucket, fileURL string) error {
	if b.err != nil {
		return b.err
	}
	return b.adapter.DeleteFileUsingURL(bucket, fileURL)
}

// TemporaryServingFile :
func (b *Builder) TemporaryServingFile(bucket, fileURL string, expiredTime time.Time, client interface{}) (string, error) {
	if b.err != nil {
		return "", b.err
	}
	return b.adapter.TemporaryServingFile(bucket, fileURL, expiredTime, client)
}

// GoogleTemporaryServingFile :
func GoogleTemporaryServingFile(bucket, fileURL string, expiredTime time.Time, client GoogleClient) (string, error) {
	builder := new(Builder)
	adapter := new(GCSAdapter)
	builder.adapter = adapter
	return builder.adapter.TemporaryServingFile(bucket, fileURL, expiredTime, client)
}

// UploadReader :
func (b *Builder) UploadReader(bucket, filename string, reader io.Reader, contentType string) (string, error) {
	if b.err != nil {
		return "", b.err
	}
	var (
		errNameIsRequired   = errors.New("storage: filename is required")
		errBucketIsRequired = errors.New("storage: bucket is required")
		errReaderIsNil      = errors.New("storage: io reader is nil")
	)

	if len(filename) == 0 {
		return "", errNameIsRequired
	}

	if len(bucket) == 0 {
		return "", errBucketIsRequired
	}

	if reader == nil {
		return "", errReaderIsNil
	}

	return b.adapter.UploadReader(bucket, filename, reader, contentType)
}

// UploadBuffer :
func (b *Builder) UploadBuffer(bucket, filename string, contentType string) (*Buffer, error) {
	if b.err != nil {
		return nil, b.err
	}
	var (
		errNameIsRequired   = errors.New("storage: filename is required")
		errBucketIsRequired = errors.New("storage: bucket is required")
	)

	if len(filename) == 0 {
		return nil, errNameIsRequired
	}

	if len(bucket) == 0 {
		return nil, errBucketIsRequired
	}

	return b.adapter.UploadBuffer(bucket, filename, contentType)
}
