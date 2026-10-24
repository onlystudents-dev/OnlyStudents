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
	"/homeworks",
	"/absences",
}
