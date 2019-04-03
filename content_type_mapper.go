package storage

import (
	"fmt"
	"strings"

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
	ContentTypeHTML = "html"
	ContentTypeCSS  = "css"
	ContentTypeJS   = "js"
	ContentTypeExcel = "xlsx"
	ContentTypeAny  = ""
)

var contentTypeMapper = map[string]func(sw *storage.Writer){
	ContentTypeCSV:  contentTypeCSV,
	ContentTypePNG:  contentTypePNG,
	ContentTypeJPEG: contentTypeJPEG,
	ContentTypePDF:  contentTypePDF,
	ContentTypeGZip: contentTypeGZip,
	ContentTypeAPK:  contentTypeApk,
	ContentTypeHTML: contentTypeHTML,
	ContentTypeCSS:  contentTypeCSS,
	ContentTypeJS:   contentTypeJS,
	ContentTypeExcel: contentTypeExcel,
	// ContentTypeAny:  contentTypeAny,
}

func contentTypeGZip(sw *storage.Writer) {
	sw.ContentType = "application/zip"
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
	filenameArr := strings.Split(sw.Name, ".")
	fileExtension := strings.TrimSpace(filenameArr[len(filenameArr)-1])

	sw.ContentDisposition = fmt.Sprintf("attachment;filename=%s", sw.Name)

	contentFunc, isExist := contentTypeMapper[fileExtension]
	if isExist {
		contentFunc(sw)
	}
}

func contentTypeApk(sw *storage.Writer) {
	sw.ContentType = "application/vnd.android.package-archive"
	sw.ContentDisposition = fmt.Sprintf("attachment;filename=%s", sw.Name)
}

func contentTypeHTML(sw *storage.Writer) {
	sw.ContentType = "text/html"
}

func contentTypeCSS(sw *storage.Writer) {
	sw.ContentType = "text/css"
}

func contentTypeJS(sw *storage.Writer) {
	sw.ContentType = "application/javascript"
}

func contentTypeExcel(sw *storage.Writer) {
	sw.ContentType = "application/vnd.ms-excel"
	sw.ContentDisposition = fmt.Sprintf("attachment;filename=%s", sw.Name)
}
