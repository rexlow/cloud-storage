package storage

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"strings"
	"time"

	s "cloud.google.com/go/storage"
)

const (
	googleGCSDomain = "https://storage.googleapis.com"
)

// GoogleClient :
type GoogleClient struct {
	Type                    string `json:"type"`
	ProjectID               string `json:"project_id"`
	PrivateKeyID            string `json:"private_key_id"`
	PrivateKey              string `json:"private_key"`
	ClientID                string `json:"client_id"`
	ClientEmail             string `json:"client_email"`
	AuthURI                 string `json:"auth_uri"`
	TokenURI                string `json:"token_uri"`
	AuthProviderX509CertURL string `json:"auth_provider_x509_cert_url"`
	ClientX509CertURL       string `json:"client_x509_cert_url"`
}

// GCSAdapter : The adapter empty because they do not have any data to share for now
type GCSAdapter struct{}

// UploadFile : Upload file to the bucket
func (adapter *GCSAdapter) UploadFile(file *multipart.FileHeader, bucket, filename string) (string, error) {
	var reader io.Reader
	name := ""

	fileExt := strings.ToLower(strings.Split(file.Header["Content-Type"][0], "/")[1])

	if ext, isExist := extensionMapper[fileExt]; isExist {
		name = fmt.Sprintf("%s.%s", filename, ext)
	} else {
		name = fmt.Sprintf("%s.%s", filename, fileExt)
	}

	src, err := file.Open()
	if err != nil {
		return "", err
	}

	defer src.Close()
	reader = src

	fileURL, err := adapter.UploadReader(bucket, name, reader, strings.ToLower(fileExt))
	if err != nil {
		return "", err
	}

	return fileURL, nil
}

// DeleteFileUsingURL : Delete file from the bucket using url
func (adapter *GCSAdapter) DeleteFileUsingURL(bucket, fileURL string) error {
	ctx := context.Background()
	storageClient, err := s.NewClient(ctx)
	if err != nil {
		return err
	}
	defer storageClient.Close()

	fileName := strings.Replace(fileURL, fmt.Sprintf("%s/%s/", googleGCSDomain, bucket), "", -1)

	return storageClient.Bucket(bucket).Object(fileName).Delete(ctx)
}

// TemporaryServingFile : TemporaryServingFile file serving
func (adapter *GCSAdapter) TemporaryServingFile(bucket, fileURL string, expiredDateTime time.Time, googleClient interface{}) (string, error) {
	credential := googleClient.(GoogleClient)

	method := "GET"

	fileName := strings.Replace(fileURL, fmt.Sprintf("%s/%s/", googleGCSDomain, bucket), "", -1)
	url, err := s.SignedURL(bucket, fileName, &s.SignedURLOptions{
		GoogleAccessID: credential.ClientEmail,
		PrivateKey:     []byte(credential.PrivateKey),
		Method:         method,
		Expires:        expiredDateTime,
	})

	if err != nil {
		return "", err
	}

	return url, nil
}

// UploadReader :
func (adapter *GCSAdapter) UploadReader(bucket, filename string, reader io.Reader, contentType string) (string, error) {

	ctx := context.Background()
	storageClient, err := s.NewClient(ctx)
	if err != nil {
		return "", err
	}

	sw := storageClient.Bucket(bucket).Object(filename).NewWriter(ctx)

	if contentType == "" {
		contentTypeAny(sw)
	} else {
		contentFunc, isExist := contentTypeMapper[contentType]
		if !isExist {
			return "", fmt.Errorf("Content type %s does not supported", contentType)
		}
		contentFunc(sw)
	}

	if _, err := io.Copy(sw, reader); err != nil {
		msg := fmt.Sprintf("Could not write file: %v", err)
		return "", errors.New(msg)
	}

	if err := sw.Close(); err != nil {
		msg := fmt.Sprintf("Could not put file: %v", err)
		return "", errors.New(msg)
	}

	return fmt.Sprintf("%s/%s/%s", googleGCSDomain, bucket, filename), nil
}

// ReadFile :
func (adapter *GCSAdapter) ReadFile(bucket, path string) ([]byte, error) {
	ctx := context.Background()
	client, err := s.NewClient(ctx)
	if err != nil {
		return nil, err
	}
	defer client.Close()
	rc, err := client.Bucket(bucket).Object(path).NewReader(ctx)
	if err != nil {
		return nil, err
	}
	defer rc.Close()
	slurp, err := ioutil.ReadAll(rc)
	if err != nil {
		return nil, err
	}
	return slurp, nil
}

// Buffer :
type Buffer struct {
	storageWriter *s.Writer
	bucket        string
}

// UploadBuffer :
func (adapter *GCSAdapter) UploadBuffer(bucket, filename string, contentType string) (*Buffer, error) {
	buf := new(Buffer)

	ctx := context.Background()
	storageClient, err := s.NewClient(ctx)
	if err != nil {
		return nil, err
	}

	sw := storageClient.Bucket(bucket).Object(filename).NewWriter(ctx)
	buf.storageWriter = sw
	buf.bucket = bucket
	contentTypeMapper[contentType](sw)

	return buf, nil
}

// Copy :
func (buf *Buffer) Copy(reader io.Reader) error {
	if _, err := io.Copy(buf.storageWriter, reader); err != nil {
		msg := fmt.Sprintf("Could not write file: %v", err)
		return errors.New(msg)
	}

	return nil
}

// CopyByte :
func (buf *Buffer) CopyByte(data []byte) error {
	d := bytes.NewReader(data)
	if _, err := io.Copy(buf.storageWriter, d); err != nil {
		msg := fmt.Sprintf("Could not write file: %v", err)
		return errors.New(msg)
	}

	return nil
}

// CopyString :
func (buf *Buffer) CopyString(data string) error {
	d := strings.NewReader(data)
	if _, err := io.Copy(buf.storageWriter, d); err != nil {
		msg := fmt.Sprintf("Could not write file: %v", err)
		return errors.New(msg)
	}

	return nil
}

// Close :
func (buf *Buffer) Close() (string, error) {
	if err := buf.storageWriter.Close(); err != nil {
		msg := fmt.Sprintf("Could not put file: %v", err)
		return "", errors.New(msg)
	}

	return fmt.Sprintf("%s/%s/%s", googleGCSDomain, buf.bucket, buf.storageWriter.Name), nil
}
