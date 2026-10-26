package helpers

import "regexp"

const AppName = "OnlyStudents"

var EmailRegex = regexp.MustCompile(`^[^@\s]+@[^@\s]+\.[^@\s]+$`)
var PhoneNumberRegex = regexp.MustCompile(`^(\+\d{1,2}\s?)?\(?\d{3}\)?[\s.-]?\d{3}[\s.-]?\d{4}$`)

var Languages = []string{"", "en-US", "hu-HU"}
var TimeFormats = []string{"", "h12", "h23"}
var FrontendPaths = []string{
	"/",
	"/admin",
	"/me",
	"/enroll",
	"/timetable",
	"/grades",
	"/exams",
	"/homework",
	"/absences",
}

var sniffedMIMEs = map[FileCategory][]string{
	CategoryVideo: {
		"video/mp4",
		"video/avi",
		"video/webm", // also matches .mkv (same EBML magic)
	},
	CategoryImage: {
		"image/jpeg",
		"image/png",
		"image/webp",
		"image/gif",
		"image/bmp",
		"image/x-icon", // .ico
	},
	CategoryAudio: {
		"audio/wave",      // .wav
		"application/ogg", // .ogg
	},
	CategoryDocument: {
		"text/plain; charset=utf-8", // txt, md, csv
		"application/zip",           // docx, odt, epub sniff as zip
		"application/pdf",
	},
	CategoryArchive: {
		"application/zip",
		"application/x-gzip",           // gzip, tar.gz, tgz
		"application/x-rar-compressed", // rar
	},
	CategorySpreadSheet: {
		"application/zip", // xlsx, ods sniff as zip
	},
	CategoryPresentation: {
		"application/zip", // pptx, odp sniff as zip
	},
}
