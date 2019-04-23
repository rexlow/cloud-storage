package storage

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"regexp"
	"strings"
	"time"

	oss "github.com/aliyun/aliyun-oss-go-sdk/oss"
)

// AliyunAdapter : The adapter empty because they do not have any data to share for now
type AliyunAdapter struct {
	Endpoint        string
	AccessKeyID     string
	AccessKeySecret string
}

// UploadFile : Upload file to the bucket
func (adapter *AliyunAdapter) UploadFile(file *multipart.FileHeader, bucket, filename string) (string, error) {
	var reader io.Reader
	name := ""

	fileExt := strings.ToLower(strings.Split(file.Header.Get("Content-Type"), "/")[1])

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
func (adapter *AliyunAdapter) DeleteFileUsingURL(bucket, fileURL string) error {
	storageClient, err := adapter.getClient()
	if err != nil {
		return err
	}

	object, err := storageClient.Bucket(bucket)
	if err != nil {
		return err
	}

	filepath := adapter.getFilePathFromURL(bucket, fileURL)

	return object.DeleteObject(filepath)
}

// TemporaryServingFile : TemporaryServingFile file serving
func (adapter *AliyunAdapter) TemporaryServingFile(bucket, fileURL string, expiredDateTime time.Time, aliClient interface{}) (string, error) {
	storageClient, err := adapter.getClient()
	if err != nil {
		return "", err
	}

	object, err := storageClient.Bucket(bucket)
	if err != nil {
		return "", err
	}

	filepath := adapter.getFilePathFromURL(bucket, fileURL)

	url, err := object.SignURL(filepath, http.MethodGet, int64(expiredDateTime.UTC().Sub(time.Now().UTC()).Seconds()))
	if err != nil {
		return "", err
	}

	return url, nil
}

// UploadReader :
func (adapter *AliyunAdapter) UploadReader(bucket, filename string, reader io.Reader, contentType string) (string, error) {
	storageClient, err := adapter.getClient()
	if err != nil {
		return "", err
	}

	object, err := storageClient.Bucket(bucket)
	if err != nil {
		return "", err
	}

	options := make([]oss.Option, 0)
	if contentType == "" {
		options = aliyunContentTypeAny(filename)
	} else {
		contentFunc, isExist := aliyunContentTypeMapper[contentType]
		if !isExist {
			return "", fmt.Errorf("Content type %s does not supported", contentType)
		}
		options = contentFunc(filename)
	}

	if err := object.PutObject(filename, reader, options...); err != nil {
		msg := fmt.Sprintf("Could not write file: %v", err)
		return "", errors.New(msg)
	}

	return getAliyunFileURL(adapter.Endpoint, bucket, filename), nil
}

// ReadFile :
func (adapter *AliyunAdapter) ReadFile(bucket, path string) ([]byte, error) {
	storageClient, err := adapter.getClient()
	if err != nil {
		return nil, err
	}

	object, err := storageClient.Bucket(bucket)
	if err != nil {
		return nil, err
	}

	rc, err := object.GetObject(path)
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

// UploadBuffer :
func (adapter *AliyunAdapter) UploadBuffer(bucket, filename string, contentType string) (*Buffer, error) {
	buf := new(Buffer)
	buf.adapter = ALIYUN

	storageClient, err := adapter.getClient()
	if err != nil {
		return nil, err
	}

	object, err := storageClient.Bucket(bucket)
	if err != nil {
		return nil, err
	}

	contentFunc, isExist := aliyunContentTypeMapper[contentType]
	if !isExist {
		return nil, fmt.Errorf("Content type %s does not supported", contentType)
	}
	options := contentFunc(filename)

	buffer := new(bytes.Buffer)

	// delete the object before append
	object.DeleteObject(filename)

	position, err := object.AppendObject(filename, buffer, buf.position, options...)
	if err != nil {
		log.Println("error: ", err)
		return nil, err
	}
	buf.position = position

	buf.position = position
	buf.object = object
	buf.filename = filename
	buf.bucket = bucket
	buf.endpoint = adapter.Endpoint
	return buf, nil
}

func (adapter *AliyunAdapter) getClient() (*oss.Client, error) {
	return oss.New(adapter.Endpoint, adapter.AccessKeyID, adapter.AccessKeySecret)
}

func getAliyunFileURL(endpoint, bucket string, filename string) string {
	endpoint = regexp.MustCompile(`^(http|https)://`).ReplaceAllString(endpoint, "")
	return fmt.Sprintf("https://%s.%s/%s", bucket, endpoint, filename)
}

func (adapter *AliyunAdapter) getFilePathFromURL(bucket, fileURL string) string {
	endpoint := regexp.MustCompile(`^(http|https)://`).ReplaceAllString(adapter.Endpoint, "")
	f := regexp.MustCompile(`^(http|https)://`).ReplaceAllString(fileURL, "")
	return strings.Replace(f, fmt.Sprintf("%s.%s/", bucket, endpoint), "", -1)
}
