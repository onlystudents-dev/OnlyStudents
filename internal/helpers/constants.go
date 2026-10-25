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
	CategoryImage: {
		"image/jpeg",
		"image/png",
		"image/webp",
	},
	CategoryDocument: {
		"text/plain; charset=utf-8",
		"application/zip", // docx sniffs as zip
	},
	CategoryArchive: {
		"application/zip",
		"application/gzip",
		"application/x-xz",
		"application/x-tar",
		"application/x-7z-compressed", // 7z
		"application/vnd.rar",         // rar
	},
	CategorySpreadSheet: {
		"application/vnd.ms-excel", // xls
		"application/zip",          // xlsx sniffs as zip
	},
	CategoryPresentation: {
		"application/vnd.ms-powerpoint", // ppt
		"application/zip",               // pptx sniffs as zip
	},
}
