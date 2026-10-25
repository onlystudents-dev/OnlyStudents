package helpers

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"
)

type FileCategory string

const (
	CategoryImage        FileCategory = "image"
	CategoryDocument     FileCategory = "document"
	CategoryArchive      FileCategory = "archive"
	CategorySpreadSheet  FileCategory = "spreadsheet"
	CategoryPresentation FileCategory = "presentation"
)

var FileCategories = map[FileCategory][]string{
	CategoryImage:        {".jpg", ".jpeg", ".png", ".webp"},
	CategoryDocument:     {".docx", ".txt"},
	CategoryArchive:      {".zip", ".7z", ".tar", ".tar.gz", ".tar.xz", ".rar"},
	CategorySpreadSheet:  {".xls", ".xlsx"},
	CategoryPresentation: {".pptx", ".ppt"},
}

var CategoryByExt = map[string]FileCategory{}

func InitFileExts() {
	for cat, exts := range FileCategories {
		for _, ext := range exts {
			CategoryByExt[strings.ToLower(ext)] = cat
		}
	}
}

var compositeExts = []string{".tar.gz", ".tar.xz"}

func FileExtOf(name string) string {
	name = strings.ToLower(name)
	for _, ext := range compositeExts {
		if strings.HasSuffix(name, ext) {
			return ext
		}
	}
	return filepath.Ext(name)
}

func ValidateFileAndCategorize(file *multipart.FileHeader) (FileCategory, error) {
	cat, ok := CategoryByExt[FileExtOf(file.Filename)]
	if !ok {
		return "", fmt.Errorf("file extension not allowed")
	}

	f, err := file.Open()
	if err != nil {
		return "", err
	}
	defer f.Close()

	buf := make([]byte, 512)
	n, err := f.Read(buf)
	if err != nil && err != io.EOF {
		return "", err
	}

	mimeType := http.DetectContentType(buf[:n])
	for _, allowed := range sniffedMIMEs[cat] {
		if mimeType == allowed {
			return cat, nil
		}
	}
	return "", fmt.Errorf("content type %s does not match %s files", mimeType, cat)
}
