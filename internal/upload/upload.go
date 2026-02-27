package upload

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const uploadDir = "./uploads"

func init() {
	os.MkdirAll(uploadDir, os.ModePerm)
}

func SaveFile(file multipart.File, header *multipart.FileHeader, fileType string) (string, error) {
	ext := filepath.Ext(header.Filename)
	filename := fmt.Sprintf("%d_%s%s", time.Now().UnixNano(), strings.ReplaceAll(fileType, "/", "_"), ext)
	filePath := filepath.Join(uploadDir, filename)

	dst, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		return "", err
	}

	return "/uploads/" + filename, nil
}
