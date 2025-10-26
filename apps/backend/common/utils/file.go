package utils

import (
	"mime/multipart"
	"path/filepath"
)

func GetFileContentType(fileHeader *multipart.FileHeader) string {
	return fileHeader.Header.Get("Content-Type")
}

func GetFileExtension(fileHeader *multipart.FileHeader) string {
	return filepath.Ext(fileHeader.Filename)
}

func OpenFile(fileHeader *multipart.FileHeader) (multipart.File, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return nil, err
	}
	return file, nil
}
