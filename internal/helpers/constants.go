package helpers

import "regexp"

const AppName = "OnlyStudents"

var MigrationRegex = regexp.MustCompile(`^[0-9]{5,}-[A-Za-z0-9_-]+\.sql$`)
var EmailRegex = regexp.MustCompile(`^[^@\s]+@[^@\s]+\.[^@\s]+$`)
var PhoneNumberRegex = regexp.MustCompile(`^(\+\d{1,2}\s?)?\(?\d{3}\)?[\s.-]?\d{3}[\s.-]?\d{4}$`)

var Languages = []string{"", "en-US", "hu-HU"}
var TimeFormats = []string{"", "h12", "h23"}
var FrontendPaths = []string{
	"/",
	"/admin",
	"/me",
	"/timetable",
	"/grades",
	"/homeworks",
	"/absences",
}
