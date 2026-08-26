package helpers

import "regexp"

var MigrationRegex = regexp.MustCompile(`^[0-9]{5,}-[A-Za-z0-9_-]+\.sql$`)
