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
	CategoryVideo        FileCategory = "video"
	CategoryImage        FileCategory = "image"
	CategoryAudio        FileCategory = "audio"
	CategoryDocument     FileCategory = "document"
	CategoryArchive      FileCategory = "archive"
	CategorySpreadSheet  FileCategory = "spreadsheet"
	CategoryPresentation FileCategory = "presentation"
)

var FileCategories = map[FileCategory][]string{
	CategoryVideo:        {".mp4", ".avi", ".mkv", ".webm"},
	CategoryImage:        {".jpg", ".jpeg", ".png", ".webp", ".gif", ".bmp"},
	CategoryAudio:        {".mp3", ".wav", ".ogg", ".m4a", ".flac", ".aac", ".opus"},
	CategoryDocument:     {".docx", ".txt", ".pdf", ".doc", ".odt", ".rtf", ".md", ".csv", ".epub"},
	CategoryArchive:      {".zip", ".7z", ".tar", ".tar.gz", ".tar.xz", ".rar", ".gz", ".xz", ".bz2", ".tgz"},
	CategorySpreadSheet:  {".xls", ".xlsx", ".ods"},
	CategoryPresentation: {".pptx", ".ppt", ".odp"},
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

var sniffOpaqueExts = map[string]struct{}{
	// archive containers without signatures in Go's sniffer
	".7z":     {},
	".tar":    {},
	".tar.xz": {},
	".xz":     {},
	".bz2":    {},

	// legacy binary Office documents (OLE compound files)
	".doc": {},
	".xls": {},
	".ppt": {},

	// richt text is text-based but may contain binary escapes
	".rtf": {},

	// audio containers (Go only sniffs ID3-tagged MP3 / WAV / OGG)
	".mp3":  {},
	".m4a":  {},
	".flac": {},
	".aac":  {},
	".opus": {},
}

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
	ext := FileExtOf(file.Filename)
	cat, ok := CategoryByExt[ext]
	if !ok {
		return "", fmt.Errorf("file extension not allowed")
	}
	if _, opaque := sniffOpaqueExts[ext]; opaque {
		return cat, nil
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
