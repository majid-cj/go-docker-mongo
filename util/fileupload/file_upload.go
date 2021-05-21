package fileupload

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/majid-cj/go-docker-mongo/util"
	"github.com/thoas/go-funk"
)

// AllowedImages ....
var AllowedImages = []string{"image/jpeg", "image/jpg", "image/png"}

// UploadFile ...
type UploadFile struct{}

// UploadFileInterface ...
type UploadFileInterface interface {
	UploadFile(*multipart.FileHeader, multipart.File, string, string) (string, error)
}

var _ UploadFileInterface = &UploadFile{}

// NewUploadFile ...
func NewUploadFile() *UploadFile {
	return &UploadFile{}
}

// UploadFile ...
func (uf *UploadFile) UploadFile(fileHeader *multipart.FileHeader, file multipart.File, dest, host string) (string, error) {
	fileHeader.Filename = FormatFile(fileHeader.Filename)
	src, err := fileHeader.Open()

	if err != nil {
		return "", util.GetError("general_error")
	}

	defer src.Close()

	size := fileHeader.Size

	if size > 10<<20 {
		return "", util.GetError("image_size_error")
	}

	buffer := make([]byte, size)
	_, err = file.Read(buffer)
	if err != nil {
		return "", util.GetError("general_error")
	}
	filetype := http.DetectContentType(buffer)

	if !funk.ContainsString(AllowedImages, filetype) {
		return "", util.GetError("not_supported_type")
	}

	os.Mkdir(dest, os.FileMode(0766))
	output, err := os.OpenFile(filepath.Join(dest, fileHeader.Filename), os.O_RDWR|os.O_CREATE, os.FileMode(0766))

	if err != nil {
		return "", util.GetError("general_error")
	}

	defer output.Close()

	fileHeader.Filename = fmt.Sprintf("%s/%s%s", host, dest, fileHeader.Filename)
	_, err = io.Copy(output, src)
	if err != nil {
		return "", util.GetError("general_error")
	}
	return fileHeader.Filename, nil
}
