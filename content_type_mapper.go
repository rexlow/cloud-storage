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
	ContentTypeAPK  = "vnd.android.package-archive"
	ContentTypeAny  = ""
)

var contentTypeMapper = map[string]func(sw *storage.Writer){
	ContentTypeCSV:  contentTypeCSV,
	ContentTypePNG:  contentTypePNG,
	ContentTypeJPEG: contentTypeJPEG,
	ContentTypePDF:  contentTypePDF,
	ContentTypeGZip: contentTypeGZip,
	ContentTypeAPK:  contentTypeApk,
	ContentTypeAny:  contentTypeAny,
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

func contentTypeApk(sw *storage.Writer) {
	sw.ContentType = "application/vnd.android.package-archive"
}
