package helpers

import "regexp"

var MigrationRegex = regexp.MustCompile(`^[0-9]{5,}-[A-Za-z0-9_-]+\.sql$`)
var EmailRegex = regexp.MustCompile(`^[^@\s]+@[^@\s]+\.[^@\s]+$`)
var Languages = []string{"", "en-US", "hu-HU"}
var Themes = []string{"", "dark", "light"}
