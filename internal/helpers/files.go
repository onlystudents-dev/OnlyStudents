package helpers

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/mojatter/s2"
	_ "github.com/mojatter/s2/azblob" // Azure Blob Storage
	_ "github.com/mojatter/s2/fs"     // osfs + memfs
	_ "github.com/mojatter/s2/gcs"    // Google Cloud Storage
	_ "github.com/mojatter/s2/s3"     // AWS S3
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

func LoadFileUploadBackend() (s2.Storage, error) {
	file_upload_backend := strings.ToLower(os.Getenv("FILE_UPLOAD_BACKEND"))

	var strg s2.Storage
	ctx := context.Background()

	switch file_upload_backend {
	case "directory":
		return s2.NewStorage(ctx, s2.Config{
			Type: s2.TypeOSFS,
			Root: GetEnvFallback("FILE_UPLOADS_DIR", "./uploads"),
		})
	case "memory":
		slog.Warn("Using memory FS for file storage, not recommended in production")
		return s2.NewStorage(ctx, s2.Config{
			Type: s2.TypeMemFS,
		})
	case "gcs":
		bucketName := GetEnvFallback("GCS_BUCKET", "onlystudents_files")

		return s2.NewStorage(ctx, s2.Config{
			Type: s2.TypeGCS,
			Root: bucketName,
			GCS: &s2.GCSConfig{
				CredentialsFile: os.Getenv("GCS_CREDENTIAL_FILE_PATH"),
			},
		})
	case "azure":
		accountName := os.Getenv("AZURE_ACCOUNT_NAME")
		accountKey := os.Getenv("AZURE_ACCOUNT_KEY")
		connectionString := os.Getenv("AZURE_CONNECTION_STRING")

		return s2.NewStorage(ctx, s2.Config{
			Type: s2.TypeAzblob,
			Azblob: &s2.AzblobConfig{
				AccountName:      accountName,
				AccountKey:       accountKey,
				ConnectionString: connectionString,
			},
		})
	case "s3":
		endpoint := os.Getenv("S3_ENDPOINT")
		accessKeyID := os.Getenv("S3_ACCESS_KEY_ID")
		secretAccessKey := os.Getenv("S3_SECRET_ACCESS_KEY")
		bucketName := GetEnvFallback("S3_BUCKET", "onlystudents_files")

		return s2.NewStorage(ctx, s2.Config{
			Type: s2.TypeS3,
			Root: bucketName,
			S3: &s2.S3Config{
				EndpointURL:     endpoint,
				AccessKeyID:     accessKeyID,
				SecretAccessKey: secretAccessKey,
			},
		})
	default:
		return strg, errors.New("Invalid file provider")
	}
}

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
