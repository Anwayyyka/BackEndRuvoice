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
	// создаём папку для загрузок, если её нет
	os.MkdirAll(uploadDir, os.ModePerm)
}

// SaveFile сохраняет файл на диск и возвращает путь к нему (относительный)
func SaveFile(file multipart.File, header *multipart.FileHeader, fileType string) (string, error) {
	// создаём уникальное имя файла
	ext := filepath.Ext(header.Filename)
	filename := fmt.Sprintf("%d_%s%s", time.Now().UnixNano(), strings.ReplaceAll(fileType, "/", "_"), ext)
	filePath := filepath.Join(uploadDir, filename)

	// создаём файл на диске
	dst, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		return "", err
	}

	// возвращаем относительный URL для доступа через статику
	return "/uploads/" + filename, nil
}
