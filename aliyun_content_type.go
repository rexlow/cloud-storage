package storage

import (
	"fmt"
	"strings"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

var aliyunContentTypeMapper = map[string]func(filename string) []oss.Option{
	ContentTypeCSV:   aliyunContentTypeCSV,
	ContentTypePNG:   aliyunContentTypePNG,
	ContentTypeJPEG:  aliyunContentTypeJPEG,
	ContentTypeJPG:   aliyunContentTypeJPEG,
	ContentTypeHEIC:  aliyunContentTypeHeic,
	ContentTypePDF:   aliyunContentTypePDF,
	ContentTypeZip:   aliyunContentTypeZip,
	ContentTypeAPK:   aliyunContentTypeApk,
	ContentTypeHTML:  aliyunContentTypeHTML,
	ContentTypeCSS:   aliyunContentTypeCSS,
	ContentTypeJS:    aliyunContentTypeJS,
	ContentTypeExcel: aliyunContentTypeExcel,
	ContentTypeSVG:   aliyunContentTypeSVG,
}

func aliyunContentTypeZip(filename string) []oss.Option {
	options := make([]oss.Option, 0)
	options = append(options, oss.ContentType("application/zip"))
	options = append(options, oss.ContentDisposition(fmt.Sprintf("attachment;filename=%s", filename)))

	return options
}

func aliyunContentTypeCSV(filename string) []oss.Option {
	options := make([]oss.Option, 0)
	options = append(options, oss.ContentType("text/csv"))
	options = append(options, oss.ContentDisposition(fmt.Sprintf("attachment;filename=%s", filename)))

	return options
}

func aliyunContentTypePNG(filename string) []oss.Option {
	options := make([]oss.Option, 0)
	options = append(options, oss.ContentType("image/png"))

	return options
}

func aliyunContentTypeJPEG(filename string) []oss.Option {
	options := make([]oss.Option, 0)
	options = append(options, oss.ContentType("image/jpeg"))

	return options
}

func aliyunContentTypeHeic(filename string) []oss.Option {
	options := make([]oss.Option, 0)
	options = append(options, oss.ContentType("image/heic"))

	return options
}

func aliyunContentTypePDF(filename string) []oss.Option {
	options := make([]oss.Option, 0)
	options = append(options, oss.ContentType("application/pdf"))

	return options
}

func aliyunContentTypeAny(filename string) []oss.Option {
	filenameArr := strings.Split(filename, ".")
	fileExtension := strings.TrimSpace(filenameArr[len(filenameArr)-1])

	options := make([]oss.Option, 0)
	// options = append(options, oss.ContentDisposition(fmt.Sprintf("attachment;filename=%s", filename)))

	contentFunc, isExist := aliyunContentTypeMapper[fileExtension]
	if isExist {
		o := contentFunc(filename)
		options = append(options, o...)
	}

	return options
}

func aliyunContentTypeApk(filename string) []oss.Option {
	options := make([]oss.Option, 0)
	options = append(options, oss.ContentType("application/vnd.android.package-archive"))
	options = append(options, oss.ContentDisposition(fmt.Sprintf("attachment;filename=%s", filename)))

	return options
}

func aliyunContentTypeHTML(filename string) []oss.Option {
	options := make([]oss.Option, 0)
	options = append(options, oss.ContentType("text/html"))

	return options
}

func aliyunContentTypeCSS(filename string) []oss.Option {
	options := make([]oss.Option, 0)
	options = append(options, oss.ContentType("text/css"))

	return options
}

func aliyunContentTypeJS(filename string) []oss.Option {
	options := make([]oss.Option, 0)
	options = append(options, oss.ContentType("application/javascript"))

	return options
}

func aliyunContentTypeExcel(filename string) []oss.Option {
	options := make([]oss.Option, 0)
	options = append(options, oss.ContentType("application/vnd.ms-excel"))
	options = append(options, oss.ContentDisposition(fmt.Sprintf("attachment;filename=%s", filename)))

	return options
}

func aliyunContentTypeSVG(filename string) []oss.Option {
	options := make([]oss.Option, 0)
	options = append(options, oss.ContentType("image/svg+xml"))

	return options
}
