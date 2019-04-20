package storage

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"

	s "cloud.google.com/go/storage"
	oss "github.com/aliyun/aliyun-oss-go-sdk/oss"
)

// Buffer :
type Buffer struct {
	adapter       string      // gcs and aliyun
	bucket        string      // gcs and aliyun
	filename      string      // gcs and aliyun
	storageWriter *s.Writer   // gcs
	object        *oss.Bucket // aliyun
	endpoint      string      // aliyun
	position      int64       // aliyun
}

// Copy :
func (buf *Buffer) Copy(reader io.Reader) error {
	switch buf.adapter {
	case GCS:
		if _, err := io.Copy(buf.storageWriter, reader); err != nil {
			msg := fmt.Sprintf("Could not write file: %v", err)
			return errors.New(msg)
		}

	case ALIYUN:
		position, err := buf.object.AppendObject(buf.filename, reader, buf.position)
		if err != nil {
			msg := fmt.Sprintf("Could not write file: %v", err)
			return errors.New(msg)
		}
		buf.position = position
	default:
		return errors.New("invalid adapter")
	}

	return nil
}

// CopyByte :
func (buf *Buffer) CopyByte(data []byte) error {
	d := bytes.NewReader(data)
	return buf.Copy(d)
}

// CopyString :
func (buf *Buffer) CopyString(data string) error {
	d := strings.NewReader(data)
	return buf.Copy(d)
}

// Close :
func (buf *Buffer) Close() (string, error) {
	switch buf.adapter {
	case GCS:
		if err := buf.storageWriter.Close(); err != nil {
			msg := fmt.Sprintf("Could not put file: %v", err)
			return "", errors.New(msg)
		}

		return fmt.Sprintf("%s/%s/%s", googleGCSDomain, buf.bucket, buf.storageWriter.Name), nil

	case ALIYUN:
		return getAliyunFileURL(buf.endpoint, buf.bucket, buf.filename), nil

	default:
		return "", errors.New("invalid adapter")
	}
}
