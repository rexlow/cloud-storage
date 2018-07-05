package storage

import (
	"fmt"

	"cloud.google.com/go/storage"
)

// Content Type
const (
	ContentTypeCSV  = "csv"
	ContentTypePNG  = "png"
	ContentTypeJPEG = "jpeg"
	ContentTypePDF  = "pdf"
	ContentTypeGZip = "gzip"
	ContentTypeAny = ""
)

var contentTypeMapper = map[string]func(sw *storage.Writer){
	ContentTypeCSV:  contentTypeCSV,
	ContentTypePNG:  contentTypePNG,
	ContentTypeJPEG: contentTypeJPEG,
	ContentTypePDF:  contentTypePDF,
	ContentTypeGZip: contentTypeGZip,
	ContentTypeAny: contentTypeAny,
}

func contentTypeGZip(sw *storage.Writer) {
	sw.ContentDisposition = fmt.Sprintf("attachment;filename=%s", sw.Name)
}

func contentTypeCSV(sw *storage.Writer) {
	sw.ContentType = "text/csv"
	sw.ContentDisposition = fmt.Sprintf("attachment;filename=%s", sw.Name)
}

func contentTypePNG(sw *storage.Writer) {
	sw.ContentType = "image/png"
}

func contentTypeJPEG(sw *storage.Writer) {
	sw.ContentType = "image/jpeg"
}

func contentTypePDF(sw *storage.Writer) {
	sw.ContentType = "application/pdf"
}

func contentTypeAny(sw *storage.Writer) {
	sw.ContentType = "application/pdf"
}
